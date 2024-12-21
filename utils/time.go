package utils

import (
	"fmt"
	"time"
)

// FormatTime 格式化时间
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime 解析时间字符串
func ParseTime(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}

// GetCurrentTimestamp 获取当前时间戳
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentMilliseconds 获取当前毫秒时间戳
func GetCurrentMilliseconds() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetCurrentMicroseconds 获取当前微秒时间戳
func GetCurrentMicroseconds() int64 {
	return time.Now().UnixNano() / 1e3
}

// GetCurrentNanoseconds 获取当前纳秒时间戳
func GetCurrentNanoseconds() int64 {
	return time.Now().UnixNano()
}

// FormatDuration 格式化时间间隔
func FormatDuration(d time.Duration) string {
	days := int64(d.Hours() / 24)
	hours := int64(d.Hours()) % 24
	minutes := int64(d.Minutes()) % 60
	seconds := int64(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd%dh%dm%ds", days, hours, minutes, seconds)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh%dm%ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
}

// GetStartOfDay 获取一天的开始时间
func GetStartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// GetEndOfDay 获取一天的结束时间
func GetEndOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, 999999999, t.Location())
}

// GetStartOfWeek 获取一周的开始时间（周一）
func GetStartOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := t.Date()
	return time.Date(year, month, day-weekday+1, 0, 0, 0, 0, t.Location())
}

// GetEndOfWeek 获取一周的结束时间（周日）
func GetEndOfWeek(t time.Time) time.Time {
	weekday := int(t.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	year, month, day := t.Date()
	return time.Date(year, month, day+(7-weekday), 23, 59, 59, 999999999, t.Location())
}

// GetStartOfMonth 获取一个月的开始时间
func GetStartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfMonth 获取一个月的结束时间
func GetEndOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	lastDay := time.Date(year, month+1, 0, 23, 59, 59, 999999999, t.Location())
	return lastDay
}

// GetStartOfYear 获取一年的开始时间
func GetStartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfYear 获取一年的结束时间
func GetEndOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 12, 31, 23, 59, 59, 999999999, t.Location())
}

// AddDays 添加天数
func AddDays(t time.Time, days int) time.Time {
	return t.AddDate(0, 0, days)
}

// AddMonths 添加月数
func AddMonths(t time.Time, months int) time.Time {
	return t.AddDate(0, months, 0)
}

// AddYears 添加年数
func AddYears(t time.Time, years int) time.Time {
	return t.AddDate(years, 0, 0)
}

// IsLeapYear 判断是否为闰年
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// GetDaysInMonth 获取某月的天数
func GetDaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetWeekday 获取星期几
func GetWeekday(t time.Time) string {
	weekdays := []string{"周日", "周一", "周二", "周三", "周四", "周五", "周六"}
	return weekdays[t.Weekday()]
}

// GetQuarter 获取季度
func GetQuarter(t time.Time) int {
	month := t.Month()
	return int((month-1)/3 + 1)
}

// GetStartOfQuarter 获取季度的开始时间
func GetStartOfQuarter(t time.Time) time.Time {
	year, month, _ := t.Date()
	quarter := (int(month)-1)/3 + 1
	firstMonth := time.Month((quarter-1)*3 + 1)
	return time.Date(year, firstMonth, 1, 0, 0, 0, 0, t.Location())
}

// GetEndOfQuarter 获取季度的结束时间
func GetEndOfQuarter(t time.Time) time.Time {
	year, month, _ := t.Date()
	quarter := (int(month)-1)/3 + 1
	lastMonth := time.Month(quarter * 3)
	return time.Date(year, lastMonth+1, 0, 23, 59, 59, 999999999, t.Location())
}

// IsSameDay 判断是否为同一天
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameWeek 判断是否为同一周
func IsSameWeek(t1, t2 time.Time) bool {
	y1, w1 := t1.ISOWeek()
	y2, w2 := t2.ISOWeek()
	return y1 == y2 && w1 == w2
}

// IsSameMonth 判断是否为同一月
func IsSameMonth(t1, t2 time.Time) bool {
	y1, m1, _ := t1.Date()
	y2, m2, _ := t2.Date()
	return y1 == y2 && m1 == m2
}

// IsSameYear 判断是否为同一年
func IsSameYear(t1, t2 time.Time) bool {
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()
	return y1 == y2
}

// IsSameQuarter 判断是否为同一季度
func IsSameQuarter(t1, t2 time.Time) bool {
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()
	q1 := GetQuarter(t1)
	q2 := GetQuarter(t2)
	return y1 == y2 && q1 == q2
}

// GetAge 计算年龄
func GetAge(birthDate time.Time) int {
	now := time.Now()
	years := now.Year() - birthDate.Year()

	// 检查是否已过生日
	if now.Month() < birthDate.Month() ||
		(now.Month() == birthDate.Month() && now.Day() < birthDate.Day()) {
		years--
	}

	return years
}
