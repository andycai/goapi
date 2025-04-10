package gameconf

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/andycai/goapi/models"
	"github.com/xuri/excelize/v2"
	"gopkg.in/yaml.v3"
)

// executeExport 执行导出任务
func executeExport(export *models.GameConfExport) {
	// 更新开始时间
	export.StartTime = time.Now()
	export.Status = "running"
	app.DB.Save(export)

	// 获取项目和配置表信息
	var project models.GameConfProject
	var table models.GameConfTable
	if err := app.DB.First(&project, export.ProjectID).Error; err != nil {
		updateExportStatus(export, "failed", fmt.Sprintf("获取项目信息失败: %v", err))
		return
	}
	if err := app.DB.First(&table, export.TableID).Error; err != nil {
		updateExportStatus(export, "failed", fmt.Sprintf("获取配置表信息失败: %v", err))
		return
	}

	// 读取源文件
	data, err := readSourceFile(&table)
	if err != nil {
		updateExportStatus(export, "failed", fmt.Sprintf("读取源文件失败: %v", err))
		return
	}

	// 导出数据
	if err := exportData(data, export, &project, &table); err != nil {
		updateExportStatus(export, "failed", fmt.Sprintf("导出数据失败: %v", err))
		return
	}

	// 生成代码
	if err := generateCode(data, export, &project, &table); err != nil {
		updateExportStatus(export, "failed", fmt.Sprintf("生成代码失败: %v", err))
		return
	}

	// 更新完成状态
	export.EndTime = time.Now()
	export.Duration = int(export.EndTime.Sub(export.StartTime).Seconds())
	export.Status = "success"
	app.DB.Save(export)
}

// updateExportStatus 更新导出状态
func updateExportStatus(export *models.GameConfExport, status string, errMsg string) {
	export.Status = status
	export.Error = errMsg
	export.EndTime = time.Now()
	export.Duration = int(export.EndTime.Sub(export.StartTime).Seconds())
	app.DB.Save(export)
}

// readSourceFile 读取源文件
func readSourceFile(table *models.GameConfTable) ([]map[string]interface{}, error) {
	switch strings.ToLower(table.FileType) {
	case "excel", "xlsx", "xls", "xlsm":
		return readExcelFile(table)
	case "csv":
		return readCSVFile(table)
	default:
		return nil, fmt.Errorf("不支持的文件类型: %s", table.FileType)
	}
}

// readExcelFile 读取Excel文件
func readExcelFile(table *models.GameConfTable) ([]map[string]interface{}, error) {
	f, err := excelize.OpenFile(table.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 获取工作表名称
	sheetName := table.SheetName
	if sheetName == "" {
		sheetName = f.GetSheetName(0)
	}

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel文件格式错误：至少需要标题行和数据行")
	}

	// 第一行为标题
	headers := rows[0]
	var data []map[string]interface{}

	// 从第二行开始读取数据
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		item := make(map[string]interface{})
		for j := 0; j < len(headers) && j < len(row); j++ {
			item[headers[j]] = row[j]
		}
		data = append(data, item)
	}

	return data, nil
}

// readCSVFile 读取CSV文件
func readCSVFile(table *models.GameConfTable) ([]map[string]interface{}, error) {
	content, err := os.ReadFile(table.FilePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("CSV文件格式错误：至少需要标题行和数据行")
	}

	// 第一行为标题
	headers := strings.Split(lines[0], ",")
	var data []map[string]interface{}

	// 从第二行开始读取数据
	for i := 1; i < len(lines); i++ {
		if lines[i] == "" {
			continue
		}
		values := strings.Split(lines[i], ",")
		item := make(map[string]interface{})
		for j := 0; j < len(headers) && j < len(values); j++ {
			item[headers[j]] = values[j]
		}
		data = append(data, item)
	}

	return data, nil
}

