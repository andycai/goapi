package citask

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/modules/system/adminlog"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"golang.org/x/text/encoding/simplifiedchinese"
)

var (
	taskProgressMap = make(map[uint]*TaskProgress)
	taskCmdMap      = make(map[uint]*exec.Cmd)
)

// getTasks 获取任务列表
func listTasksHandler(c *fiber.Ctx) error {
	var tasks []models.Task
	if err := app.DB.Order("created_at desc").Find(&tasks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取任务列表失败: %v", err),
		})
	}
	return c.JSON(tasks)
}

// createTaskHandler 创建任务
func createTaskHandler(c *fiber.Ctx) error {
	var task models.Task
	if err := c.BodyParser(&task); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 如果启用了定时执行，验证cron表达式
	if task.EnableCron == 1 {
		if task.CronExpr == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "启用定时执行时必须提供Cron表达式",
			})
		}
		if _, err := cron.ParseStandard(task.CronExpr); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("无效的Cron表达式: %v", err),
			})
		}
	}

	if err := app.DB.Create(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建任务失败: %v", err),
		})
	}

	// 如果启用了定时执行，添加到调度器
	if task.EnableCron == 1 {
		if err := scheduleCronTask(&task); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": fmt.Sprintf("设置定时任务失败: %v", err),
			})
		}
	}

	return c.JSON(task)
}

// getTaskHandler 获取任务详情
func getTaskHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := app.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}
	return c.JSON(task)
}

// updateTaskHandler 更新任务
func updateTaskHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少任务ID",
		})
	}

	var updates models.Task
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的请求数据: %v", err),
		})
	}

	// 获取原有任务信息
	var task models.Task
	if err := app.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}

	// 检查定时任务状态变化
	cronChanged := task.EnableCron != updates.EnableCron ||
		(updates.EnableCron == 1 && task.CronExpr != updates.CronExpr)

	// 如果启用了定时执行，验证cron表达式
	if updates.EnableCron == 1 {
		if updates.CronExpr == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "启用定时执行时必须提供Cron表达式",
			})
		}
		if _, err := cron.ParseStandard(updates.CronExpr); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": fmt.Sprintf("无效的Cron表达式: %v", err),
			})
		}
	}

	// 更新任务信息
	if err := app.DB.Model(&task).Updates(updates).Update("enable_cron", updates.EnableCron).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("更新任务失败: %v", err),
		})
	}

	// 如果定时配置发生变化，处理定时任务
	if cronChanged {
		progressMutex.Lock()
		// 如果存在旧的定时任务，先移除
		if entryID, ok := cronEntries[task.ID]; ok {
			cronScheduler.Remove(entryID)
			delete(cronEntries, task.ID)
			fmt.Printf("已移除任务 [%d] 的定时配置\n", task.ID)
		}
		progressMutex.Unlock()

		// 如果启用了定时执行，添加新的定时任务
		if updates.EnableCron == 1 {
			if err := scheduleCronTask(&task); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": fmt.Sprintf("更新定时任务失败: %v", err),
				})
			}
			fmt.Printf("已为任务 [%d] 添加新的定时配置: %s\n", task.ID, updates.CronExpr)
		}
	}

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": task,
	})
}

// deleteTaskHandler 删除任务
func deleteTaskHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := app.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}

	if err := app.DB.Delete(&task).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("删除任务失败: %v", err),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "task", task.ID, fmt.Sprintf("删除任务：%s", task.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// runTaskHandler 执行任务
func runTaskHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var task models.Task
	if err := app.DB.First(&task, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("任务不存在: %v", err),
		})
	}

	// 创建任务日志
	taskLog := models.TaskLog{
		TaskID:    task.ID,
		Status:    "running",
		StartTime: time.Now(),
	}
	if err := app.DB.Create(&taskLog).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("创建任务日志失败: %v", err),
		})
	}

	// 初始化进度信息
	progress := &TaskProgress{
		ID:        taskLog.ID,
		TaskID:    task.ID,
		TaskName:  task.Name,
		Status:    "running",
		StartTime: taskLog.StartTime,
		Progress:  0,
	}
	progressMutex.Lock()
	taskProgressMap[taskLog.ID] = progress
	progressMutex.Unlock()

	// 记录操作日志
	adminlog.WriteLog(c, "run", "task", task.ID, fmt.Sprintf("执行任务：%s", task.Name))

	// 异步执行任务
	go executeTask(&task, &taskLog)

	return c.JSON(taskLog)
}

