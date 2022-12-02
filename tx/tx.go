package tx

import (
	"database/sql"

	"github.com/yangtizi/log/zaplog"
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

func (m *TTx) Exec(strQuery string, args ...interface{}) (sql.Result, error) {
	zaplog.Map("Tx").Debugf("strQuery = [%s]", strQuery)
	zaplog.Map("Tx").Debug("[+] ", args)
	return m.tx.Exec(strQuery, args...)
}

func (m *TTx) Commit() error {
	zaplog.Map("Tx").Debugf("Commit")
	return m.tx.Commit()
}
func (m *TTx) Rollback() error {
	zaplog.Map("Tx").Debugf("Rollback")
	return m.tx.Rollback()
}
