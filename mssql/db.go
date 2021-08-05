package mssql

import (
	"database/sql"
	"errors"
	"sync"

	"github.com/yangtizi/log/zaplog"
	"github.com/yangtizi/sql/scanner"

	_ "github.com/denisenkom/go-mssqldb" // mssql 数据库
)

var mapMSSQL sync.Map

// QueryRow (agent 代理商编号, strQuery sql脚本, args 脚本参数)
func QueryRow(agent interface{}, strQuery string, args ...interface{}) (*scanner.TRow, error) {
	zaplog.Debugf("agent = [%v], strQuery = [%s]", agent, strQuery)
	zaplog.Debug("[+] ", args)
	v, ok := mapMSSQL.Load(agent)
	if !ok {
		zaplog.Errorf("QueryRow 不存在索引")
		return nil, errors.New("不存在的DB索引")
	}

	return v.(*TMsSQLDB).queryRow(strQuery, args...)
}

// QueryRows (agent 代理商编号, strQuery sql脚本, args 脚本参数)
func QueryRows(agent interface{}, strQuery string, args ...interface{}) (*scanner.TRows, error) {
	zaplog.Debugf("agent = [%v], strQuery = [%s]", agent, strQuery)
	zaplog.Debug("[+] ", args)
	v, ok := mapMSSQL.Load(agent)
	if !ok {
		zaplog.Errorf("QueryRows 不存在索引")
		return nil, errors.New("不存在的DB索引")
	}

	return v.(*TMsSQLDB).queryRows(strQuery, args...)
}

// Exec (agent 代理商编号, strQuery sql脚本, args 脚本参数)
func Exec(agent interface{}, strQuery string, args ...interface{}) (*scanner.TResult, error) {
	zaplog.Debugf("agent = [%v], strQuery = [%s]", agent, strQuery)
	zaplog.Debug("[+] ", args)
	v, ok := mapMSSQL.Load(agent)
	if !ok {
		zaplog.Errorf("Exec 不存在索引")
		return nil, errors.New("不存在的DB索引")
	}

	return v.(*TMsSQLDB).exec(strQuery, args...)
}

// Transaction 事务
func Transaction(agent interface{}) (*sql.Tx, error) {
	zaplog.Debugf("Transaction begin")
	v, ok := mapMSSQL.Load(agent)
	if !ok {
		zaplog.Errorf("Transaction 不存在索引")
		return nil, errors.New("不存在的DB索引")
	}

	return v.(*TMsSQLDB).transaction()
}

// InitDB 初始化DB (agent 代理商编号, strConnect 从库连接字符串)
func InitDB(agent interface{}, strConnect string) {
	_, ok := mapMSSQL.Load(agent)
	if !ok {
		// * 创建新的DB指针
		pMsSQL := NewDB(strConnect)

		zaplog.Infof("正在连接数据库 agent = [%v], strConnect = [%s]", agent, strConnect)
		mapMSSQL.Store(agent, pMsSQL)
		return
	}

	// zaplog.Println("已经存在确有重复创建")
}