// getTaskLogsHandler 获取任务日志
func getTaskLogsHandler(c *fiber.Ctx) error {
	taskID := c.Params("id")
	var logs []models.TaskLog
	if err := app.DB.Where("task_id = ?", taskID).Order("created_at desc").Find(&logs).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("获取任务日志失败: %v", err),
		})
	}
	return c.JSON(logs)
}

// getTaskStatus 获取任务状态
func getTaskStatus(c *fiber.Ctx) error {
	logID := c.Query("log_id")
	var log models.TaskLog
	if err := app.DB.First(&log, logID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": fmt.Sprintf("日志不存在: %v", err),
		})
	}
	return c.JSON(log)
}

// getTaskProgressHandler 获取任务进度
func getTaskProgressHandler(c *fiber.Ctx) error {
	logId := c.Params("logId")
	if logId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少logId参数",
		})
	}

	// 将字符串转换为uint
	id, err := strconv.ParseUint(logId, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的logId参数",
		})
	}

	progressMutex.RLock()
	progress, exists := taskProgressMap[uint(id)]
	progressMutex.RUnlock()

	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "找不到任务进度信息",
		})
	}

	return c.JSON(progress)
}

// executeTask 执行任务
func executeTask(task *models.Task, log *models.TaskLog) {
	progress := taskProgressMap[log.ID]
	defer func() {
		log.EndTime = time.Now()
		log.Duration = int(log.EndTime.Sub(log.StartTime).Seconds())

		// 执行完成任务，保存任务日志到数据库
		app.DB.Save(log)

		// 更新并清理进度信息
		if progress != nil {
			progress.Status = log.Status
			progress.EndTime = log.EndTime
			progress.Duration = log.Duration
			progress.Progress = 100

			// 延迟删除进度信息
			time.AfterFunc(time.Hour*2, func() {
				progressMutex.Lock()
				delete(taskProgressMap, log.ID)
				progressMutex.Unlock()
			})
		}
	}()

	switch task.Type {
	case "script":
		executeScriptTask(task, log, progress)
	case "http":
		executeHTTPTask(task, log, progress)
	default:
		log.Status = "failed"
		log.Error = "未知的任务类型"
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
	}
}

