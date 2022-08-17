package cache

import "time"

// 记录项目的具体内容
type TItem struct {
	Object     interface{} // 对象
	Expiration int64       // 超时时间
}

// 如果项目已过期，则返回true
func (m *TItem) Expired() bool {
	if m.Expiration == 0 {
		return false
	}
	return time.Now().UnixNano() > m.Expiration
}
