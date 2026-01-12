package sbus

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/wwengg/threego/core/slog"
	"go.uber.org/zap"
)

type MsgData struct {
	MsgID uint32
	data  []byte
}

var MsgDataPool = new(sync.Pool)

func init() {
	MsgDataPool.New = func() interface{} {
		return allocateMsgData()
	}
}

func allocateMsgData() *MsgData {
	msgData := new(MsgData)
	return msgData
}
func (md *MsgData) Reset(msgId uint32, data []byte) {
	md.MsgID = msgId
	md.data = data
}

func GetMsgData(msgId uint32, data []byte) *MsgData {

	// 根据当前模式判断是否使用对象池

	// 从对象池中取得一个 Request 对象,如果池子中没有可用的 Request 对象则会调用 allocateRequest 函数构造一个新的对象分配
	r := MsgDataPool.Get().(*MsgData)
	// 因为取出的 Request 对象可能是已存在也可能是新构造的,无论是哪种情况都应该初始化再返回使用
	r.Reset(msgId, data)
	return r
}

func PutMsgData(msgData *MsgData) {
	MsgDataPool.Put(msgData)
}

type SConnection interface {
	// Start the connection, make the current connection start working
	// (启动连接，让当前连接开始工作)
	Start()
	// Stop the connection and end the current connection state
	// (停止连接，结束当前连接状态)
	Stop()

	// Returns ctx, used by user-defined go routines to obtain connection exit status
	// (返回ctx，用于用户自定义的go程获取连接退出状态)
	Context() context.Context

	GetConnID() uint64            // Get the current connection ID (获取当前连接ID)
	GetConnIdStr() string         // Get the current connection ID for string (获取当前字符串连接ID)
	GetTaskHandler() STaskHandler // Get the message handler (获取消息处理器)
	RemoteAddr() net.Addr         // Get the remote address information of the connection (获取链接远程地址信息)
	LocalAddr() net.Addr          // Get the local address information of the connection (获取链接本地地址信息)
	LocalAddrString() string      // Get the local address information of the connection as a string
	RemoteAddrString() string     // Get the remote address information of the connection as a string
	GetConnVersion() int32

	SendData(data []byte) error // Send data to the message queue to be sent to the remote TCP client later
	SendMsg(msg SMsg) error

	SetProperty(key string, value string)   // Set connection property
	GetProperty(key string) (string, error) // Get connection property
	RemoveProperty(key string)              // Remove connection property
	IsAlive() bool                          // Check if the current connection is alive(判断当前连接是否存活)
	SetHeartBeat(checker SHeartbeatChecker) // Set the heartbeat detector (设置心跳检测器)

	// 返回当前连接是否存在FrameDecoder
	HasFrameDecoder() bool

	//AddCloseCallback(handler, key interface{}, callback func()) // Add a close callback function (添加关闭回调函数)
	//RemoveCloseCallback(handler, key interface{})               // Remove a close callback function (删除关闭回调函数)
	//InvokeCloseCallbacks()                                      // Trigger the close callback function (触发关闭回调函数，独立协程完成)
}

type Connection struct {
	Conn net.Conn
	// The ID of the current connection, also known as SessionID, globally unique, used by server Connection
	// uint64 range: 0~18,446,744,073,709,551,615
	// This is the maximum number of connID theoretically supported by the process
	// (当前连接的ID 也可以称作为SessionID，ID全局唯一 ，服务端Connection使用
	// uint64 取值范围：0 ~ 18,446,744,073,709,551,615
	// 这个是理论支持的进程connID的最大数量)
	ConnID uint64
	// connection id for string
	// (字符串的连接id)
	ConnIdStr string
	// 连接版本
	ConnVersion int32
	// The message management module that manages MsgID and the corresponding processing method
	// (消息管理MsgID和对应处理方法的消息管理模块)
	TaskHandler STaskHandler
	// onConnStart is the Hook function when the current connection is created.
	// (当前连接创建时Hook函数)
	OnConnStart func(conn SConnection)
	// onConnStop is the Hook function when the current connection is created.
	// (当前连接断开时的Hook函数)
	OnConnStop func(conn SConnection)
	// ctx and cancel are used to notify that the connection has exited/stopped.
	// (告知该链接已经退出/停止的channel)
	ctx    context.Context
	cancel context.CancelFunc
	// Which Connection Manager the current connection belongs to
	// (当前链接是属于哪个Connection Manager的)
	connManager SConnManager

	// frameDecoder is the decoder for splitting or splicing data packets.
	// (断粘包解码器)
	FrameDecoder SFrameDecoder

	Datapack SDataPack

	// msgLock is used for locking when users send and receive messages.
	// (用户收发消息的Lock)
	msgLock sync.RWMutex

	// property is the connection attribute. (链接属性)
	Property map[string]string

	// propertyLock protects the current property lock. (保护当前property的锁)
	propertyLock sync.Mutex

	IOReadBuffSize uint32

	// Last activity time
	// (最后一次活动时间)
	lastActivityTime time.Time

	hc SHeartbeatChecker

	heartBeatDuration time.Duration
}

