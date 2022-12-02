package table

import (
	"github.com/yangtizi/sql/mysql"
	"github.com/yangtizi/sql/scanner"
)

// import "github.com/yangtizi/sql/scanner"

type IBaseInfo interface {
	Params() string
	Dataname() string  // 数据库名称
	Tablename() string // 表名称
	ScanFrom(scanner.IScanner) error
}

// 查询
func Select(m IBaseInfo, strWhere string, params ...interface{}) error {
	row, err := mysql.QueryRow(m.Dataname(), "SELECT "+m.Params()+" FROM "+m.Tablename()+" "+strWhere, params...)
	if err != nil {
		return err
	}
	return m.ScanFrom(row)
}

func Update(m IBaseInfo, strSet string, strWhere string, params ...any) (*scanner.TResult, error) {
	r, err := mysql.Exec(m.Dataname(), "UPDATE "+m.Tablename()+" SET "+strSet+" "+strWhere, params...)
	return r, err
}

type IBaseList interface {
	Params() string
	Dataname() string  // 数据库名称
	Tablename() string // 表名称
	ScanFrom(scanner.IScanner) error
	GetList()
}

func SelectList(m IBaseList, strWhere string, params ...any) error {
	rows, err := mysql.QueryRows("dbaa", "SELECT "+m.Params()+" FROM "+m.Tablename()+" "+strWhere, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
