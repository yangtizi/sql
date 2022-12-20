package redis

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	r "github.com/go-redis/redis/v8"
	"github.com/yangtizi/log/zaplog"
)

// TValues 快速解析
type TValues []interface{}

func toString(val interface{}) (string, error) {
	switch val := val.(type) {
	case string:
		return val, nil
	default:
		err := fmt.Errorf("redis: unexpected type=%T for String", val)
		return "", err
	}
}

func toInt64(val interface{}) (int64, error) {
	switch val := val.(type) {
	case int64:
		return val, nil
	case string:
		return strconv.ParseInt(val, 10, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Int64", val)
		return 0, err
	}
}

func toFloat64(val interface{}) (float64, error) {
	switch val := val.(type) {
	case int64:
		return float64(val), nil
	case string:
		return strconv.ParseFloat(val, 64)
	default:
		err := fmt.Errorf("redis: unexpected type=%T for Float64", val)
		return 0, err
	}
}

// I 快速获取整数
func (m TValues) I(n int) int64 {
	if n > len(m) {
		return 0
	}
	if m[n] == nil {
		return 0
	}
	i64, err := toInt64(m[n])
	if err != nil {
		return 0
	}

	return i64
}

// S 快速获取字符串
func (m TValues) S(n int) string {
	if n > len(m) {
		return ""
	}
	if m[n] == nil {
		return ""
	}
	s, err := toString(m[n])
	if err != nil {
		return ""
	}
	return s
}

var mapRedis sync.Map

// Do (strAgent 代理商编号, strCommand sql脚本, args 脚本参数)
func Do(strAgent string, args ...any) (*r.Cmd, error) {
	v, ok := mapRedis.Load(strAgent)
	if !ok {
		return nil, errors.New("不存在的DB索引")
	}

	cmd := v.(*TRedisDB).do(args...)
	return cmd, cmd.Err()
}

func Client(strAgent string) *r.Client {
	v, ok := mapRedis.Load(strAgent)
	if !ok {
		return nil
	}

	cli := v.(*TRedisDB).client()
	return cli
}

// // HMGet (strAgent 代理商编号, args 脚本参数)
// func HMGet(strAgent string, args ...interface{}) (TValues, error) {
// 	v, ok := mapRedis.Load(strAgent)
// 	if !ok {
// 		return nil, errors.New("不存在的DB索引")
// 	}

// 	return v.(*TRedisDB).do("hmget", args...)
// }

// func HMSet(strAgent string, args ...interface{}) (TValues, error) {
// 	v, ok := mapRedis.Load(strAgent)
// 	if !ok {
// 		return nil, errors.New("不存在的DB索引")
// 	}

// 	return v.(*TRedisDB).do("hmset", args...)
// }

// InitDB 初始化DB (strAgent 代理商编号, strReadConnect 从库连接字符串, strWriteConnect 主库连接字符串)
func InitDB(strAgent string, strConnect string) {
	zaplog.Map("redis").Debugf("strAgent = [%s], strConnect = [%s]", strAgent, strConnect)
	_, ok := mapRedis.Load(strAgent)
	if !ok {
		// * 创建新的DB指针
		pRedis := &TRedisDB{}
		pRedis.init(strConnect)
		mapRedis.Store(strAgent, pRedis)
		return
	}

	zaplog.Map("Redis").Warnf("已经存在确有重复创建")
}
