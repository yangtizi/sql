package cache

import (
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	// 永不过期
	NoExpiration time.Duration = -1
	// 默认的过期时间,  具体要看过期时间值
	DefaultExpiration time.Duration = 0
)

type TCache struct {
	FileName          string                    // 文件名
	defaultExpiration time.Duration             // 默认超时
	items             map[string]TItem          // 内容
	mu                sync.RWMutex              // 互斥锁
	onEvicted         func(string, interface{}) // 删除回调
	timer             *TTimer                   // 自动垃圾回收定时器
}

// 添加缓存, 并且, 并且更新时间
func (m *TCache) Set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = m.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	m.mu.Lock()
	m.items[k] = TItem{
		Object:     x,
		Expiration: e,
	}
	m.mu.Unlock()
}

func (m *TCache) set(k string, x interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = m.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	m.items[k] = TItem{
		Object:     x,
		Expiration: e,
	}
}

// 使用默认超时时间
func (m *TCache) SetDefault(k string, x interface{}) {
	m.Set(k, x, DefaultExpiration)
}

// Add 和 Set的区别是. 如果存在 Set 会替换, 但是Add会报错
func (m *TCache) Add(k string, x interface{}, d time.Duration) error {
	m.mu.Lock()
	_, found := m.get(k)
	if found {
		m.mu.Unlock()
		return fmt.Errorf("已经存在 Item %s ", k)
	}
	m.set(k, x, d)
	m.mu.Unlock()
	return nil
}

// Replace 和 Set 的区别是  Replace 必须要存在, 不然会报错
func (m *TCache) Replace(k string, x interface{}, d time.Duration) error {
	m.mu.Lock()
	_, found := m.get(k)
	if !found {
		m.mu.Unlock()
		return fmt.Errorf("不存在 Item %s", k)
	}
	m.set(k, x, d)
	m.mu.Unlock()
	return nil
}

// 第一个是返回缓存值,  第二个是是否有值
func (m *TCache) Get(k string) (interface{}, bool) {
	m.mu.RLock()
	item, found := m.items[k]
	if !found {
		m.mu.RUnlock()
		return nil, false
	}
	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			m.mu.RUnlock()
			return nil, false
		}
	}
	m.mu.RUnlock()
	return item.Object, true
}

// 同时返回值和过期时间,  第三个参数返回是否有值
func (m *TCache) GetWithExpiration(k string) (interface{}, time.Time, bool) {
	m.mu.RLock()
	item, found := m.items[k]
	if !found {
		m.mu.RUnlock()
		return nil, time.Time{}, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			m.mu.RUnlock()
			return nil, time.Time{}, false
		}

		m.mu.RUnlock()
		return item.Object, time.Unix(0, item.Expiration), true
	}

	m.mu.RUnlock()
	return item.Object, time.Time{}, true
}

func (m *TCache) get(k string) (interface{}, bool) {
	item, found := m.items[k]
	if !found {
		return nil, false
	}

	if item.Expiration > 0 {
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}
	}
	return item.Object, true
}

// 增加, 如果不可以增加会返回错误, 不建议这样用, 性能不好
func (m *TCache) Increment(k string, n int64) error {
	m.mu.Lock()
	v, found := m.items[k]
	if !found || v.Expired() {
		m.mu.Unlock()
		return fmt.Errorf("找不到 Item %s ", k)
	}
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) + int(n)
	case int8:
		v.Object = v.Object.(int8) + int8(n)
	case int16:
		v.Object = v.Object.(int16) + int16(n)
	case int32:
		v.Object = v.Object.(int32) + int32(n)
	case int64:
		v.Object = v.Object.(int64) + n
	case uint:
		v.Object = v.Object.(uint) + uint(n)
	case uintptr:
		v.Object = v.Object.(uintptr) + uintptr(n)
	case uint8:
		v.Object = v.Object.(uint8) + uint8(n)
	case uint16:
		v.Object = v.Object.(uint16) + uint16(n)
	case uint32:
		v.Object = v.Object.(uint32) + uint32(n)
	case uint64:
		v.Object = v.Object.(uint64) + uint64(n)
	case float32:
		v.Object = v.Object.(float32) + float32(n)
	case float64:
		v.Object = v.Object.(float64) + float64(n)
	default:
		m.mu.Unlock()
		return fmt.Errorf("[%s] 的值不是数值", k)
	}
	m.items[k] = v
	m.mu.Unlock()
	return nil
}

