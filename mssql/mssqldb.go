package mssql

import (
	"database/sql"
	"errors"

	"github.com/yangtizi/log/zaplog"
	"github.com/yangtizi/sql/scanner"
)

// TMsSQLDB 单个的数据库
type TMsSQLDB struct {
	chpool     chan int
	strConnect string
	pDB        *sql.DB
}

// NewDB 创建新的MSSQL数据库类
func NewDB(strReadConnect string) *TMsSQLDB {
	p := &TMsSQLDB{}
	p.init(strReadConnect)
	return p
}

func (m *TMsSQLDB) init(strConnect string) {
	db, err := sql.Open("mssql", strConnect)
	if err == nil {
		m.pDB = db
		m.strConnect = strConnect
		m.chpool = make(chan int, 30)
		return
	}

	zaplog.Ins.Errorf("数据库连接出现问题 connect = [%s] err = [%v]", strConnect, err)
}

func (m *TMsSQLDB) queryRow(strQuery string, args ...interface{}) (*scanner.TRow, error) {
	if m.pDB == nil {
		zaplog.Ins.Errorf("queryRow [db == nil]")
		return nil, errors.New("不存在DB")
	}

	m.chpool <- 1
	row := m.pDB.QueryRow(strQuery, args...)
	<-m.chpool

	return scanner.NewRow(row), nil
}

func (m *TMsSQLDB) queryRows(strQuery string, args ...interface{}) (*scanner.TRows, error) {
	if m.pDB == nil {
		zaplog.Ins.Errorf("queryRows [db == nil]")
		return nil, errors.New("不存在DB")
	}
	m.chpool <- 1
	rows, err := m.pDB.Query(strQuery, args...)
	<-m.chpool
	return scanner.NewRows(rows), err
}

func (m *TMsSQLDB) exec(strQuery string, args ...interface{}) (*scanner.TResult, error) {
	if m.pDB == nil {
		zaplog.Ins.Errorf("exec [db == nil]")
		return nil, errors.New("不存在DB")
	}

	m.chpool <- 1
	rs, err := m.pDB.Exec(strQuery, args...)
	<-m.chpool

	if err != nil {
		return nil, err
	}

	nInsert, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}

	nCount, err := rs.RowsAffected()

	return scanner.NewResult(nInsert, nCount), err
}

func (m *TMsSQLDB) transaction() (*sql.Tx, error) {
	if m.pDB == nil {
		zaplog.Ins.Errorf("transaction")
		return nil, errors.New("不存在DB")
	}

	m.chpool <- 1
	tx, err := m.pDB.Begin()
	<-m.chpool
	return tx, err
}
