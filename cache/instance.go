package cache

import "time"

// 单例, 方便直接调用函数使用
var instance = &TCache{}

func init() {
	if instance == nil {
		instance = New(12*time.Hour, 10*time.Minute) // 缓存默认时间12小时,   10分钟检测一次无用的缓存
	}
}

// 获取缓存
func Get(k string) (interface{}, bool) {
	return instance.Get(k)
}

// 设置缓存
func Set(k string, x interface{}) {
	instance.SetDefault(k, x)
}

// 设置缓存, 独立时间
func SetWithTime(k string, x interface{}, d time.Duration) {
	instance.Set(k, x, d)
}

// 删除
func Delete(k string) {
	instance.Delete(k)
}

// 清空
func Clean() {
	instance.Flush()
}

// 保存到文件
func SaveToFile(strFilename string) error {
	instance.FileName = strFilename
	return instance.SaveFile()
}

//
func LoadFromFile(strFilename string) error {
	instance.FileName = strFilename
	return instance.LoadFile()
}
