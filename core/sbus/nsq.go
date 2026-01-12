package sbus

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/slog"
	"go.uber.org/zap"
)

type TopicEnum int32

var TopicEnum_name = map[int32]string{}

var TopicEnum_value = map[string]int32{}

func (t TopicEnum) String() string {
	return enumName(TopicEnum_name, int32(t))
}

func enumName(m map[int32]string, v int32) string {
	s, ok := m[v]
	if ok {
		return s
	}
	return strconv.Itoa(int(v))
}

func SetTopicEnumName(v map[int32]string) {
	TopicEnum_name = v
}

func SetTopicEnumValue(v map[string]int32) {
	TopicEnum_value = v
}

type NsqConsumer struct {
	topic         string
	channel       string
	nsqLookupAddr string
	concurrency   int
	nsqConsumer   *nsq.Consumer
	handler       nsq.Handler
}

func NewNsqConsumer(topic, channel, nsqLookupAddr string, concurrency, maxInFlight int) (*NsqConsumer, error) {
	nsqConsumer := &NsqConsumer{
		nsqLookupAddr: nsqLookupAddr,
		nsqConsumer:   nil,
		concurrency:   concurrency,
	}
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second
	cfg.LookupdPollTimeout = time.Millisecond * 25
	cfg.MaxInFlight = maxInFlight
	if cfg.MaxInFlight <= 0 {
		cfg.MaxInFlight = 100
	}
	if c, err := nsq.NewConsumer(topic, channel, cfg); err != nil {
		return nil, err
	} else {
		nsqConsumer.nsqConsumer = c
		return nsqConsumer, nil
	}
}

func (c *NsqConsumer) StartReader(handler nsq.Handler) error {
	c.nsqConsumer.AddConcurrentHandlers(handler, c.concurrency)
	return c.nsqConsumer.ConnectToNSQLookupd(c.nsqLookupAddr)
}

func (c *NsqConsumer) Stop() {
	c.nsqConsumer.Stop()
}

type NsqProducer struct {
	producer *nsq.Producer
}

func NewProducer(nsqdAddr string) (*NsqProducer, error) {
	cfg := nsq.NewConfig()
	cfg.LookupdPollInterval = time.Second // 设置重连时间
	cfg.OutputBufferTimeout = time.Millisecond * 25
	cfg.MaxInFlight = 64

	p, err := nsq.NewProducer(nsqdAddr, cfg)
	if err != nil {
		slog.Ins().Errorf("create nsq producer failed, err:%v", err)
		return nil, err
	}
	return &NsqProducer{producer: p}, nil
}

func (p *NsqProducer) PublishDirect(topic string, data []byte) error {
	if p.producer != nil {
		if data == nil { //不能发布空串，否则会导致error
			return fmt.Errorf("data is nil")
		}

		err := p.producer.Publish(topic, data) // 发布消息
		return err
	}
	return fmt.Errorf("producer is nil")
}

type NsqData struct {
	Topic string
	data  []byte
}

var NsqDataPool = new(sync.Pool)

func init() {
	NsqDataPool.New = func() interface{} {
		return allocateNsqData()
	}
}

func allocateNsqData() *NsqData {
	nsqData := new(NsqData)
	return nsqData
}
func (nd *NsqData) Reset(topic string, data []byte) {
	nd.Topic = topic
	nd.data = data
}

func GetNsqData(topic string, data []byte) *NsqData {

	// 根据当前模式判断是否使用对象池

	// 从对象池中取得一个 Request 对象,如果池子中没有可用的 Request 对象则会调用 allocateRequest 函数构造一个新的对象分配
	r := NsqDataPool.Get().(*NsqData)
	// 因为取出的 Request 对象可能是已存在也可能是新构造的,无论是哪种情况都应该初始化再返回使用
	r.Reset(topic, data)
	return r
}

func PutNsqData(nsqData *NsqData) {
	NsqDataPool.Put(nsqData)
}

type Nsq struct {
	//BaseConnection
	// The message management module that manages MsgID and the corresponding processing method
	// (消息管理MsgID和对应处理方法的消息管理模块)
	//taskHandler STaskHandler
	producers []*NsqProducer
	Consumers []*NsqConsumer

	// Buffered channel used for message communication between the read and write goroutines
	// (有缓冲管道，用于读、写两个goroutine之间的消息通信)
	NsqDataBuffChan   chan *NsqData
	MaxNsqDataChanLen uint32
	//发布完管道内所有数据
	wg sync.WaitGroup

	// Channel to notify that the connection has exited/stopped
	// (告知nsq退出/停止的channel)
	ctx    context.Context
	cancel context.CancelFunc

	//info
	channel         string
	nsqLookupAddr   string
	concurrency     int
	maxInFlight     int
	startWriterFlag int32
	dataPack        SDataPack

	Apis map[int32]SRouter
}

