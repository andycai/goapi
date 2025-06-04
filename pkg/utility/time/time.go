package time

import (
	"fmt"
	"time"
)

// GetCurrentTimestamp 获取当前Unix时间戳(秒)
func GetCurrentTimestamp() int64 {
	return time.Now().Unix()
}

// GetCurrentTimestampMS 获取当前Unix时间戳(毫秒)
func GetCurrentTimestampMS() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// GetCurrentTimestampNS 获取当前Unix时间戳(纳秒)
func GetCurrentTimestampNS() int64 {
	return time.Now().UnixNano()
}

// FormatTimestamp 格式化Unix时间戳为字符串
func FormatTimestamp(timestamp int64, layout string) string {
	return time.Unix(timestamp, 0).Format(layout)
}

// FormatTimestampMS 格式化毫秒级Unix时间戳为字符串
func FormatTimestampMS(timestampMS int64, layout string) string {
	sec := timestampMS / 1000
	nsec := (timestampMS % 1000) * 1000000
	return time.Unix(sec, nsec).Format(layout)
}

// FormatTime 将时间格式化为字符串
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime 解析时间字符串为时间对象
func ParseTime(timeStr, layout string) (time.Time, error) {
	return time.Parse(layout, timeStr)
}

// FormatToISO8601 格式化时间为ISO8601标准格式
func FormatToISO8601(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatToRFC3339 格式化时间为RFC3339标准格式
func FormatToRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}

// FormatToRFC822 格式化时间为RFC822标准格式
func FormatToRFC822(t time.Time) string {
	return t.Format(time.RFC822)
}

// FormatToDate 格式化时间为日期格式(yyyy-mm-dd)
func FormatToDate(t time.Time) string {
	return t.Format("2006-01-02")
}

// FormatToDateTime 格式化时间为日期时间格式(yyyy-mm-dd HH:MM:SS)
func FormatToDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// GetDayStart 获取指定日期的开始时间(00:00:00)
func GetDayStart(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

// GetDayEnd 获取指定日期的结束时间(23:59:59.999999999)
func GetDayEnd(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// GetWeekStart 获取指定日期所在周的开始时间(周日或周一作为一周的开始)
func GetWeekStart(t time.Time, sundayStart bool) time.Time {
	weekday := int(t.Weekday())
	if !sundayStart && weekday == 0 {
		weekday = 7
	}
	year, month, day := t.Date()
	return time.Date(year, month, day-weekday, 0, 0, 0, 0, t.Location())
}

// GetWeekEnd 获取指定日期所在周的结束时间
func GetWeekEnd(t time.Time, sundayStart bool) time.Time {
	start := GetWeekStart(t, sundayStart)
	return start.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Nanosecond*999999999)
}

// GetMonthStart 获取指定日期所在月的开始时间
func GetMonthStart(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

// GetMonthEnd 获取指定日期所在月的结束时间
func GetMonthEnd(t time.Time) time.Time {
	year, month, _ := t.Date()
	last := time.Date(year, month+1, 0, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
	return last
}

// GetYearStart 获取指定日期所在年的开始时间
func GetYearStart(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

// GetYearEnd 获取指定日期所在年的结束时间
func GetYearEnd(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 12, 31, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
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

// DiffInSeconds 计算两个时间相差的秒数
func DiffInSeconds(t1, t2 time.Time) int64 {
	return int64(t1.Sub(t2).Seconds())
}

// DiffInMinutes 计算两个时间相差的分钟数
func DiffInMinutes(t1, t2 time.Time) int64 {
	return int64(t1.Sub(t2).Minutes())
}

// DiffInHours 计算两个时间相差的小时数
func DiffInHours(t1, t2 time.Time) int64 {
	return int64(t1.Sub(t2).Hours())
}

// DiffInDays 计算两个时间相差的天数
func DiffInDays(t1, t2 time.Time) int {
	hours := t1.Sub(t2).Hours()
	return int(hours / 24)
}

// IsLeapYear 判断是否闰年
func IsLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// GetDaysInMonth 获取某月的天数
func GetDaysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// GetHumanReadableDuration 获取人类可读的时间间隔
func GetHumanReadableDuration(d time.Duration) string {
	days := int64(d.Hours() / 24)
	hours := int64(d.Hours()) % 24
	minutes := int64(d.Minutes()) % 60
	seconds := int64(d.Seconds()) % 60

	if days > 0 {
		return fmt.Sprintf("%d天%d小时%d分钟%d秒", days, hours, minutes, seconds)
	} else if hours > 0 {
		return fmt.Sprintf("%d小时%d分钟%d秒", hours, minutes, seconds)
	} else if minutes > 0 {
		return fmt.Sprintf("%d分钟%d秒", minutes, seconds)
	}
	return fmt.Sprintf("%d秒", seconds)
}

// IsSameDay 判断两个时间是否是同一天
func IsSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameMonth 判断两个时间是否是同一月
func IsSameMonth(t1, t2 time.Time) bool {
	y1, m1, _ := t1.Date()
	y2, m2, _ := t2.Date()
	return y1 == y2 && m1 == m2
}

// IsSameYear 判断两个时间是否是同一年
func IsSameYear(t1, t2 time.Time) bool {
	y1, _, _ := t1.Date()
	y2, _, _ := t2.Date()
	return y1 == y2
}

// GetWeekday 获取星期几(0-6，0表示周日)
func GetWeekday(t time.Time) int {
	return int(t.Weekday())
}

// GetWeekdayName 获取星期几的名称
func GetWeekdayName(t time.Time, locale string) string {
	weekday := t.Weekday()

	switch locale {
	case "zh-CN", "zh_CN":
		names := []string{"星期日", "星期一", "星期二", "星期三", "星期四", "星期五", "星期六"}
		return names[weekday]
	case "en-US", "en_US", "en":
		names := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
		return names[weekday]
	default:
		// 默认使用英文
		names := []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"}
		return names[weekday]
	}
}

// GetQuarterStart 获取季度开始时间
func GetQuarterStart(t time.Time) time.Time {
	year, month, _ := t.Date()
	quarter := (int(month) - 1) / 3
	firstMonthOfQuarter := time.Month(quarter*3 + 1)
	return time.Date(year, firstMonthOfQuarter, 1, 0, 0, 0, 0, t.Location())
}

// GetQuarterEnd 获取季度结束时间
func GetQuarterEnd(t time.Time) time.Time {
	year, month, _ := t.Date()
	quarter := (int(month) - 1) / 3
	firstMonthOfNextQuarter := time.Month(quarter*3 + 4)
	if firstMonthOfNextQuarter > 12 {
		firstMonthOfNextQuarter = 1
		year++
	}
	return time.Date(year, firstMonthOfNextQuarter, 0, 23, 59, 59, int(time.Second-time.Nanosecond), t.Location())
}

// GetCurrentQuarter 获取当前季度(1-4)
func GetCurrentQuarter(t time.Time) int {
	_, month, _ := t.Date()
	return (int(month)-1)/3 + 1
}
