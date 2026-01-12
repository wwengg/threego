package sbus

import (
	"errors"
	"strconv"

	"github.com/wwengg/threego/core/slog"
	"github.com/wwengg/threego/core/utils"
)

type ConnManager struct {
	Connections utils.ShardLockMaps
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		Connections: utils.NewShardLockMaps(),
	}
}

func (connMgr *ConnManager) Add(conn SConnection) {

	connMgr.Connections.Set(conn.GetConnIdStr(), conn) // 将conn连接添加到ConnManager中

	slog.Ins().Debugf("connection add to ConnManager successfully: conn num = %d", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn SConnection) {

	connMgr.Connections.Remove(conn.GetConnIdStr()) // 删除连接信息

	slog.Ins().Debugf("connection Remove ConnID=%d successfully: conn num = %d", conn.GetConnID(), connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint64) (SConnection, error) {

	strConnId := strconv.FormatUint(connID, 10)
	if conn, ok := connMgr.Connections.Get(strConnId); ok {
		return conn.(SConnection), nil
	}

	return nil, errors.New("connection not found")
}

// Get2 It is recommended to use this method to obtain connection instances
func (connMgr *ConnManager) Get2(strConnId string) (SConnection, error) {

	if conn, ok := connMgr.Connections.Get(strConnId); ok {
		return conn.(SConnection), nil
	}

	return nil, errors.New("connection not found")
}

func (connMgr *ConnManager) Len() int {

	length := connMgr.Connections.Count()

	return length
}

func (connMgr *ConnManager) ClearConn() {

	// Stop and delete all connection information
	for item := range connMgr.Connections.IterBuffered() {
		val := item.Val
		if conn, ok := val.(SConnection); ok {
			// stop will eventually trigger the deletion of the connection,
			// no additional deletion is required
			conn.Stop()
		}
	}

	slog.Ins().Infof("Clear All Connections successfully: conn num = %d", connMgr.Len())
}

func (connMgr *ConnManager) GetAllConnID() []uint64 {

	strConnIdList := connMgr.Connections.Keys()
	ids := make([]uint64, 0, len(strConnIdList))

	for _, strId := range strConnIdList {
		connId, err := strconv.ParseUint(strId, 10, 64)
		if err == nil {
			ids = append(ids, connId)
		} else {
			slog.Ins().Infof("GetAllConnID Id: %d, error: %v", connId, err)
		}
	}

	return ids
}

func (connMgr *ConnManager) GetAllConnIdStr() []string {
	return connMgr.Connections.Keys()
}

func (connMgr *ConnManager) Range(cb func(uint64, SConnection, interface{}) error, args interface{}) (err error) {

	connMgr.Connections.IterCb(func(key string, v interface{}) {
		conn, _ := v.(SConnection)
		connId, _ := strconv.ParseUint(key, 10, 64)
		err = cb(connId, conn, args)
		if err != nil {
			slog.Ins().Infof("Range key: %v, v: %v, error: %v", key, v, err)
		}
	})

	return err
}

// Range2 It is recommended to use this method to 'Range'
func (connMgr *ConnManager) Range2(cb func(string, SConnection, interface{}) error, args interface{}) (err error) {

	connMgr.Connections.IterCb(func(key string, v interface{}) {
		conn, _ := v.(SConnection)
		err = cb(conn.GetConnIdStr(), conn, args)
		if err != nil {
			slog.Ins().Infof("Range2 key: %v, v: %v, error: %v", key, v, err)
		}
	})

	return err
}
