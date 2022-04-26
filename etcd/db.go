package etcd

import (
	"errors"
	"sync"

	"github.com/yangtizi/log/zaplog"
	clientv3 "go.etcd.io/etcd/client/v3"
	// 成功
)

var mapETCD sync.Map
var instance *TEtcdDB

func InitDB(agent interface{}, endpoints []string) {
	if agent == nil {
		instance = NewDB(endpoints)

		return
	}

	_, ok := mapETCD.Load(agent)
	if !ok {

		pEtcd := NewDB(endpoints)
		zaplog.Ins.Infof("正在连接Etcd 数据库")

		mapETCD.Store(agent, pEtcd)
		return
	}
}

func Put(agent interface{}, strKey string, strValue string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if agent == nil {
		if instance == nil {
			return nil, errors.New("")
		}

		return instance.Put(strKey, strValue, opts...)
	}

	v, ok := mapETCD.Load(agent)
	if !ok {
		return nil, errors.New("")
	}
	return v.(*TEtcdDB).Put(strKey, strValue, opts...)
}