// isUnsafeCommand 检查命令是否不安全
func isUnsafeCommand(script string) (bool, string) {
	// 转换为小写以进行大小写不敏感的检查
	lowerScript := strings.ToLower(script)

	// 通用的危险模式（适用于所有操作系统）
	dangerousPatterns := []string{
		"$(", "`", // 命令替换
		"&&", "||", // 命令链接
		"../", // 目录遍历
		"/*",  // 根目录操作
	}

	// 检查危险的模式
	for _, pattern := range dangerousPatterns {
		if strings.Contains(script, pattern) {
			return true, fmt.Sprintf("检测到危险的命令模式: %s", pattern)
		}
	}

	// 根据操作系统选择不安全命令列表
	var unsafeCommands []string
	if runtime.GOOS == "windows" {
		// Windows 特定的危险命令
		unsafeCommands = []string{
			"format",        // 格式化磁盘
			"del /", "rd /", // 删除系统文件
			"net user",            // 用户管理
			"net localgroup",      // 用户组管理
			"reg",                 // 注册表操作
			"taskkill /im",        // 结束系统进程
			"shutdown", "restart", // 系统关机重启
			"netsh",      // 网络配置
			"sc",         // 服务控制
			"powershell", // PowerShell 命令
			"wmic",       // Windows 管理工具
			"cipher",     // 文件加密
			"diskpart",   // 磁盘分区
			"chkdsk",     // 磁盘检查
			"attrib",     // 文件属性修改
			"runas",      // 提权运行
		}
	} else {
		// Unix/Linux 特定的危险命令
		unsafeCommands = []string{
			"rm -rf", "rm -r", // 递归删除
			"mkfs",                         // 格式化
			":(){:|:&};:", ":(){ :|:& };:", // Fork炸弹
			"dd",                // 磁盘操作
			"> /dev/", ">/dev/", // 设备文件操作
			"wget", "curl", // 外部下载
			"chmod 777", "chmod -R 777", // 危险的权限设置
			"sudo", "su", // 提权命令
			"nc", "netcat", // 网络工具
			"telnet",          // 远程连接
			"|mail", "|email", // 邮件命令
			"tcpdump",        // 网络抓包
			"chown -R",       // 递归改变所有者
			"mv /* ", "mv /", // 移动根目录
			"cp /* ", "cp /", // 复制根目录
			"shutdown", "reboot", "halt", "poweroff", // 系统关机重启
			"passwd",             // 修改密码
			"useradd", "userdel", // 用户管理
			"mkfs", "fdisk", "fsck", // 磁盘管理
			"iptables", "firewall", // 防火墙
			"nmap", // 端口扫描
			"eval", // 命令注入
		}
	}

	// 检查不安全的命令
	for _, cmd := range unsafeCommands {
		if strings.Contains(lowerScript, cmd) {
			return true, fmt.Sprintf("检测到不安全的命令: %s", cmd)
		}
	}

	// 检查环境变量操作
	if strings.Contains(lowerScript, "export") || strings.Contains(lowerScript, "env") ||
		strings.Contains(lowerScript, "set ") { // Windows 环境变量设置
		return true, "不允许修改环境变量"
	}

	// 检查系统特定的危险重定向
	if runtime.GOOS == "windows" {
		windowsPatterns := []string{
			"> con:", ">con:", // Windows 设置文件
			"> prn:", ">prn:",
			"> aux:", ">aux:",
			"con:", "prn:", "aux:", // 设备名称
			"> nul", ">nul", // Windows null 设备
		}
		for _, pattern := range windowsPatterns {
			if strings.Contains(lowerScript, pattern) {
				return true, fmt.Sprintf("检测到危险的 Windows 设备操作: %s", pattern)
			}
		}
	} else {
		unixPatterns := []string{
			"> /", ">/", // 重定向到系统目录
			"2> /", "2>/", // 错误重定向到系统目录
			">> /", ">>/", // 追加重定向到系统目录
			"< /", "</", // 从系统目录读取
		}
		for _, pattern := range unixPatterns {
			if strings.Contains(script, pattern) {
				return true, fmt.Sprintf("检测到危险的文件重定向: %s", pattern)
			}
		}
	}

	return false, ""
}

// 添加 GBK 输出处理器
type gbkOutputWriter struct {
	buffer *bytes.Buffer
}

func (w *gbkOutputWriter) Write(p []byte) (n int, err error) {
	// 尝试将 GBK 转换为 UTF-8
	utf8Bytes, err := simplifiedchinese.GB18030.NewDecoder().Bytes(p)
	if err != nil {
		// 如果转换失败，尝试使用 GBK 解码
		if gbkBytes, err := simplifiedchinese.GBK.NewDecoder().Bytes(p); err == nil {
			return w.buffer.Write(gbkBytes)
		}
		// 如果 GBK 也失败，尝试使用 HZ-GB2312 解码
		if hzBytes, err := simplifiedchinese.HZGB2312.NewDecoder().Bytes(p); err == nil {
			return w.buffer.Write(hzBytes)
		}
		// 所有转换都失败，直接写入原始数据
		return w.buffer.Write(p)
	}
	// 写入转换后的 UTF-8 数据
	return w.buffer.Write(utf8Bytes)
}

