package etcd

import (
	"context"
	"errors"
	"time"

	"github.com/yangtizi/log/zaplog"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type TEtcdDB struct {
	pClient *clientv3.Client
}

func NewDB(endpoints []string) *TEtcdDB {
	p := &TEtcdDB{}
	p.init(endpoints)
	return p
}

func (m *TEtcdDB) init(endpoints []string) {
	pClient, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		zaplog.Ins.Errorf("etcd连接数据库失败 endpoints = [%v] err = [%v]",
			endpoints, err)
		return
	}
	m.pClient = pClient
}

func (m *TEtcdDB) Put(strKey string, strValue string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	if m.pClient == nil {
		zaplog.Ins.Errorf("Put [pClient == nil]")
		return nil, errors.New("不存在DB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()
	resp, err := m.pClient.Put(ctx, strKey, strValue, opts...)

	return resp, err
}