func NewConnection(conn net.Conn, connId uint64, connVersion int32, taskHandler STaskHandler, OnConnStart, OnConnStop func(conn SConnection), frameDecoder SFrameDecoder, datapack SDataPack, connManager SConnManager, IOReadBuffSize uint32, heartbeatDuration time.Duration) SConnection {
	return &Connection{
		Conn:              conn,
		ConnID:            connId,
		ConnIdStr:         fmt.Sprintf("%d", connId),
		ConnVersion:       connVersion,
		TaskHandler:       taskHandler,
		OnConnStart:       OnConnStart,
		OnConnStop:        OnConnStop,
		FrameDecoder:      frameDecoder,
		Datapack:          datapack,
		Property:          nil,
		IOReadBuffSize:    IOReadBuffSize,
		connManager:       connManager,
		heartBeatDuration: heartbeatDuration,
	}
}

// 更新心跳检测时间
func (c *Connection) updateActivity() {
	c.lastActivityTime = time.Now()
}

func (bc *Connection) callOnConnStart() {
	if bc.OnConnStart != nil {
		slog.Ins().Info("CallOnConnStart....")
		bc.OnConnStart(bc)
	}
}

func (bc *Connection) callOnConnStop() {
	if bc.OnConnStop != nil {
		slog.Ins().Info("callOnConnStop....")
		bc.OnConnStop(bc)
	}
}

func (bc *Connection) isClosed() bool {
	return bc.ctx == nil || bc.ctx.Err() != nil
}

func (bc *Connection) StartReader() {
	slog.Ins().Infof("[Reader Goroutine is running]")
	defer slog.Ins().Infof("%s [conn Reader exit!]", bc.ConnIdStr)
	defer bc.Stop()
	defer func() {
		if err := recover(); err != nil {
			slog.Ins().Errorf("Reader connID=%d, panic err=%v", bc.GetConnID(), err)
		}
	}()
	//Reduce buffer allocation times to improve efficiency
	// add by ray 2023-02-03
	buffer := make([]byte, bc.IOReadBuffSize)

	for {
		select {
		case <-bc.ctx.Done():
			// 停止循环 不读了，连接断开啦！！
			return
		default:
			if n, err := bc.Conn.Read(buffer); err != nil {
				slog.Ins().Errorf("read msg head [read datalen=%d], error = %s", n, err)
				return
			} else {
				if n == 0 {
					continue
				}
				if n > 0 && bc.hc != nil {
					bc.updateActivity()
				}
				// Deal with the custom protocol fragmentation problem, added by uuxia 2023-03-21
				// (处理自定义协议断粘包问题)
				if bc.FrameDecoder != nil {
					// Decode the 0-n bytes of data read
					// (为读取到的0-n个字节的数据进行解码)
					bufArrays, err2 := bc.FrameDecoder.Decode(buffer[0:n])
					if bufArrays == nil {
						continue
					}
					for _, bytes := range bufArrays {
						msg, err := bc.Datapack.Unpack(bytes)
						if err != nil {
							slog.Ins().Error(err.Error())
							continue
						}
						// Get the current client's Request data
						// (得到当前客户端请求的Request数据)
						task := GetTask(bc, msg)
						// 如果cmd为心跳包，不走后续逻辑，直接心跳保活 发送心跳包给客户端
						if task.GetCmd() == bc.hc.Cmd() {
							err := bc.hc.SendHeartBeatMsg()
							if err != nil {
								slog.Ins().Error("SendHeartBeatMsg", zap.Error(err))
								return
							}
							continue
						}
						bc.TaskHandler.SendTaskToTaskQueue(task)
					}
					if err2 != nil {
						slog.Ins().Error(err2.Error())
						continue // 发送过长数据包或协议错误
					}
				} else {
					msg, err := bc.Datapack.Unpack(buffer[0:n])
					if err != nil {
						slog.Ins().Error(err.Error())
						continue
					}
					// Get the current client's Request data
					// (得到当前客户端请求的Request数据)
					task := GetTask(bc, msg)
					// 如果cmd为心跳包，不走后续逻辑，直接心跳保活 发送心跳包给客户端
					if task.GetCmd() == bc.hc.Cmd() {
						err := bc.hc.SendHeartBeatMsg()
						if err != nil {
							slog.Ins().Error("SendHeartBeatMsg", zap.Error(err))
							return
						}
						continue
					}
					bc.TaskHandler.SendTaskToTaskQueue(task)
				}
			}

		}
	}
}