// executeScriptTask 执行脚本任务
func executeScriptTask(task *models.Task, log *models.TaskLog, progress *TaskProgress) {
	fmt.Printf("开始执行脚本任务: %s (ID: %d)\n", task.Name, task.ID)

	// 首先检查脚本安全性
	if unsafe, reason := isUnsafeCommand(task.Script); unsafe {
		log.Status = "failed"
		log.Error = fmt.Sprintf("脚本包含不安全的命令: %s", reason)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("脚本安全检查失败: %s\n", reason)
		return
	}

	// 创建临时脚本文件
	ext := ".sh"
	if runtime.GOOS == "windows" {
		ext = ".bat"
	}

	tmpFile, err := os.CreateTemp("", "task_*"+ext)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建临时文件失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("创建临时文件失败: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())
	fmt.Printf("创建临时脚本文件: %s\n", tmpFile.Name())

	// 添加安全限制的shell选项（仅用于Unix系统）
	scriptContent := task.Script
	if runtime.GOOS != "windows" {
		scriptContent = "set -euo pipefail\ntrap 'exit 1' INT TERM\n" + scriptContent
	}

	if _, err := tmpFile.WriteString(scriptContent); err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("写入脚本失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("写入脚本内容失败: %v\n", err)
		return
	}
	tmpFile.Close()
	fmt.Printf("写入脚本内容:\n%s\n", scriptContent)

	// 设置脚本可执行权限
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tmpFile.Name(), 0755); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("设置脚本权限失败: %v", err)
			if progress != nil {
				progress.Status = "failed"
				progress.Error = log.Error
			}
			fmt.Printf("设置脚本权限失败: %v\n", err)
			return
		}
		fmt.Println("设置脚本可执行权限成功")
	}

	// 执行脚本
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		// Windows 下���置代码页和环境变量
		// 使用 /V:ON 启用延迟变量扩展，/D 禁用自动运行命令，/S 禁用命令回显
		// < NUL 防止等待输入，> NUL 2>&1 重定向标准输出和错误输出
		cmdStr := fmt.Sprintf("cmd /V:ON /D /S /C \"chcp 65001 > NUL && type \"%s\" > NUL && \"%s\" < NUL\"", tmpFile.Name(), tmpFile.Name())
		cmd = exec.Command("cmd", "/C", cmdStr)

		// 设置环境变量
		cmd.Env = append(os.Environ(),
			"PYTHONIOENCODING=utf8",
			"PYTHONLEGACYWINDOWSSTDIO=1",
			"PYTHONUNBUFFERED=1",
			"JAVA_TOOL_OPTIONS=-Dfile.encoding=UTF-8",
		)
		fmt.Printf("Windows 命令: %s\n", cmdStr)
	} else {
		cmd = exec.Command("/bin/bash", tmpFile.Name())
		fmt.Printf("Unix 命令: /bin/bash %s\n", tmpFile.Name())
	}

	// 设置工作目录为临时目录
	tmpDir := os.TempDir()
	cmd.Dir = tmpDir
	fmt.Printf("工作目录: %s\n", tmpDir)

	// 创建输出缓冲区
	var outputBuffer bytes.Buffer
	var errorBuffer bytes.Buffer

	// 设置标准输入为 null 设备
	if runtime.GOOS == "windows" {
		nullFile, err := os.OpenFile("NUL", os.O_RDWR, 0)
		if err == nil {
			cmd.Stdin = nullFile
			defer nullFile.Close()
		}
	} else {
		nullFile, err := os.OpenFile("/dev/null", os.O_RDWR, 0)
		if err == nil {
			cmd.Stdin = nullFile
			defer nullFile.Close()
		}
	}

	cmd.Stdout = io.MultiWriter(&outputBuffer, os.Stdout)
	cmd.Stderr = io.MultiWriter(&errorBuffer, os.Stderr)

	// 设置超时
	timeout := time.Duration(task.Timeout) * time.Second
	if timeout == 0 {
		timeout = 300 * time.Second // 默认5分钟超时
	}
	fmt.Printf("设置超时时间: %v\n", timeout)

	// 创建一个带有超时的context
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 使用context创建新的命令
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	cmd.Dir = tmpDir
	cmd.Stdout = io.MultiWriter(&outputBuffer, os.Stdout)
	cmd.Stderr = io.MultiWriter(&errorBuffer, os.Stderr)

	// 保存命令到映射中
	progressMutex.Lock()
	taskCmdMap[log.ID] = cmd
	progressMutex.Unlock()

	// 启动命令
	if err := cmd.Start(); err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("启动命令失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("启动命令失败: %v\n", err)
		return
	}
	fmt.Println("命令启动成功")

	// 等待命令完成
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	// 清理函数
	cleanup := func() {
		progressMutex.Lock()
		delete(taskCmdMap, task.ID)
		progressMutex.Unlock()
	}
	defer cleanup()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				log.Status = "failed"
				log.Error = fmt.Sprintf("执行超时（%d秒）\n%s\n%s", task.Timeout, outputBuffer.String(), errorBuffer.String())
				if progress != nil {
					progress.Status = "failed"
					progress.Error = log.Error
					progress.Output = outputBuffer.String() + "\nError: " + errorBuffer.String()
				}
				fmt.Printf("命令执行超时: %v\n", ctx.Err())
				return
			}
		case err := <-done:
			output := outputBuffer.String()
			errorOutput := errorBuffer.String()

			if err != nil {
				log.Status = "failed"
				log.Error = fmt.Sprintf("执行失败: %v\n%s\n%s", err, output, errorOutput)
				if progress != nil {
					progress.Status = "failed"
					progress.Error = log.Error
					progress.Output = output + "\nError: " + errorOutput
				}
				fmt.Printf("命令执行失败: %v\n", err)
				return
			}

			// 更新任务状态和输出
			log.Status = "success"
			log.Output = outputBuffer.String()
			log.Error = errorBuffer.String()

			if progress != nil {
				progress.Status = "success"
				progress.Output = log.Output
				progress.Progress = 100
			}
			fmt.Printf("任务执行成功完成: %s (ID: %d)\n", task.Name, task.ID)
			return
		case <-ticker.C:
			if progress != nil {
				progress.Output = outputBuffer.String()
				if errorBuffer.Len() > 0 {
					progress.Output += "\nError: " + errorBuffer.String()
				}
			}
		}
	}
}