// 减法
func (m *TCache) Decrement(k string, n int64) error {
	m.mu.Lock()
	v, found := m.items[k]
	if !found || v.Expired() {
		m.mu.Unlock()
		return fmt.Errorf("找不到 Item %s", k)
	}
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) - int(n)
	case int8:
		v.Object = v.Object.(int8) - int8(n)
	case int16:
		v.Object = v.Object.(int16) - int16(n)
	case int32:
		v.Object = v.Object.(int32) - int32(n)
	case int64:
		v.Object = v.Object.(int64) - n
	case uint:
		v.Object = v.Object.(uint) - uint(n)
	case uintptr:
		v.Object = v.Object.(uintptr) - uintptr(n)
	case uint8:
		v.Object = v.Object.(uint8) - uint8(n)
	case uint16:
		v.Object = v.Object.(uint16) - uint16(n)
	case uint32:
		v.Object = v.Object.(uint32) - uint32(n)
	case uint64:
		v.Object = v.Object.(uint64) - uint64(n)
	case float32:
		v.Object = v.Object.(float32) - float32(n)
	case float64:
		v.Object = v.Object.(float64) - float64(n)
	default:
		m.mu.Unlock()
		return fmt.Errorf("[%s] 的值不是数值", k)
	}
	m.items[k] = v
	m.mu.Unlock()
	return nil
}

// 删除缓存, 没有的话并不会报错
func (m *TCache) Delete(k string) {
	m.mu.Lock()
	v, evicted := m.delete(k)
	m.mu.Unlock()
	if evicted {
		m.onEvicted(k, v)
	}
}

func (m *TCache) delete(k string) (interface{}, bool) {
	if m.onEvicted != nil {
		if v, found := m.items[k]; found {
			delete(m.items, k)
			return v.Object, true
		}
	}
	delete(m.items, k)
	return nil, false
}

// 删除超时
func (m *TCache) DeleteExpired() {
	var evictedItems []TKeyAndValue
	now := time.Now().UnixNano()
	m.mu.Lock()
	for k, v := range m.items {
		// "Inlining" of expired
		if v.Expiration > 0 && now > v.Expiration {
			ov, evicted := m.delete(k)
			if evicted {
				evictedItems = append(evictedItems, TKeyAndValue{k, ov})
			}
		}
	}
	m.mu.Unlock()
	for _, v := range evictedItems {
		m.onEvicted(v.Key, v.Value)
	}
}

// 设置删除回调, 如果要删除 [删除回调], 传nil
func (m *TCache) SetOnEvicted(f func(string, interface{})) {
	m.mu.Lock()
	m.onEvicted = f
	m.mu.Unlock()
}

// Copies all unexpired items in the cache into a new map and returns it.
func (m *TCache) Items() map[string]TItem {
	m.mu.RLock()
	defer m.mu.RUnlock()
	items := make(map[string]TItem, len(m.items))
	now := time.Now().UnixNano()
	for k, v := range m.items {
		// "Inlining" of Expired
		if v.Expiration > 0 {
			if now > v.Expiration {
				continue
			}
		}
		items[k] = v
	}
	return items
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func (m *TCache) ItemCount() int {
	m.mu.RLock()
	n := len(m.items)
	m.mu.RUnlock()
	return n
}

// Delete all items from the cache.
func (m *TCache) Flush() {
	m.mu.Lock()
	m.items = map[string]TItem{}
	m.mu.Unlock()
}

func (m *TCache) Load(r io.Reader) error {
	dec := gob.NewDecoder(r)
	items := map[string]TItem{}
	err := dec.Decode(&items)
	if err == nil {
		m.mu.Lock()
		defer m.mu.Unlock()
		for k, v := range items {
			ov, found := m.items[k]
			if !found || ov.Expired() {
				m.items[k] = v
			}
		}
	}
	return err
}
func (m *TCache) Save(w io.Writer) (err error) {
	enc := gob.NewEncoder(w)
	defer func() {
		if x := recover(); x != nil {
			err = fmt.Errorf("Error registering item types with Gob library")
		}
	}()
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, v := range m.items {
		gob.Register(v.Object)
	}
	err = enc.Encode(&m.items)
	return
}

func (m *TCache) LoadFile() error {
	if m.FileName == "" {
		return errors.New("no file name")
	}
	fp, err := os.Open(m.FileName)
	if err != nil {
		return err
	}
	err = m.Load(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

func (m *TCache) SaveFile() error {
	if m.FileName == "" {
		return errors.New("no file name")
	}
	fp, err := os.Create(m.FileName)
	if err != nil {
		return err
	}
	err = m.Save(fp)
	if err != nil {
		fp.Close()
		return err
	}
	return fp.Close()
}

// 返回具有给定默认过期时间和清理的新缓存
// 间隔。如果过期时间小于一（或无过期），
// 缓存中的项目永远不会过期（默认情况下），必须删除
// 手动。如果清理间隔小于一，则不包括过期项目
// 在调用c.DeleteExpired（）之前从缓存中删除。
func New(de time.Duration, ci time.Duration) *TCache {
	if de == 0 {
		de = -1
	}
	pCache := &TCache{
		defaultExpiration: de,
		items:             map[string]TItem{},
	}

	// 自动释放
	if ci > 0 {
		pCache.timer = &TTimer{
			Interval: ci,
			stop:     make(chan bool),
		}
		go pCache.timer.Run(pCache)

		// 设置析构
		runtime.SetFinalizer(pCache, stopTimer)
	}
	return pCache
}
