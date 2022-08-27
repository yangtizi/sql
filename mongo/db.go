package mongo

import (
	"errors"
	"sync"

	_ "github.com/denisenkom/go-mssqldb" // mssql 数据库
	"github.com/qiniu/qmgo"
	"github.com/yangtizi/log/zaplog"
)

var mapMongo sync.Map
var instance *TMongoDB

// GetTable 获取表格
// func GetTable(agent interface{}, strDatabase, strCollection string) (*mgo.Collection, error) {
// 	if agent == nil {
// 		if instance == nil {
// 			return nil, errors.New("不存在的DB索引")
// 		}
// 		return instance.pDB.DB(strDatabase).C(strCollection), nil
// 	}
// 	zaplog.Ins.Debugf("%v, %s, %s", agent, strDatabase, strCollection)
// 	v, ok := mapMongo.Load(agent)
// 	if !ok {
// 		zaplog.Ins.Errorf("Exec 不存在索引")
// 		return nil, errors.New("不存在的DB索引")
// 	}
// 	return v.(*TMongoDB).pDB.DB(strDatabase).C(strCollection), nil
// }

// GetDB 获取DB
func GetDB(agent interface{}, D string, C string) (*qmgo.Collection, error) {

	if agent == nil {
		if instance == nil {
			zaplog.Ins.Errorf("不存在的DB索引")
			return nil, errors.New("不存在的DB索引")
		}

		instance.pDB.Database(D).Collection(C)
	}

	v, ok := mapMongo.Load(agent)
	if !ok {
		zaplog.Ins.Errorf("Exec 不存在索引")
		return nil, errors.New("不存在的DB索引")
	}
	return v.(*TMongoDB).pDB.Database(D).Collection(C), nil
}

// InitDB 初始化DB (strAgent 代理商编号, strConnect 从库连接字符串)
func InitDB(agent interface{}, strConnect string) {
	if agent == nil {
		instance = NewDB(strConnect)
		zaplog.Ins.Infof("正在连接数据库 agent = [%v], strConnect = [%s]", "默认", strConnect)
		return
	}

	_, ok := mapMongo.Load(agent)

	if !ok {
		// * 创建新的DB指针
		pMongo := NewDB(strConnect)

		zaplog.Ins.Infof("正在连接数据库 agent = [%v], strConnect = [%s]", agent, strConnect)
		mapMongo.Store(agent, pMongo)
		return
	}

	zaplog.Ins.Println("已经存在确有重复创建")
}