// executeHTTPTask 执行HTTP任务
func executeHTTPTask(task *models.Task, log *models.TaskLog, progress *TaskProgress) {
	fmt.Printf("开始执行HTTP任务: %s (ID: %d)\n", task.Name, task.ID)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: time.Duration(task.Timeout) * time.Second,
	}

	// 创建请求
	var body io.Reader
	if task.Body != "" {
		body = strings.NewReader(task.Body)
	}
	req, err := http.NewRequest(task.Method, task.URL, body)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("创建请求失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("创建HTTP请求失败: %v\n", err)
		return
	}

	// 添加请求头
	if task.Headers != "" {
		var headers map[string]string
		if err := json.Unmarshal([]byte(task.Headers), &headers); err != nil {
			log.Status = "failed"
			log.Error = fmt.Sprintf("解析请求头失败: %v", err)
			if progress != nil {
				progress.Status = "failed"
				progress.Error = log.Error
			}
			fmt.Printf("解析请求头失败: %v\n", err)
			return
		}
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("发送请求失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("发送HTTP请求失败: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Status = "failed"
		log.Error = fmt.Sprintf("读取响应失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("读取HTTP响应失败: %v\n", err)
		return
	}

	// 仅在 Windows 系统下尝试转换响应内容的编码
	if runtime.GOOS == "windows" {
		gbkOutput := &gbkOutputWriter{buffer: &bytes.Buffer{}}
		if _, err := gbkOutput.Write(respBody); err == nil && gbkOutput.buffer.Len() > 0 {
			respBody = gbkOutput.buffer.Bytes()
		}
	}

	// 检查响应状态码
	if resp.StatusCode >= 400 {
		log.Status = "failed"
		log.Error = fmt.Sprintf("HTTP请求失败: %s\n响应内容: %s", resp.Status, string(respBody))
		if progress != nil {
			progress.Status = "failed"
			progress.Error = log.Error
		}
		fmt.Printf("HTTP请求返回错误状态码: %d\n", resp.StatusCode)
		return
	}

	// 更新任务状态
	log.Status = "success"
	log.Output = string(respBody)
	if progress != nil {
		progress.Status = "success"
		progress.Output = log.Output
		progress.Progress = 100
	}
	fmt.Printf("HTTP任务执行成功完成: %s (ID: %d)\n", task.Name, task.ID)
}

// stopTaskHandler 停止正在执行的任务
func stopTaskHandler(c *fiber.Ctx) error {
	logId := c.Params("logId")
	if logId == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少日志ID",
		})
	}

	// 将字符串ID转换为uint
	id, err := strconv.ParseUint(logId, 10, 32)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的日志ID",
		})
	}
	taskId := uint(id)

	// 获取进度信息
	progressMutex.Lock()
	progress, exists := taskProgressMap[taskId]
	progressMutex.Unlock()

	if !exists {
		return c.Status(404).JSON(fiber.Map{
			"error": "任务不存在或已结束",
		})
	}

	// 如果任务不是运行状态，返回错误
	if progress.Status != "running" {
		return c.Status(400).JSON(fiber.Map{
			"error": "任务不在运行状态",
		})
	}

	// 获取命令进程
	progressMutex.Lock()
	cmd := taskCmdMap[taskId]
	progressMutex.Unlock()

	if cmd == nil || cmd.Process == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "无法获取任务进程",
		})
	}

	// 停止进程
	if runtime.GOOS == "windows" {
		// Windows 下使用 taskkill 强制结束进程树
		exec.Command("taskkill", "/F", "/T", "/PID", strconv.Itoa(cmd.Process.Pid)).Run()
	} else {
		// Unix 系统下发送 SIGTERM 信号
		err = cmd.Process.Signal(syscall.SIGTERM)
		if err != nil {
			// 如果 SIGTERM 失败，尝试 SIGKILL
			err = cmd.Process.Kill()
		}
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("停止任务失败: %v", err),
		})
	}

	// 更新任务状态
	progressMutex.Lock()
	progress.Status = "failed"
	progress.Error = "任务被手动停止"
	progress.EndTime = time.Now()
	progressMutex.Unlock()

	// 清理命令映射
	progressMutex.Lock()
	delete(taskCmdMap, taskId)
	progressMutex.Unlock()

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "任务已停止",
	})
}

