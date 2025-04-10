package luban

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/andycai/goapi/models"
)

// executeExport 执行导出任务
func executeExport(export *models.ConfigExport, project *models.ConfigProject, table *models.ConfigTable) {
	progress := exportProgressMap[export.ID]
	defer func() {
		export.EndTime = time.Now()
		export.Duration = int(export.EndTime.Sub(export.StartTime).Seconds())

		// 执行完成任务，保存导出记录到数据库
		app.DB.Save(export)

		// 更新并清理进度信息
		if progress != nil {
			progress.Status = export.Status
			progress.EndTime = export.EndTime
			progress.Duration = export.Duration
			progress.Progress = 100

			// 延迟删除进度信息
			time.AfterFunc(time.Hour*2, func() {
				delete(exportProgressMap, export.ID)
			})
		}
	}()

	// 创建输出目录
	outputDir := filepath.Join(project.OutputPath, table.Name)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		export.Status = "failed"
		export.Error = fmt.Sprintf("创建输出目录失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = export.Error
		}
		return
	}

	// 准备Luban命令参数
	args := []string{
		"--input_data_dir", project.RootPath,
		"--output_data_dir", outputDir,
		"--gen_types", export.Format,
		"--target_language", export.Language,
	}

	// 根据文件类型添加特定参数
	switch table.FileType {
	case "excel", "csv":
		args = append(args, "--input", table.FilePath)
		if table.SheetName != "" {
			args = append(args, "--sheet_name", table.SheetName)
		}
	case "json":
		args = append(args, "--json", table.FilePath)
	case "xml":
		args = append(args, "--xml", table.FilePath)
	case "yaml":
		args = append(args, "--yaml", table.FilePath)
	default:
		export.Status = "failed"
		export.Error = fmt.Sprintf("不支持的文件类型: %s", table.FileType)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = export.Error
		}
		return
	}

	// 如果有数据验证规则，添加验证参数
	if table.Validators != "" {
		var validators map[string]interface{}
		if err := json.Unmarshal([]byte(table.Validators), &validators); err == nil {
			validatorFile := filepath.Join(os.TempDir(), fmt.Sprintf("validator_%d.json", export.ID))
			if err := os.WriteFile(validatorFile, []byte(table.Validators), 0644); err == nil {
				args = append(args, "--validator", validatorFile)
				defer os.Remove(validatorFile)
			}
		}
	}

	// 创建命令
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("luban.exe", args...)
	} else {
		cmd = exec.Command("luban", args...)
	}

	// 设置工作目录
	cmd.Dir = project.RootPath

	// 创建输出缓冲区
	var outputBuffer bytes.Buffer
	var errorBuffer bytes.Buffer
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errorBuffer

	// 设置超时
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
	cmd.Dir = project.RootPath
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errorBuffer

	// 启动命令
	if err := cmd.Start(); err != nil {
		export.Status = "failed"
		export.Error = fmt.Sprintf("启动导出命令失败: %v", err)
		if progress != nil {
			progress.Status = "failed"
			progress.Error = export.Error
		}
		return
	}

	// 等待命令完成
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
	}()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				export.Status = "failed"
				export.Error = fmt.Sprintf("导出超时（30分钟）\n%s\n%s", outputBuffer.String(), errorBuffer.String())
				if progress != nil {
					progress.Status = "failed"
					progress.Error = export.Error
					progress.Output = outputBuffer.String() + "\nError: " + errorBuffer.String()
				}
				return
			}
		case err := <-done:
			output := outputBuffer.String()
			errorOutput := errorBuffer.String()

			if err != nil {
				export.Status = "failed"
				export.Error = fmt.Sprintf("导出失败: %v\n%s\n%s", err, output, errorOutput)
				if progress != nil {
					progress.Status = "failed"
					progress.Error = export.Error
					progress.Output = output + "\nError: " + errorOutput
				}
				return
			}

			// 检查输出文件是否生成
			outputFiles, err := filepath.Glob(filepath.Join(outputDir, "*"))
			if err != nil || len(outputFiles) == 0 {
				export.Status = "failed"
				export.Error = fmt.Sprintf("未找到导出文件\n%s\n%s", output, errorOutput)
				if progress != nil {
					progress.Status = "failed"
					progress.Error = export.Error
					progress.Output = output + "\nError: " + errorOutput
				}
				return
			}

			// 更新导出状态和输出
			export.Status = "success"
			export.Output = strings.Join([]string{
				fmt.Sprintf("导出文件目录: %s", outputDir),
				fmt.Sprintf("导出文件列表: %s", strings.Join(outputFiles, ", ")),
				"导出日志:",
				output,
			}, "\n")

			if progress != nil {
				progress.Status = "success"
				progress.Output = export.Output
				progress.Progress = 100
			}
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
