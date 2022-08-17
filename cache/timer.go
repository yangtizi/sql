package cache

import "time"

// 定时器模块, 主要是为了定时删除缓存
type TTimer struct {
	Interval time.Duration // 定时器
	stop     chan bool
}

func (m *TTimer) Run(c *TCache) {
	ticker := time.NewTicker(m.Interval)
	for {
		select {
		case <-ticker.C: // 到时间, 定期清理
			c.DeleteExpired()
		case <-m.stop: // 删除计时器
			ticker.Stop()
			return
		}
	}
}

func stopTimer(pCache *TCache) {
	// 通道关闭
	pCache.timer.stop <- true
}
