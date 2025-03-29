package citask

import (
	"fmt"
	"sync"
	"time"

	"github.com/andycai/unitool/models"
	"github.com/robfig/cron/v3"
)

var (
	cronEntries   = make(map[uint]cron.EntryID)
	progressMutex sync.RWMutex
	cronScheduler *cron.Cron
)

// 初始化定时任务调度器
func initCron() {
	cronScheduler = cron.New()
	cronScheduler.Start()

	// 从数据库加载定时任务
	var tasks []models.Task
	if err := app.DB.Where("enable_cron = ? AND status = ?", true, "active").Find(&tasks).Error; err != nil {
		fmt.Printf("加载定时任务失败: %v\n", err)
		return
	}

	for _, task := range tasks {
		if err := scheduleCronTask(&task); err != nil {
			fmt.Printf("调度任务失败 [%d]: %v\n", task.ID, err)
		} else {
			fmt.Printf("成功加载定时任务 [%d]: %s\n", task.ID, task.Name)
		}
	}
}

// 调度定时任务
func scheduleCronTask(task *models.Task) error {
	if task.EnableCron == 0 || task.CronExpr == "" {
		return nil
	}

	entryID, err := cronScheduler.AddFunc(task.CronExpr, func() {
		// 创建任务日志
		taskLog := &models.TaskLog{
			TaskID:    task.ID,
			StartTime: time.Now(),
			Status:    "running",
		}

		if err := app.DB.Create(taskLog).Error; err != nil {
			fmt.Printf("创建任务日志失败: %v\n", err)
			return
		}

		// 创建进度信息
		progress := &TaskProgress{
			Status:    "running",
			StartTime: time.Now(),
		}

		progressMutex.Lock()
		taskProgressMap[taskLog.ID] = progress
		progressMutex.Unlock()

		// 执行任务
		go func() {
			if task.Type == "script" {
				executeScriptTask(task, taskLog, progress)
			} else {
				executeHTTPTask(task, taskLog, progress)
			}

			// 更新任务日志
			taskLog.EndTime = time.Now()
			taskLog.Status = progress.Status
			taskLog.Output = progress.Output
			taskLog.Error = progress.Error
			if err := app.DB.Save(taskLog).Error; err != nil {
				fmt.Printf("更新任务日志失败: %v\n", err)
			}
		}()
	})

	if err != nil {
		return err
	}

	// 保存定时任务ID
	progressMutex.Lock()
	cronEntries[task.ID] = entryID
	progressMutex.Unlock()

	return nil
}