func NewNsqByConf(nsq2 sconfig.Nsq, dataPack SDataPack) (*Nsq, error) {
	//taskHandler := NewTaskHandler(nsq2.WorkerPoolSize, nsq2.MaxTaskChanLen)
	n := &Nsq{
		//BaseConnection: BaseConnection{
		//	TaskHandler: taskHandler,
		//},
		Apis: make(map[int32]SRouter),
		//taskHandler:       taskHandler,
		startWriterFlag:   0,
		producers:         make([]*NsqProducer, 0),
		Consumers:         make([]*NsqConsumer, 0),
		MaxNsqDataChanLen: nsq2.MaxNsqDataChanLen,
		channel:           nsq2.Channel,
		nsqLookupAddr:     nsq2.NsqlookupdAddr,
		concurrency:       nsq2.Concurrency,
		maxInFlight:       nsq2.MaxInFlight,
		dataPack:          dataPack,
	}
	for i, addr := range nsq2.NsqdAddrList {
		if p, err := NewProducer(addr); err != nil {
			return nil, err
		} else {
			n.producers = append(n.producers, p)
			slog.Ins().Infof("[nsq] add producer [%d]", i)
		}
	}
	n.ctx, n.cancel = context.WithCancel(context.Background())
	return n, nil
}

func NewNsq(workPoolSize, maxTaskQueueLen, maxNsqDataChanLen uint32, channel, nsqLookupAddr string, concurrency, maxInFlight int, nsqdList []string) *Nsq {
	//taskHandler := NewTaskHandler(workPoolSize, maxTaskQueueLen)
	n := &Nsq{
		//BaseConnection: BaseConnection{
		//	TaskHandler: taskHandler,
		//},
		//taskHandler:       taskHandler,
		startWriterFlag:   0,
		producers:         make([]*NsqProducer, 0),
		Consumers:         make([]*NsqConsumer, 0),
		MaxNsqDataChanLen: maxNsqDataChanLen,
		channel:           channel,
		nsqLookupAddr:     nsqLookupAddr,
		concurrency:       concurrency,
		maxInFlight:       maxInFlight,
	}
	for i, addr := range nsqdList {
		if p, err := NewProducer(addr); err != nil {
			panic(err)
		} else {
			n.producers = append(n.producers, p)
			slog.Ins().Infof("[nsq] add producer [%d]", i)
		}
	}
	n.ctx, n.cancel = context.WithCancel(context.Background())
	return n
}

func (n *Nsq) setStartWriterFlag() bool {
	return atomic.CompareAndSwapInt32(&n.startWriterFlag, 0, 1)
}

func (n *Nsq) addRouter(msgID int32, router SRouter) {
	// 1. Check whether the current API processing method bound to the msgID already exists
	// (判断当前msg绑定的API处理方法是否已经存在)
	if _, ok := n.Apis[msgID]; ok {
		msgErr := fmt.Sprintf("repeated api , msgID = %+v\n", msgID)
		panic(msgErr)
	}
	// 2. Add the binding relationship between msg and API
	// (添加msg与api的绑定关系)
	n.Apis[msgID] = router
	slog.Ins().Infof("Add Router msgID = %d", msgID)
}

func (n *Nsq) AddRouter(topic string, msgID int32, router SRouter) {
	//n.taskHandler.AddRouter(msgID, router)
	n.addRouter(msgID, router)
	if v, ok := TopicEnum_value[topic]; ok {
		slog.Ins().Infof("already created Consumer,Topic=%s msgId=%d,repeatedMsgId=%d", topic, msgID, v)
	} else {
		if c, err := NewNsqConsumer(topic, n.channel, n.nsqLookupAddr, n.concurrency, n.maxInFlight); err != nil {
			panic(err)
		} else {
			n.Consumers = append(n.Consumers, c)
			slog.Ins().Infof("created Consumer Success,Topic=%s, msgId=%d", topic, msgID)
		}
	}
}

func (n *Nsq) SendToMsgBuffChan(topic string, data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("data is nil")
	}
	if len(topic) == 0 {
		return fmt.Errorf("topic is nil")
	}

	nsqData := GetNsqData(topic, data)

	if n.NsqDataBuffChan == nil && n.setStartWriterFlag() {
		n.NsqDataBuffChan = make(chan *NsqData, n.MaxNsqDataChanLen)
		// Start a Goroutine to write data back to the client
		// This method only reads data from the MsgBuffChan without allocating memory or starting a Goroutine
		// (开启用于写回客户端数据流程的Goroutine
		// 此方法只读取MsgBuffChan中的数据没调用SendBuffMsg可以分配内存和启用协程)
		for _, producer := range n.producers {
			n.wg.Add(1)
			go n.StartWriter(producer)
		}
	}
	idleTimeout := time.NewTimer(5 * time.Millisecond)
	defer idleTimeout.Stop()
	// 要让数据发出去，先停止rpcx服务，释放rpcx端口，再停止nsq服务，等管道内所有消息发出再关闭本服务
	//if n.isClosed() == true {
	//	return errors.New("nsqd closed when send buff msg")
	//}

	// Send timeout
	select {
	case <-idleTimeout.C:
		return errors.New("send buff msg timeout")
	case n.NsqDataBuffChan <- nsqData:
		return nil
	}

}

