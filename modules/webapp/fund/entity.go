package fund

import "time"

// Fund 基金信息
type Fund struct {
	Code        string    `json:"code" gorm:"primaryKey;size:6"`    // 基金代码
	Name        string    `json:"name" gorm:"size:100"`             // 基金名称
	Type        string    `json:"type" gorm:"size:20"`              // 基金类型
	NetValue    float64   `json:"net_value"`                        // 最新净值
	NetValueDay time.Time `json:"net_value_day"`                    // 净值日期
	DayGrowth   float64   `json:"day_growth"`                       // 日增长率
	WeekGrowth  float64   `json:"week_growth"`                      // 周增长率
	MonthGrowth float64   `json:"month_growth"`                     // 月增长率
	YearGrowth  float64   `json:"year_growth"`                      // 年增长率
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"` // 更新时间
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"` // 创建时间
}

// MarketIndex 市场指数
type MarketIndex struct {
	Code      string    `json:"code" gorm:"primaryKey;size:20"` // 指数代码
	Name      string    `json:"name" gorm:"size:50"`            // 指数名称
	Current   float64   `json:"current"`                        // 当前点数
	Change    float64   `json:"change"`                         // 涨跌点数
	Growth    float64   `json:"growth"`                         // 涨跌幅
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// FundConfig 配置信息
type FundConfig struct {
	ConfigPath string `json:"-"`        // 配置文件路径
	Source     string `json:"source"`   // 数据来源
	ApiKey     string `json:"api_key"`  // API密钥
	Interval   int    `json:"interval"` // 自动更新间隔（分钟）
}