// exportData 导出数据
func exportData(data []map[string]interface{}, export *models.GameConfExport, project *models.GameConfProject, table *models.GameConfTable) error {
	fileName := fmt.Sprintf("%s.%s", table.Name, export.Format)
	filePath := filepath.Join(project.DataPath, fileName)

	switch strings.ToLower(export.Format) {
	case "json":
		return exportJSON(data, filePath)
	case "xml":
		return exportXML(data, filePath)
	case "yaml":
		return exportYAML(data, filePath)
	case "lua":
		return exportLua(data, filePath)
	case "binary":
		return exportBinary(data, filePath)
	default:
		return fmt.Errorf("不支持的导出格式: %s", export.Format)
	}
}

// exportJSON 导出JSON格式
func exportJSON(data []map[string]interface{}, filePath string) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0644)
}

// exportXML 导出XML格式
func exportXML(data []map[string]interface{}, filePath string) error {
	type Item struct {
		XMLName xml.Name
		Fields  []struct {
			Name  string `xml:"name,attr"`
			Value string `xml:",chardata"`
		} `xml:"field"`
	}

	var items []Item
	for _, row := range data {
		item := Item{XMLName: xml.Name{Local: "item"}}
		for key, value := range row {
			item.Fields = append(item.Fields, struct {
				Name  string `xml:"name,attr"`
				Value string `xml:",chardata"`
			}{
				Name:  key,
				Value: fmt.Sprintf("%v", value),
			})
		}
		items = append(items, item)
	}

	xmlData, err := xml.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}

	xmlData = append([]byte(xml.Header), xmlData...)
	return os.WriteFile(filePath, xmlData, 0644)
}

// exportYAML 导出YAML格式
func exportYAML(data []map[string]interface{}, filePath string) error {
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, yamlData, 0644)
}

// exportLua 导出Lua格式
func exportLua(data []map[string]interface{}, filePath string) error {
	var sb strings.Builder
	sb.WriteString("return {\n")

	for i, item := range data {
		sb.WriteString("  {\n")
		for key, value := range item {
			sb.WriteString(fmt.Sprintf("    %s = %v,\n", key, formatLuaValue(value)))
		}
		sb.WriteString("  }")
		if i < len(data)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}

	sb.WriteString("}\n")
	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// formatLuaValue 格式化Lua值
func formatLuaValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return fmt.Sprintf("%q", v)
	case nil:
		return "nil"
	default:
		return fmt.Sprintf("%v", v)
	}
}

// exportBinary 导出二进制格式
func exportBinary(data []map[string]interface{}, filePath string) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(filePath, jsonData, 0644)
}

// generateCode 生成代码
func generateCode(data []map[string]interface{}, export *models.GameConfExport, project *models.GameConfProject, table *models.GameConfTable) error {
	switch strings.ToLower(export.Language) {
	case "cs":
		return generateCSharpCode(data, export, project, table)
	case "java":
		return generateJavaCode(data, export, project, table)
	case "go":
		return generateGoCode(data, export, project, table)
	default:
		return fmt.Errorf("不支持的目标语言: %s", export.Language)
	}
}

