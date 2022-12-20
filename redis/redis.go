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

func (m *TRedisDB) do(args ...any) *r.Cmd {
	return m.redisClient.Do(context.Background(), args...)
}

func (m *TRedisDB) client() *r.Client {
	return m.redisClient
}

/*
{
			url: "redis://localhost:123/1",
			o:   &Options{Addr: "localhost:123", DB: 1},
		}, {
			url: "redis://localhost:123",
			o:   &Options{Addr: "localhost:123"},
		}, {
			url: "redis://localhost/1",
			o:   &Options{Addr: "localhost:6379", DB: 1},
		}, {
			url: "redis://12345",
			o:   &Options{Addr: "12345:6379"},
		}, {
			url: "rediss://localhost:123",
			o:   &Options{Addr: "localhost:123", TLSConfig: &tls.Config{  }},
		}, {
			url: "redis://:bar@localhost:123",
			o:   &Options{Addr: "localhost:123", Password: "bar"},
		}, {
			url: "redis://foo@localhost:123",
			o:   &Options{Addr: "localhost:123", Username: "foo"},
		}, {
			url: "redis://foo:bar@localhost:123",
			o:   &Options{Addr: "localhost:123", Username: "foo", Password: "bar"},
		}, {
			// multiple params
			url: "redis://localhost:123/?db=2&read_timeout=2&pool_fifo=true",
			o:   &Options{Addr: "localhost:123", DB: 2, ReadTimeout: 2 * time.Second, PoolFIFO: true},
		}, {
			// special case handling for disabled timeouts
			url: "redis://localhost:123/?db=2&idle_timeout=0",
			o:   &Options{Addr: "localhost:123", DB: 2, IdleTimeout: -1},
		}, {
			// negative values disable timeouts as well
			url: "redis://localhost:123/?db=2&idle_timeout=-1",
			o:   &Options{Addr: "localhost:123", DB: 2, IdleTimeout: -1},
		}, {
			// absent timeout values will use defaults
			url: "redis://localhost:123/?db=2&idle_timeout=",
			o:   &Options{Addr: "localhost:123", DB: 2, IdleTimeout: 0},
		}, {
			url: "redis://localhost:123/?db=2&idle_timeout", // missing "=" at the end
			o:   &Options{Addr: "localhost:123", DB: 2, IdleTimeout: 0},
		}, {
			url: "unix:///tmp/redis.sock",
			o:   &Options{Addr: "/tmp/redis.sock"},
		}, {
			url: "unix://foo:bar@/tmp/redis.sock",
			o:   &Options{Addr: "/tmp/redis.sock", Username: "foo", Password: "bar"},
		}, {
			url: "unix://foo:bar@/tmp/redis.sock?db=3",
			o:   &Options{Addr: "/tmp/redis.sock", Username: "foo", Password: "bar", DB: 3},
		}, {
			// invalid db format
			url: "unix://foo:bar@/tmp/redis.sock?db=test",
			err: errors.New(`redis: invalid database number: strconv.Atoi: parsing "test": invalid syntax`),
		}, {
			// invalid int value
			url: "redis://localhost/?pool_size=five",
			err: errors.New(`redis: invalid pool_size number: strconv.Atoi: parsing "five": invalid syntax`),
		}, {
			// invalid bool value
			url: "redis://localhost/?pool_fifo=yes",
			err: errors.New(`redis: invalid pool_fifo boolean: expected true/false/1/0 or an empty string, got "yes"`),
		}, {
			// it returns first error
			url: "redis://localhost/?db=foo&pool_size=five",
			err: errors.New(`redis: invalid database number: strconv.Atoi: parsing "foo": invalid syntax`),
		}, {
			url: "redis://localhost/?abc=123",
			err: errors.New("redis: unexpected option: abc"),
		}, {
			url: "redis://foo@localhost/?username=bar",
			err: errors.New("redis: unexpected option: username"),
		}, {
			url: "redis://localhost/?wrte_timout=10s&abc=123",
			err: errors.New("redis: unexpected option: abc, wrte_timout"),
		}, {
			url: "http://google.com",
			err: errors.New("redis: invalid URL scheme: http"),
		}, {
			url: "redis://localhost/1/2/3/4",
			err: errors.New("redis: invalid URL path: /1/2/3/4"),
		}, {
			url: "12345",
			err: errors.New("redis: invalid URL scheme: "),
		}, {
			url: "redis://localhost/iamadatabase",
			err: errors.New(`redis: invalid database number: "iamadatabase"`),
		},
*/
