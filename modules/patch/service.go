package patch

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andycai/unitool/core/utility/path"
	"github.com/andycai/unitool/models"
)

var config *PatchConfig

// initService 初始化服务
func initService() {
	config = &PatchConfig{ConfigPath: "./data/patch_config.json"}

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
func updateConfig(conf *PatchConfig) error {
	// 验证路径
	if !path.IsValid(conf.SourceDir) || !path.IsValid(conf.TargetDir) {
		return errors.New("无效的目录路径")
	}

	// 更新配置
	config = conf

	// 保存配置
	return saveConfig()
}

// getConfig 获取配置
func getConfig() *PatchConfig {
	return config
}

// GeneratePatch 生成补丁包
func GeneratePatch(oldVersion, newVersion, description string) (*models.PatchRecord, error) {
	if config == nil {
		return nil, errors.New("配置为空")
	}

	// 比较文件差异
	changes, err := compareDirectories(config.SourceDir, config.TargetDir)
	if err != nil {
		return nil, err
	}

	// 生成补丁包
	outputZip := filepath.Join("./data/patches", fmt.Sprintf("%s_%s.zip", oldVersion, newVersion))
	if err := createPatchZip(changes, outputZip); err != nil {
		return nil, err
	}

	// 创建补丁记录
	record := &models.PatchRecord{
		OldVersion:  oldVersion,
		NewVersion:  newVersion,
		Version:     fmt.Sprintf("%s_%s", oldVersion, newVersion),
		PatchFile:   outputZip,
		Size:        0, // 文件大小将在创建zip文件后更新
		Description: description,
		FileCount:   len(changes),
		CreatedAt:   time.Now(),
	}

	// 获取补丁包文件大小
	if fileInfo, err := os.Stat(outputZip); err == nil {
		record.Size = fileInfo.Size()
	}

	// 保存到数据库
	if err := app.DB.Create(record).Error; err != nil {
		return nil, err
	}

	return record, nil
}

// compareDirectories 比较两个目录的差异
func compareDirectories(sourceDir, targetDir string) ([]FileChange, error) {
	var changes []FileChange

	// 遍历源目录
	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}

		// 计算目标文件路径
		targetPath := filepath.Join(targetDir, relPath)

		// 检查目标文件是否存在
		targetInfo, err := os.Stat(targetPath)
		if os.IsNotExist(err) {
			// 文件在目标目录不存在，标记为新增
			checksum, size, err := getFileInfo(path)
			if err != nil {
				return err
			}

			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: "A",
				Checksum:   checksum,
				Size:       size,
			})
			return nil
		} else if err != nil {
			return err
		}

		// 比较文件内容
		if info.Size() != targetInfo.Size() || info.ModTime() != targetInfo.ModTime() {
			checksum, size, err := getFileInfo(path)
			if err != nil {
				return err
			}

			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: "M",
				Checksum:   checksum,
				Size:       size,
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 查找已删除的文件
	err = filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 计算相对路径
		relPath, err := filepath.Rel(targetDir, path)
		if err != nil {
			return err
		}

		// 检查源文件是否存在
		sourcePath := filepath.Join(sourceDir, relPath)
		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			// 文件在源目录不存在，标记为删除
			changes = append(changes, FileChange{
				Path:       relPath,
				ChangeType: "D",
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return changes, nil
}

// getFileInfo 获取文件信息
func getFileInfo(path string) (string, int64, error) {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	// 计算MD5
	hash := md5.New()
	size, err := io.Copy(hash, file)
	if err != nil {
		return "", 0, err
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum, size, nil
}

// createPatchZip 创建补丁包
func createPatchZip(changes []FileChange, outputZip string) error {
	// 确保输出目录存在
	if err := os.MkdirAll(filepath.Dir(outputZip), 0755); err != nil {
		return err
	}

	// 创建zip文件
	zipFile, err := os.Create(outputZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// 添加文件到zip
	for _, change := range changes {
		if change.ChangeType == "D" {
			continue // 跳过删除的文件
		}

		// 打开源文件
		sourcePath := filepath.Join(config.SourceDir, change.Path)
		sourceFile, err := os.Open(sourcePath)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// 构建zip中的文件路径，添加 publish/android/ 前缀
		zipPath := filepath.Join("publish/android", change.Path)

		// 创建zip文件条目
		fileWriter, err := zipWriter.Create(zipPath)
		if err != nil {
			return err
		}

		// 复制文件内容
		if _, err := io.Copy(fileWriter, sourceFile); err != nil {
			return err
		}
	}

	return nil
}

// GetPatchRecords 获取补丁记录列表
func GetPatchRecords(limit, page int) ([]models.PatchRecord, int, error) {
	var records []models.PatchRecord
	var totalCount int64

	// 获取总记录数
	if err := app.DB.Model(&models.PatchRecord{}).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * limit
	if err := app.DB.Order("created_at desc").Offset(offset).Limit(limit).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, int(totalCount), nil
}

// ApplyPatch 应用补丁包
func ApplyPatch(recordID uint) error {
	// 获取补丁记录
	var record models.PatchRecord
	if err := app.DB.First(&record, recordID).Error; err != nil {
		return err
	}

	// 检查补丁包是否存在
	if _, err := os.Stat(record.PatchFile); os.IsNotExist(err) {
		return errors.New("补丁包文件不存在")
	}

	// 打开zip文件
	zip, err := zip.OpenReader(record.PatchFile)
	if err != nil {
		return err
	}
	defer zip.Close()

	// 应用补丁
	for _, file := range zip.File {
		// 处理路径，移除publish/android前缀
		filePath := file.Name
		if strings.HasPrefix(filePath, "publish/android/") {
			filePath = strings.TrimPrefix(filePath, "publish/android/")
		}

		// 构建目标文件路径
		targetPath := filepath.Join(config.TargetDir, filePath)

		// 确保目标目录存在
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// 创建目标文件
		targetFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer targetFile.Close()

		// 打开zip中的文件
		sourceFile, err := file.Open()
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// 复制文件内容
		if _, err := io.Copy(targetFile, sourceFile); err != nil {
			return err
		}
	}

	// 更新补丁记录状态为已应用
	record.Status = 1 // 1-已应用
	record.UpdatedAt = time.Now()

	// 更新补丁状态
	return app.DB.Save(&record).Error
}