func (n *Nsq) HandleMessage(message *nsq.Message) error {
	defer func() {
		if err := recover(); err != nil {
			slog.Ins().Errorf("Nsq HandleMessage error: %v", err)
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			slog.Ins().Errorf("panic in HandleMessage: %v, stack: %s", err, errStack[:n])
		}
	}()
	// 支持自定义dataPack
	dataPack := n.dataPack
	if dataPack == nil {
		dataPack = NsqDataPackObj
	}
	if msg, err := dataPack.Unpack(message.Body); err != nil {
		slog.Ins().Error("Nsq Consumer Unpack Data err", zap.Error(err))
		return nil
	} else {
		task := GetTask(nil, msg)
		defer PutTask(task)
		task.GetMessage().SetNsqMessage(message)

		msgId := task.GetMsgID()
		handler, ok := n.Apis[msgId]
		//n.taskHandler.SendTaskToTaskQueue(task)
		if !ok {
			slog.Ins().Errorf("api msgID = %d is not FOUND!", task.GetMsgID())
			// 返回报错，让其他版本的服务接收数据再试试
			return fmt.Errorf("api msgID = %d is not FOUND!", task.GetMsgID())
		}

		// Bind the Task request to the corresponding Router relationship
		// (Request请求绑定Router对应关系)
		task.BindRouter(handler)

		// Execute the corresponding processing method
		err = task.Call()
		if err != nil {
			slog.Ins().Error("task.Call error", zap.Error(err), zap.Int32("msgId", msgId))
			// 返回报错，让其他版本的服务接收数据再试试
			return err
		}
	}

	return nil
}

func (n *Nsq) StartWriter(p *NsqProducer) {
	slog.Ins().Infof("Nsq Writer Goroutine is running")
	defer slog.Ins().Infof("[Nsq Writer exit!]")
	defer n.wg.Done()
	for {
		select {
		case nsqData, ok := <-n.NsqDataBuffChan:
			if ok {
				if err := p.PublishDirect(nsqData.Topic, nsqData.data); err != nil {
					slog.Ins().Errorf("Send Buff Data error:, %s NsqProducer Publish error", err)
					// 失败的消息丢回管道 重新发
					if err = n.SendToMsgBuffChan(nsqData.Topic, nsqData.data); err != nil {
						slog.Ins().Errorf("SendToMsgBuffChan error:%s,", err.Error())
					}
					PutNsqData(nsqData)
					break
				}
				PutNsqData(nsqData)
			} else {
				slog.Ins().Errorf("msgBuffChan is Closed")
				break
			}
		case <-n.ctx.Done():
			l := len(n.NsqDataBuffChan)
			slog.Ins().Infof("[Nsq Writer exit! ctx.Done],NsqDataBuffChanLen:%d", l)
			return
		}
	}
}

func (n *Nsq) Start() {
	defer func() {
		if err := recover(); err != nil {
			var errStack = make([]byte, 1024)
			n := runtime.Stack(errStack, true)
			slog.Ins().Errorf("panic in Nsq Start: %v, stack: %s", err, errStack[:n])
		}
	}()
	//n.ctx, n.cancel = context.WithCancel(context.Background())  New的时候就创建
	// 开启taskWorkPool
	//n.taskHandler.StartWorkerPool()
	// 启动nsq consumer
	for i, consumer := range n.Consumers {
		err := consumer.StartReader(n)
		if err != nil {
			panic(err)
		} else {
			slog.Ins().Infof("[nsq] consumer.StartReader [%d]", i)
		}
	}

	select {
	case <-n.ctx.Done():
		// 停止所有消费
		for _, consumer := range n.Consumers {
			consumer.Stop()
		}
		// 让taskHandler停止
		//n.taskHandler.Stop()
		return
	}
}

// Stop stops the connection and ends the current connection state.
// (停止连接，结束当前连接状态)
func (n *Nsq) Stop() {
	// 1.立即停止所有消费者，没有新的消息进来
	// 2. 通知清空taskChan内的task后 关闭所有工作池
	n.cancel()
	// 等管道内所有数据发布完 ，结束所有Writer
	n.wg.Wait()
	// 让所有producer停止工作
	for i, producer := range n.producers {
		producer.producer.Stop()
		slog.Ins().Infof("[nsq] producer.Stop [%d]", i)
	}

}
