package mongo

import (
	"context"

	"github.com/qiniu/qmgo"
	"github.com/yangtizi/log/zaplog"
)

// TMongoDB .
type TMongoDB struct {
	chpool     chan int
	strConnect string
	pDB        *qmgo.Client
}

// NewDB 新DB
func NewDB(strConnect string) *TMongoDB {
	p := &TMongoDB{}
	p.init(strConnect)
	return p
}

func (m *TMongoDB) init(strConnect string) {
	client, err := qmgo.NewClient(context.Background(), &qmgo.Config{Uri: strConnect})

	// mongo, err := mgo.Dial(strConnect)
	if err != nil {
		zaplog.Ins.Errorf("%v", err)
		return
	}

	m.pDB = client
	m.strConnect = strConnect
}
