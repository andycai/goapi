package unitool

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	guidRegex = regexp.MustCompile(`guid:\s*([a-zA-Z0-9]+)`)
)

// initService 初始化服务
func initService() {
	// 初始化服务逻辑
}

// FindDuplicateGuids 在指定目录中查找重复的GUID
func FindDuplicateGuids(targetPath, notificationURL string) (*FindGuidLog, error) {
	// 创建日志记录
	log := &FindGuidLog{
		TargetPath:      targetPath,
		NotificationURL: notificationURL,
		Status:          0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := SaveFindGuidLog(log); err != nil {
		return nil, fmt.Errorf("创建日志记录失败: %v", err)
	}

	// 使用map存储GUID和文件路径的映射
	guidMap := make(map[string][]string)
	var mutex sync.Mutex

	// 遍历目录
	err := filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 只处理.meta文件
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".meta") {
			// 读取文件内容
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			reader := bufio.NewReader(file)
			for {
				line, err := reader.ReadString('\n')
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}

				// 查找GUID
				matches := guidRegex.FindStringSubmatch(line)
				if len(matches) > 1 {
					guid := matches[1]
					mutex.Lock()
					guidMap[guid] = append(guidMap[guid], path)
					mutex.Unlock()
				}
			}
		}
		return nil
	})

	if err != nil {
		log.Status = 2
		log.Message = fmt.Sprintf("查找GUID失败: %v", err)
		log.UpdatedAt = time.Now()
		UpdateFindGuidLog(log.ID, map[string]interface{}{
			"status":     log.Status,
			"message":    log.Message,
			"updated_at": log.UpdatedAt,
		})
		return log, err
	}

	// 找出重复的GUID
	var duplicates []DuplicateGuid
	for guid, paths := range guidMap {
		if len(paths) > 1 {
			for _, path := range paths {
				duplicates = append(duplicates, DuplicateGuid{
					LogID:     log.ID,
					GUID:      guid,
					FilePath:  path,
					CreatedAt: time.Now(),
				})
			}
		}
	}

	// 保存重复的GUID记录
	if len(duplicates) > 0 {
		if err := SaveDuplicateGuids(duplicates); err != nil {
			log.Status = 2
			log.Message = fmt.Sprintf("保存重复GUID记录失败: %v", err)
			log.UpdatedAt = time.Now()
			UpdateFindGuidLog(log.ID, map[string]interface{}{
				"status":     log.Status,
				"message":    log.Message,
				"updated_at": log.UpdatedAt,
			})
			return log, err
		}
	}

	// 发送通知
	notificationData := map[string]interface{}{
		"id":  log.ID,
		"msg": duplicates,
	}
	jsonData, err := json.Marshal(notificationData)
	if err != nil {
		log.Status = 2
		log.Message = fmt.Sprintf("准备通知数据失败: %v", err)
		log.UpdatedAt = time.Now()
		UpdateFindGuidLog(log.ID, map[string]interface{}{
			"status":     log.Status,
			"message":    log.Message,
			"updated_at": log.UpdatedAt,
		})
		return log, err
	}

	// 发送HTTP请求
	resp, err := http.Post(notificationURL, "application/json", strings.NewReader(string(jsonData)))
	if err != nil {
		log.Status = 2
		log.Message = fmt.Sprintf("发送通知失败: %v", err)
		log.UpdatedAt = time.Now()
		UpdateFindGuidLog(log.ID, map[string]interface{}{
			"status":     log.Status,
			"message":    log.Message,
			"updated_at": log.UpdatedAt,
		})
		return log, err
	}
	defer resp.Body.Close()

	// 更新日志状态
	log.Status = 1
	log.DuplicateCount = len(duplicates)
	log.Message = fmt.Sprintf("成功找到 %d 个重复的GUID", len(duplicates))
	log.UpdatedAt = time.Now()
	UpdateFindGuidLog(log.ID, map[string]interface{}{
		"status":          log.Status,
		"duplicate_count": log.DuplicateCount,
		"message":         log.Message,
		"updated_at":      log.UpdatedAt,
	})

	return log, nil
}

// GetFindGuidLogs 获取查找GUID的日志列表
func GetFindGuidLogs(page, limit int) ([]FindGuidLog, int64, error) {
	return GetFindGuidLogsFromDB(page, limit)
}

// GetDuplicateGuids 获取指定日志ID的重复GUID列表
func GetDuplicateGuids(logID uint) ([]DuplicateGuid, error) {
	return GetDuplicateGuidsByLogID(logID)
}
