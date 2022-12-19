package redis

import (
	"context"

	r "github.com/go-redis/redis/v8"
)

// TRedisDB 单个的数据库
type TRedisDB struct {
	strConnect  string
	redisClient *r.Client
}

func (m *TRedisDB) init(strConnect string) error {
	opt, err := r.ParseURL(strConnect)
	if err != nil {
		return err
	}
	m.strConnect = strConnect
	rdb := r.NewClient(opt)
	m.redisClient = rdb
	return nil
}

func (m *TRedisDB) do(strCommand string, args ...interface{}) *r.Cmd {
	return m.redisClient.Do(context.Background(), args...)
}
