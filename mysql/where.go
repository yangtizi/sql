package mysql

import (
	"errors"
	"strings"

	"github.com/yangtizi/log/zaplog"
)

type TWhere struct {
	strWhere string
	params   []interface{}
}

func (m *TWhere) Append(s string, args ...interface{}) {
	m.strWhere += s + " "
	m.params = append(m.params, args...)
}

// 检查
func (m *TWhere) Check() error {
	if len(m.params) != strings.Count(m.strWhere, "?") {
		zaplog.Ins.Errorf("%s, %v", m.strWhere, m.params)
		return errors.New("参数数量不对")
	}
	return nil
}

// 获取名称
func (m *TWhere) String() string {
	return m.strWhere
}

// 获取参数
func (m *TWhere) Params() []interface{} {
	return m.params
}
