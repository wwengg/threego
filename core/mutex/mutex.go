// @Title
// @Description
// @Author  Wangwengang  2024/1/8 13:09
// @Update  Wangwengang  2024/1/8 13:09
package mutex

import (
	"context"
	"sync"

	"github.com/wwengg/threego/core/setcd"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdMutex struct {
	ctx     context.Context
	cancel  context.CancelFunc
	session *setcd.Session
	//mutex   *setcd.Mutex
	mutexMap sync.Map
}

// NewEtcdMutex default 10s
func NewEtcdMutex(client *clientv3.Client) (*EtcdMutex, error) {
	session, err := setcd.NewSession(client)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &EtcdMutex{
		ctx:     ctx,
		cancel:  cancel,
		session: session,
		//mutex:   setcd.NewMutex(session, key),
	}, nil
}

// NewEtcdLeaseMutex ttl
func NewEtcdLeaseMutex(key string, client *clientv3.Client, ttl int64) (*EtcdMutex, error) {
	res, err := client.Grant(context.Background(), ttl)
	if err != nil {
		return nil, err
	}

	session, err := setcd.NewSession(client, setcd.WithLease(res.ID))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &EtcdMutex{
		ctx:     ctx,
		cancel:  cancel,
		session: session,
		//mutex:   setcd.NewMutex(session, key),
	}, nil
}

// Lock get lock
func (em *EtcdMutex) Lock(key string) error {
	if v, ok := em.mutexMap.Load(key); ok {
		return v.(*setcd.Mutex).Lock(em.ctx)
	} else {
		mutex := setcd.NewMutex(em.session, key)
		em.mutexMap.Store(key, mutex)
		return mutex.Lock(em.ctx)
	}
}

func (em *EtcdMutex) Unlock(key string) error {
	if v, ok := em.mutexMap.Load(key); ok {
		return v.(*setcd.Mutex).Unlock(em.ctx)
	} else {
		mutex := setcd.NewMutex(em.session, key)
		em.mutexMap.Store(key, mutex)
		return mutex.Unlock(em.ctx)
	}
}

func (em *EtcdMutex) Close() {
	em.cancel()
	em.session.Close()
}