// Start()
func (bc *Connection) Start() {
	bc.ctx, bc.cancel = context.WithCancel(context.Background())

	bc.callOnConnStart()
	// Start heartbeating detection
	if bc.hc != nil {
		bc.hc.Start()
		bc.updateActivity()
	}

	// Start the Goroutine for reading data from the client
	// (开启用户从客户端读取数据流程的Goroutine)
	go bc.StartReader()

	select {
	case <-bc.ctx.Done():
		// If the user has registered a close callback for the connection, it should be called explicitly at this moment.
		// (如果用户注册了该链接的	关闭回调业务，那么在此刻应该显示调用)
		bc.callOnConnStop()

		if bc.hc != nil {
			bc.hc.Stop()
		}

		_ = bc.Conn.Close()
		if bc.connManager != nil {
			bc.connManager.Remove(bc)
		}
		slog.Ins().Debugf("Conn Stop() ...ConnID = %d", bc.ConnID)
		return
	}
}
func (bc *Connection) Stop() {
	bc.cancel()
}
func (bc *Connection) Context() context.Context {
	return bc.ctx
}
func (bc *Connection) GetConnID() uint64 {
	return bc.ConnID
}
func (bc *Connection) GetConnIdStr() string {
	return bc.ConnIdStr
}
func (bc *Connection) GetTaskHandler() STaskHandler {
	return bc.TaskHandler
}
func (bc *Connection) RemoteAddr() net.Addr     { return bc.Conn.RemoteAddr() }
func (bc *Connection) LocalAddr() net.Addr      { return bc.Conn.LocalAddr() }
func (bc *Connection) LocalAddrString() string  { return bc.Conn.LocalAddr().String() }
func (bc *Connection) RemoteAddrString() string { return bc.Conn.RemoteAddr().String() }
func (bc *Connection) GetConnVersion() int32    { return bc.ConnVersion }
func (bc *Connection) HasFrameDecoder() bool    { return bc.FrameDecoder != nil }
func (bc *Connection) SendData(data []byte) error {
	bc.msgLock.RLock()
	defer bc.msgLock.RUnlock()
	defer func() {
		if err := recover(); err != nil {
			slog.Ins().Errorf("SendData connID=%d, panic err=%v", bc.GetConnID(), err)
		}
	}()
	if bc.isClosed() == true {
		return errors.New("Connection closed when send Data")
	}
	_, err := bc.Conn.Write(data)
	if err != nil {
		slog.Ins().Errorf("SendMsg err data = %+v, err = %+v", data, err)
		return err
	} else {
		slog.Ins().Debug("SendMsg data success")
	}
	return nil
}
func (bc *Connection) SendMsg(msg SMsg) error {
	msg.SetHasFrameDecoder(bc.HasFrameDecoder()) // 判断该连接是否需要编码器，pack的时候就可以选择性pack
	if data, err := bc.Datapack.Pack(msg); err == nil {
		slog.Ins().Debug("pack", zap.Any("msg", msg))
		return bc.SendData(data)
	} else {
		return err
	}
}
func (bc *Connection) SetProperty(key string, value string) {
	bc.propertyLock.Lock()
	defer bc.propertyLock.Unlock()
	if bc.Property == nil {
		bc.Property = make(map[string]string)
	}

	bc.Property[key] = value
}
func (bc *Connection) GetProperty(key string) (string, error) {
	bc.propertyLock.Lock()
	defer bc.propertyLock.Unlock()

	if value, ok := bc.Property[key]; ok {
		return value, nil
	}

	return "", errors.New("no property found")
}
func (bc *Connection) RemoveProperty(key string) {
	bc.propertyLock.Lock()
	defer bc.propertyLock.Unlock()

	delete(bc.Property, key)
}
func (bc *Connection) IsAlive() bool {
	if bc.isClosed() {
		return false
	}
	// Check the last activity time of the connection. If it's beyond the heartbeat interval,
	// then the connection is considered dead.
	// (检查连接最后一次活动时间，如果超过心跳间隔，则认为连接已经死亡)
	return time.Now().Sub(bc.lastActivityTime) < bc.heartBeatDuration
}
func (bc *Connection) SetHeartBeat(checker SHeartbeatChecker) {
	bc.hc = checker
}

// func (bc *BaseConnection) AddCloseCallback(handler, key interface{}, callback func()) {}
// func (bc *BaseConnection) RemoveCloseCallback(handler, key interface{})               {}
// func (bc *BaseConnection) InvokeCloseCallbacks()
