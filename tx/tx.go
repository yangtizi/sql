package tx

import (
	"database/sql"

	"github.com/yangtizi/log/zaplog"
	"github.com/yangtizi/sql/scanner"
)

// TTx 封装的
type TTx struct {
	tx *sql.Tx
}

// NewTx 新建
func NewTx(tx *sql.Tx) *TTx {
	t := &TTx{}
	t.tx = tx
	return t
}

func (m *TTx) Exec(strQuery string, args ...interface{}) (*scanner.TResult, error) {
	zaplog.Ins.Debugf("strQuery = [%s]", strQuery)
	zaplog.Ins.Debug("[+] ", args)
	rs, err := m.tx.Exec(strQuery, args...)

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

func (m *TTx) Commit() error {
	zaplog.Ins.Debugf("Commit")
	return m.tx.Commit()
}
func (m *TTx) Rollback() error {
	zaplog.Ins.Debugf("Rollback")
	return m.tx.Rollback()
}