// listRunningTasksHandler 获取正在执行的任务列表
func listRunningTasksHandler(c *fiber.Ctx) error {
	progressMutex.Lock()
	defer progressMutex.Unlock()

	// 从内存中获取所有正在执行的任务
	var runningTasks []fiber.Map
	for id, progress := range taskProgressMap {
		if progress.Status == "running" {
			// 查询任务信息
			var taskLog models.TaskLog
			if err := app.DB.First(&taskLog, id).Error; err != nil {
				continue
			}

			// 构建返回数据
			runningTasks = append(runningTasks, fiber.Map{
				"id":         id,
				"name":       progress.TaskName,
				"status":     progress.Status,
				"progress":   progress.Progress,
				"output":     progress.Output,
				"error":      progress.Error,
				"start_time": progress.StartTime.Unix(),
			})
		}
	}

	// 按开始时间倒序排序
	sort.Slice(runningTasks, func(i, j int) bool {
		return runningTasks[i]["start_time"].(int64) > runningTasks[j]["start_time"].(int64)
	})

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": runningTasks,
	})
}

// getNextRunTimeHandler 计算下次执行时间
func getNextRunTimeHandler(c *fiber.Ctx) error {
	expr := c.Query("expr")
	if expr == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "缺少cron表达式",
		})
	}

	// 解析cron表达式
	schedule, err := cron.ParseStandard(expr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": fmt.Sprintf("无效的cron表达式: %v", err),
		})
	}

	// 计算下次执行时间
	nextTime := schedule.Next(time.Now())

	return c.JSON(fiber.Map{
		"code": 0,
		"msg":  "success",
		"data": fiber.Map{
			"next_run":      nextTime.Unix(),
			"next_run_text": nextTime.Format("2006-01-02 15:04:05"),
		},
	})
}

// searchTasksHandler 搜索任务列表
func searchTasksHandler(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	if keyword == "" {
		return c.JSON([]models.Task{})
	}

	var tasks []models.Task
	if err := app.DB.Where("name LIKE ?", "%"+keyword+"%").
		Order("created_at desc").
		Limit(10).
		Find(&tasks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": fmt.Sprintf("搜索任务失败: %v", err),
		})
	}
	return c.JSON(tasks)
}