// generateCSharpCode 生成C#代码
func generateCSharpCode(data []map[string]interface{}, export *models.GameConfExport, project *models.GameConfProject, table *models.GameConfTable) error {
	className := ToPascalCase(table.Name)
	fileName := fmt.Sprintf("%s.cs", className)
	filePath := filepath.Join(project.CodePath, "csharp", fileName)

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("using System;\n")
	sb.WriteString("using System.Collections.Generic;\n\n")
	sb.WriteString("namespace GameConfig\n{\n")

	// 生成类定义
	sb.WriteString(fmt.Sprintf("    public class %s\n    {\n", className))

	// 从第一行数据推断字段类型
	if len(data) > 0 {
		for key, value := range data[0] {
			fieldName := ToPascalCase(key)
			fieldType := getCSharpType(value)
			sb.WriteString(fmt.Sprintf("        public %s %s { get; set; }\n", fieldType, fieldName))
		}
	}

	sb.WriteString("    }\n")
	sb.WriteString("}\n")

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// generateJavaCode 生成Java代码
func generateJavaCode(data []map[string]interface{}, export *models.GameConfExport, project *models.GameConfProject, table *models.GameConfTable) error {
	className := ToPascalCase(table.Name)
	fileName := fmt.Sprintf("%s.java", className)
	filePath := filepath.Join(project.CodePath, "java", fileName)

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("package com.game.config;\n\n")
	sb.WriteString("import java.util.*;\n\n")

	// 生成类定义
	sb.WriteString(fmt.Sprintf("public class %s {\n", className))

	// 从第一行数据推断字段类型
	if len(data) > 0 {
		for key, value := range data[0] {
			fieldName := ToCamelCase(key)
			fieldType := getJavaType(value)
			sb.WriteString(fmt.Sprintf("    private %s %s;\n", fieldType, fieldName))
		}
	}

	// 生成getter和setter方法
	if len(data) > 0 {
		sb.WriteString("\n")
		for key, value := range data[0] {
			fieldName := ToCamelCase(key)
			fieldType := getJavaType(value)

			// Getter
			sb.WriteString(fmt.Sprintf("    public %s get%s() {\n", fieldType, ToPascalCase(key)))
			sb.WriteString(fmt.Sprintf("        return %s;\n", fieldName))
			sb.WriteString("    }\n\n")

			// Setter
			sb.WriteString(fmt.Sprintf("    public void set%s(%s %s) {\n", ToPascalCase(key), fieldType, fieldName))
			sb.WriteString(fmt.Sprintf("        this.%s = %s;\n", fieldName, fieldName))
			sb.WriteString("    }\n\n")
		}
	}

	sb.WriteString("}\n")

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// generateGoCode 生成Go代码
func generateGoCode(data []map[string]interface{}, export *models.GameConfExport, project *models.GameConfProject, table *models.GameConfTable) error {
	structName := ToPascalCase(table.Name)
	fileName := fmt.Sprintf("%s.go", ToSnakeCase(table.Name))
	filePath := filepath.Join(project.CodePath, "go", fileName)

	// 创建目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return err
	}

	var sb strings.Builder
	sb.WriteString("package config\n\n")

	// 生成结构体定义
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))

	// 从第一行数据推断字段类型
	if len(data) > 0 {
		for key, value := range data[0] {
			fieldName := ToPascalCase(key)
			fieldType := getGoType(value)
			sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\"`\n", fieldName, fieldType, key))
		}
	}

	sb.WriteString("}\n")

	return os.WriteFile(filePath, []byte(sb.String()), 0644)
}

// getCSharpType 获取C#类型
func getCSharpType(value interface{}) string {
	switch value.(type) {
	case int, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case bool:
		return "bool"
	default:
		return "string"
	}
}

// getJavaType 获取Java类型
func getJavaType(value interface{}) string {
	switch value.(type) {
	case int, int32, int64:
		return "int"
	case float32, float64:
		return "float"
	case bool:
		return "boolean"
	default:
		return "String"
	}
}

// getGoType 获取Go类型
func getGoType(value interface{}) string {
	switch value.(type) {
	case int, int32, int64:
		return "int"
	case float32, float64:
		return "float64"
	case bool:
		return "bool"
	default:
		return "string"
	}
}

// ToPascalCase 将字符串转换为帕斯卡命名
func ToPascalCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return !strings.ContainsRune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789", r)
	})
	for i := 0; i < len(words); i++ {
		word := []rune(strings.ToLower(words[i]))
		if len(word) > 0 {
			word[0] = unicode.ToUpper(word[0])
		}
		words[i] = string(word)
	}
	return strings.Join(words, "")
}

// ToCamelCase 将字符串转换为驼峰命名
func ToCamelCase(s string) string {
	s = ToPascalCase(s)
	if len(s) == 0 {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// ToSnakeCase 将字符串转换为蛇形命名
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && (unicode.IsUpper(r) || unicode.IsNumber(r)) &&
			((i+1 < len(s) && unicode.IsLower(rune(s[i+1]))) || unicode.IsLower(rune(s[i-1]))) {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}
