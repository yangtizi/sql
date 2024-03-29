package mysql

import (
	"database/sql"
	"errors"
	"strings"
	"sync"

	"github.com/yangtizi/log/zaplog"
	"github.com/yangtizi/sql/scanner"
	"github.com/yangtizi/sql/tx"

	_ "github.com/go-sql-driver/mysql" // mysql 数据库
)

var mapMYSQL sync.Map

// QueryRow (strAgent 代理商编号, strQuery sql脚本, args 脚本参数)
func QueryRow(strAgent string, strQuery string, args ...interface{}) (*sql.Row, error) {
	strQuery2 := strings.ReplaceAll(strQuery, "\n", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\r", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\t", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")

	zaplog.Map("mysql").Debugf("strAgent = [%s], strQuery = [%s]", strAgent, strQuery2)
	zaplog.Map("mysql").Debug("[+] ", args)
	v, ok := mapMYSQL.Load(strAgent)
	if !ok {
		return nil, errors.New("不存在的DB索引")
	}

	return v.(*TMySQLDB).queryRow(strQuery, args...)
}

// QueryRows (strAgent 代理商编号, strQuery sql脚本, args 脚本参数)
func QueryRows(strAgent string, strQuery string, args ...interface{}) (*sql.Rows, error) {
	strQuery2 := strings.ReplaceAll(strQuery, "\n", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\r", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\t", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")

	zaplog.Map("mysql").Debugf("strAgent = [%s], strQuery = [%s]", strAgent, strQuery2)
	zaplog.Map("mysql").Debug("[+] ", args)
	v, ok := mapMYSQL.Load(strAgent)
	if !ok {
		return nil, errors.New("不存在的DB索引")
	}

	return v.(*TMySQLDB).queryRows(strQuery, args...)
}

// 高级Rows请求, 直接返回MAP
func QueryRowsEx(strAgent string, strQuery string, args ...interface{}) ([]map[string]string, error) {
	strQuery2 := strings.ReplaceAll(strQuery, "\n", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\r", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\t", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	zaplog.Map("mysql").Debugf("strAgent = [%s], strQuery = [%s]", strAgent, strQuery2)
	zaplog.Map("mysql").Debug("[+] ", args)
	v, ok := mapMYSQL.Load(strAgent)
	if !ok {
		return nil, errors.New("不存在的DB索引")
	}

	rows, err := v.(*TMySQLDB).queryRows(strQuery, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// 获取字段名的string切片
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	results := make([]map[string]string, 0)

	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		record := make(map[string]string)
		if err != nil {
			return nil, err
		}
		for i, col := range values {
			value := string(col)
			record[columns[i]] = value
		}
		results = append(results, record)
	}
	return results, nil

}

// Exec (strAgent 代理商编号, strQuery sql脚本, args 脚本参数)
func Exec(strAgent string, strQuery string, args ...interface{}) (*scanner.TResult, error) {
	strQuery2 := strings.ReplaceAll(strQuery, "\n", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\r", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "\t", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	strQuery2 = strings.ReplaceAll(strQuery2, "  ", " ")
	zaplog.Map("mysql").Debugf("strAgent = [%s], strQuery = [%s]", strAgent, strQuery2)
	zaplog.Map("mysql").Debug("[+] ", args)
	v, ok := mapMYSQL.Load(strAgent)
	if !ok {
		return nil, errors.New("不存在的DB索引")
	}

	return v.(*TMySQLDB).exec(strQuery, args...)
}

// InitDB 初始化DB (strAgent 代理商编号, strConnect 从库连接字符串)
func InitDB(strAgent string, strConnect string) {
	_, ok := mapMYSQL.Load(strAgent)
	if !ok {
		// * 创建新的DB指针
		pMsSQL := NewDB(strConnect)

		zaplog.Map("mysql").Infof("正在连接数据库, 代理编号=[%s], 连接字符串=[%s]", strAgent, strConnect)
		mapMYSQL.Store(strAgent, pMsSQL)
		return
	}

	zaplog.Map("mysql").Warnf("已经存在确有重复创建")
}
func BeginTX(strAgent string) (*tx.TTx, error) {
	v, ok := mapMYSQL.Load(strAgent)
	if !ok {
		return nil, errors.New("不存在的DB索引")
	}

	t, err := v.(*TMySQLDB).beginTX()
	return tx.NewTx(t), err
}
