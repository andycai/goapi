package fund

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var config *FundConfig

// initService 初始化服务
func initService() {
	config = &FundConfig{
		ConfigPath: "./data/fund_config.json",
		Source:     "eastmoney", // 默认使用东方财富数据
		Interval:   5,           // 默认5分钟更新一次
	}

	// 尝试加载配置
	loadConfig()
}

// saveConfig 保存配置
func saveConfig() error {
	if config == nil {
		return errors.New("配置为空")
	}

	// 确保目录存在
	dir := filepath.Dir(config.ConfigPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 序列化配置
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(config.ConfigPath, data, 0644)
}

// loadConfig 加载配置
func loadConfig() error {
	// 检查文件是否存在
	if _, err := os.Stat(config.ConfigPath); os.IsNotExist(err) {
		// 配置文件不存在，使用默认配置
		return nil
	}

	// 读取文件
	data, err := os.ReadFile(config.ConfigPath)
	if err != nil {
		return err
	}

	// 反序列化配置
	return json.Unmarshal(data, config)
}

// updateConfig 更新配置
func updateConfig(conf *FundConfig) error {
	if conf.Interval <= 0 {
		return errors.New("更新间隔必须大于0")
	}

	// 更新配置
	config.Source = conf.Source
	config.ApiKey = conf.ApiKey
	config.Interval = conf.Interval

	// 保存配置
	return saveConfig()
}

// getConfig 获取配置
func getConfig() *FundConfig {
	return config
}

// SyncData 同步数据
func SyncData() error {
	if config == nil {
		return errors.New("配置为空，请先设置配置")
	}

	// 记录错误，但不立即返回，尝试执行所有同步操作
	var syncErrors []string

	// 同步市场指数
	if err := syncMarketData(); err != nil {
		syncErrors = append(syncErrors, fmt.Sprintf("同步市场指数失败: %v", err))
	}

	// 同步热门基金
	if err := syncHotFunds(); err != nil {
		syncErrors = append(syncErrors, fmt.Sprintf("同步热门基金失败: %v", err))
	}

	// 如果有错误，合并返回
	if len(syncErrors) > 0 {
		return errors.New(strings.Join(syncErrors, "; "))
	}

	return nil
}

// syncMarketData 同步市场指数数据
func syncMarketData() error {
	// 获取现有指数列表
	indices, err := GetMarketIndices()
	if err != nil {
		return err
	}

	if len(indices) == 0 {
		return errors.New("没有配置市场指数")
	}

	// 记录更新成功的指数数量
	successCount := 0

	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// 遍历更新每个指数
	for i := range indices {
		// 构建API URL
		url := fmt.Sprintf("http://push2.eastmoney.com/api/qt/stock/get?"+
			"secid=%s&fields=f43,f44,f45,f46,f47", indices[i].Code)

		// 发送请求
		resp, err := client.Get(url)
		if err != nil {
			// 记录错误但继续处理其他指数
			fmt.Printf("获取指数 %s 数据失败: %v\n", indices[i].Name, err)
			continue
		}

		// 读取响应
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close() // 确保关闭响应体
		if err != nil {
			fmt.Printf("读取指数 %s 数据失败: %v\n", indices[i].Name, err)
			continue
		}

		// 解析JSON
		var result struct {
			Data struct {
				F43 float64 `json:"f43"` // 当前价
				F44 float64 `json:"f44"` // 涨跌额
				F45 float64 `json:"f45"` // 涨跌幅
			} `json:"data"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Printf("解析指数 %s 数据失败: %v\n", indices[i].Name, err)
			continue
		}

		// 更新数据
		indices[i].Current = result.Data.F43
		indices[i].Change = result.Data.F44
		indices[i].Growth = result.Data.F45
		indices[i].UpdatedAt = time.Now()

		// 保存到数据库
		if err := SaveMarketIndex(&indices[i]); err != nil {
			fmt.Printf("保存指数 %s 数据失败: %v\n", indices[i].Name, err)
			continue
		}

		successCount++
	}

	if successCount == 0 {
		return errors.New("所有指数更新失败")
	}

	return nil
}

// syncHotFunds 同步热门基金数据 - 使用真实数据
func syncHotFunds() error {
	// 清空现有数据
	if err := DeleteAllFunds(); err != nil {
		return fmt.Errorf("清空基金数据失败: %v", err)
	}

	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 构建请求URL - 使用东方财富最新API
	url := "https://api.fund.eastmoney.com/FundGuZhi/GetFundGZList" +
		"?type=1&sort=3&orderType=desc&canbuy=0&pageIndex=1&pageSize=100" +
		"&callback=jQuery&_=" + fmt.Sprintf("%d", time.Now().UnixNano()/1e6)

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	// 添加请求头，模拟浏览器请求
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Referer", "https://fund.eastmoney.com/")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求热门基金数据失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求热门基金数据失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取热门基金数据失败: %v", err)
	}

	// 处理JSONP响应 (jQuery(...))
	jsonStr := string(body)
	start := strings.Index(jsonStr, "(")
	end := strings.LastIndex(jsonStr, ")")

	if start >= 0 && end > start {
		jsonStr = jsonStr[start+1 : end]
	} else {
		return fmt.Errorf("解析JSONP数据失败")
	}

	// 解析JSON
	var result struct {
		Data struct {
			List []struct {
				FCODE     string `json:"FCODE"`     // 基金代码
				SHORTNAME string `json:"SHORTNAME"` // 基金名称
				FTYPE     string `json:"FTYPE"`     // 基金类型
				NAV       string `json:"NAV"`       // 单位净值
				PDATE     string `json:"PDATE"`     // 净值日期
				DAYZSYL   string `json:"DAYZSYL"`   // 日增长率
				RZDF      string `json:"RZDF"`      // 日增长率
				SYL_1M    string `json:"SYL_1M"`    // 月增长率
				SYL_3M    string `json:"SYL_3M"`    // 季增长率
				SYL_6M    string `json:"SYL_6M"`    // 半年增长率
				SYL_1Y    string `json:"SYL_1Y"`    // 年增长率
			} `json:"list"`
		} `json:"Data"`
		ErrCode int    `json:"ErrCode"`
		Success bool   `json:"Success"`
		ErrMsg  string `json:"ErrMsg"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return fmt.Errorf("解析JSON数据失败: %v", err)
	}

	// 检查API返回是否成功
	if !result.Success || result.ErrCode != 0 {
		return fmt.Errorf("API返回错误: %s", result.ErrMsg)
	}

	// 检查是否有数据
	if len(result.Data.List) == 0 {
		return fmt.Errorf("未获取到基金数据")
	}

	// 批量保存数据
	var funds []Fund
	for _, item := range result.Data.List {
		// 过滤无效数据
		if item.FCODE == "" {
			continue
		}

		// 解析数值
		netValue, _ := strconv.ParseFloat(item.NAV, 64)
		dayGrowth, _ := strconv.ParseFloat(item.DAYZSYL, 64)
		if dayGrowth == 0 {
			dayGrowth, _ = strconv.ParseFloat(item.RZDF, 64)
		}
		monthGrowth, _ := strconv.ParseFloat(item.SYL_1M, 64)
		yearGrowth, _ := strconv.ParseFloat(item.SYL_1Y, 64)

		// 解析日期
		var netValueDay time.Time
		if item.PDATE != "" {
			netValueDay, _ = time.Parse("2006-01-02", item.PDATE)
		} else {
			netValueDay = time.Now()
		}

		// 基金类型映射
		fundType := "混合型" // 默认类型
		switch item.FTYPE {
		case "001":
			fundType = "股票型"
		case "002":
			fundType = "混合型"
		case "003":
			fundType = "债券型"
		case "004":
			fundType = "指数型"
		case "005":
			fundType = "QDII"
		case "006":
			fundType = "ETF联接"
		}

		fund := Fund{
			Code:        item.FCODE,
			Name:        item.SHORTNAME,
			Type:        fundType,
			NetValue:    netValue,
			NetValueDay: netValueDay,
			DayGrowth:   dayGrowth,
			MonthGrowth: monthGrowth,
			YearGrowth:  yearGrowth,
			UpdatedAt:   time.Now(),
		}
		funds = append(funds, fund)
	}

	// 如果第一个API失败，尝试第二个API
	if len(funds) == 0 {
		return fetchAlternativeFundData()
	}

	return BatchSaveFunds(funds)
}

// fetchAlternativeFundData 获取备用基金数据
func fetchAlternativeFundData() error {
	// 创建HTTP客户端，设置超时
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// 构建备用API URL
	url := "https://fundmobapi.eastmoney.com/FundMNewApi/FundMNRank" +
		"?pageIndex=1&pageSize=100&plat=Android&appType=ttjj&product=EFund" +
		"&version=1&sortColumn=RZDF&sortType=desc"

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("创建备用请求失败: %v", err)
	}

	// 添加请求头，模拟移动端请求
	req.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1")
	req.Header.Add("Accept", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("请求备用基金数据失败: %v", err)
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("请求备用基金数据失败，状态码: %d", resp.StatusCode)
	}

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取备用基金数据失败: %v", err)
	}

	// 解析JSON
	var result struct {
		Datas []struct {
			FCODE     string `json:"FCODE"`     // 基金代码
			SHORTNAME string `json:"SHORTNAME"` // 基金简称
			FTYPE     string `json:"FTYPE"`     // 基金类型
			NAV       string `json:"NAV"`       // 单位净值
			PDATE     string `json:"PDATE"`     // 净值日期
			RZDF      string `json:"RZDF"`      // 日增长率
			SYL_1W    string `json:"SYL_1W"`    // 周增长率
			SYL_1M    string `json:"SYL_1M"`    // 月增长率
			SYL_1Y    string `json:"SYL_1Y"`    // 年增长率
		} `json:"Datas"`
		ErrCode int    `json:"ErrCode"`
		Success bool   `json:"Success"`
		Message string `json:"Message"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析备用基金数据失败: %v", err)
	}

	// 检查是否有数据
	if len(result.Datas) == 0 {
		return fmt.Errorf("未获取到备用基金数据")
	}

	// 清空现有数据
	if err := DeleteAllFunds(); err != nil {
		return fmt.Errorf("清空基金数据失败: %v", err)
	}

	// 批量保存新数据
	var funds []Fund
	for _, item := range result.Datas {
		// 过滤无效数据
		if item.FCODE == "" {
			continue
		}

		// 解析数值
		netValue, _ := strconv.ParseFloat(item.NAV, 64)
		dayGrowth, _ := strconv.ParseFloat(item.RZDF, 64)
		weekGrowth, _ := strconv.ParseFloat(item.SYL_1W, 64)
		monthGrowth, _ := strconv.ParseFloat(item.SYL_1M, 64)
		yearGrowth, _ := strconv.ParseFloat(item.SYL_1Y, 64)

		// 解析日期
		var netValueDay time.Time
		if item.PDATE != "" {
			netValueDay, _ = time.Parse("2006-01-02", item.PDATE)
		} else {
			netValueDay = time.Now()
		}

		fund := Fund{
			Code:        item.FCODE,
			Name:        item.SHORTNAME,
			Type:        item.FTYPE,
			NetValue:    netValue,
			NetValueDay: netValueDay,
			DayGrowth:   dayGrowth,
			WeekGrowth:  weekGrowth,
			MonthGrowth: monthGrowth,
			YearGrowth:  yearGrowth,
			UpdatedAt:   time.Now(),
		}
		funds = append(funds, fund)
	}

	return BatchSaveFunds(funds)
}
