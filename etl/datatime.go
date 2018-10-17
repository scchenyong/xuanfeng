package etl

import (
	"time"
)

// 处理时间单位
type TimeUnit string

const (
	// 年
	Year = "year"
	// 月
	Month = "month"
	// 日
	Day = "day"
	// 小时
	Hour = "hour"
	// 分钟
	Minute = "minute"
	// 秒钟
	Second = "second"
)

// 数据处理时间结构
type DataTime struct {
	// 时间单位
	Unit TimeUnit
	// 处理偏移量
	Offset int
	// 数据沉积量
	Deposit int
}

// 获取处理偏移量时间
func (dt *DataTime) OffsetTime(t time.Time) time.Time {
	return AddTime(t, dt.Offset, dt.Unit)
}

// 获取数据沉积量时间
func (dt *DataTime) DepositTime(t time.Time) time.Time {
	return AddTime(t, dt.Deposit, dt.Unit)
}

// 转换时间
func AddTime(t time.Time, offset int, unit TimeUnit) time.Time {
	switch unit {
	case Year:
		return t.AddDate(offset, 0, 0)
	case Month:
		return t.AddDate(0, offset, 0)
	case Day:
		return t.AddDate(0, 0, offset)
	case Hour:
		return t.Add(time.Hour * time.Duration(offset))
	case Minute:
		return t.Add(time.Minute * time.Duration(offset))
	case Second:
		return t.Add(time.Second * time.Duration(offset))
	}
	// 原样返回
	return t
}
