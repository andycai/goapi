package browse

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/andycai/unitool/modules/system/adminlog"
	"github.com/gofiber/fiber/v2"
	"github.com/jlaffaye/ftp"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/ianaindex"
)

// WriteToBinaryFile 写入二进制文件
func WriteToBinaryFile(filename, username, password string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 将用户名和密码进行 Base64 编码
	usernameBase64 := base64.StdEncoding.EncodeToString([]byte(username))
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))

	// 写入用户名长度和内容
	usernameLength := uint32(len(usernameBase64))
	if err := binary.Write(file, binary.LittleEndian, usernameLength); err != nil {
		return err
	}
	if _, err := file.Write([]byte(usernameBase64)); err != nil {
		return err
	}

	// 写入密码长度和内容
	passwordLength := uint32(len(passwordBase64))
	if err := binary.Write(file, binary.LittleEndian, passwordLength); err != nil {
		return err
	}
	if _, err := file.Write([]byte(passwordBase64)); err != nil {
		return err
	}

	return nil
}

// ReadFromBinaryFile 从二进制文件中读取用户名和密码
func ReadFromBinaryFile(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	// 读取用户名
	var usernameLength uint32
	if err := binary.Read(file, binary.LittleEndian, &usernameLength); err != nil {
		return "", "", err
	}
	usernameBase64 := make([]byte, usernameLength)
	if _, err := file.Read(usernameBase64); err != nil {
		return "", "", err
	}
	usernameBytes, err := base64.StdEncoding.DecodeString(string(usernameBase64))
	if err != nil {
		return "", "", err
	}
	username := string(usernameBytes)

	// 读取密码
	var passwordLength uint32
	if err := binary.Read(file, binary.LittleEndian, &passwordLength); err != nil {
		return "", "", err
	}
	passwordBase64 := make([]byte, passwordLength)
	if _, err := file.Read(passwordBase64); err != nil {
		return "", "", err
	}
	passwordBytes, err := base64.StdEncoding.DecodeString(string(passwordBase64))
	if err != nil {
		return "", "", err
	}
	password := string(passwordBytes)

	return username, password, nil
}

// 规范化路径分隔符，将反斜杠转换为正斜杠
func normalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

// formatFileSize 格式化文件大小
func formatFileSize(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%d B", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// handleBrowseDirectory 处理目录浏览请求
func handleBrowseDirectory(c *fiber.Ctx, path string) error {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	var entries []FileEntry
	for _, file := range files {
		// 获取文件类型
		fileType := ""
		if !file.IsDir() {
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".apk" {
				fileType = "apk"
			} else if ext == ".zip" {
				fileType = "zip"
			}
		}

		entries = append(entries, FileEntry{
			Name:         file.Name(),
			Size:         file.Size(),
			FormatedSize: formatFileSize(file.Size()),
			ModTime:      file.ModTime(),
			IsDir:        file.IsDir(),
			FileType:     fileType,
		})
	}

	// 按照文件夹在前，文件在后的顺序排序
	// 同类型按修改时间倒序排序
	sort.Slice(entries, func(i, j int) bool {
		// 如果一个是目录一个是文件，目录排在前面
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir
		}
		// 如果都是目录或都是文件，按修改时间倒序排序
		return entries[i].ModTime.After(entries[j].ModTime)
	})

	// 获取相对于根目录的路径
	rootDir := app.Config.Server.Output
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return c.Status(500).SendString("Error resolving root path")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return c.Status(500).SendString("Error resolving directory path")
	}

	relPath, err := filepath.Rel(absRootDir, absPath)
	if err != nil {
		return c.Status(500).SendString("Error calculating relative path")
	}

	if relPath == "." {
		relPath = ""
	}

	// 规范化路径分隔符
	relPath = normalizePath(relPath)

	rootPath := "/admin/browse"

	return c.Render("admin/directory", fiber.Map{
		"Title":    "目录浏览",
		"RootPath": rootPath,
		"Path":     relPath,
		"Entries":  entries,
		"Scripts": []string{
			"/static/js/admin/directory.js",
		},
	}, "admin/layout")
}

// isTextFile 检查文件是否为文本文件
func isTextFile(filename string) bool {
	textExtensions := map[string]bool{
		".txt":  true,
		".log":  true,
		".json": true,
		".xml":  true,
		".yml":  true,
		".yaml": true,
		".md":   true,
		".ini":  true,
		".conf": true,
		".cfg":  true,
		".html": true,
		".css":  true,
		".js":   true,
		".ts":   true,
		".go":   true,
		".py":   true,
		".java": true,
		".php":  true,
		".sh":   true,
		".bat":  true,
		".ps1":  true,
		".sql":  true,
		".csv":  true,
	}
	ext := strings.ToLower(filepath.Ext(filename))
	return textExtensions[ext]
}

// handleBrowseFile 处理文件浏览请求
func handleBrowseFile(c *fiber.Ctx, path string) error {
	if path == "" {
		return fiber.NewError(fiber.StatusBadRequest, "文件路径不能为空")
	}

	// 获取相对于根目录的路径
	rootDir := app.Config.Server.Output
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "获取根目录失败")
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "获取文件路径失败")
	}

	// 检查文件是否存在并获取文件信息
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fiber.NewError(fiber.StatusInternalServerError, "文件不存在")
	}

	// 检查是否为文本文件
	if isTextFile(path) {
		// 设置文件大小限制（20MB）
		const maxFileSize = 20 * 1024 * 1024
		if fileInfo.Size() > maxFileSize {
			return fiber.NewError(fiber.StatusBadRequest, "文件太大，不支持在线预览")
		}

		// 分块读取文件
		file, err := os.Open(path)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "打开文件失败")
		}
		defer file.Close()

		// 使用缓冲读取
		const bufferSize = 1024 * 1024 // 1MB 缓冲区
		buffer := make([]byte, bufferSize)
		content := make([]byte, 0, fileInfo.Size())

		for {
			n, err := file.Read(buffer)
			if err != nil && err != io.EOF {
				return fiber.NewError(fiber.StatusInternalServerError, "读取文件失败")
			}
			if n == 0 {
				break
			}
			content = append(content, buffer[:n]...)
		}

		relPath, err := filepath.Rel(absRootDir, absPath)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "计算相对路径失败")
		}

		// 获取目录路径并规范化分隔符
		dirPath := normalizePath(filepath.Dir(relPath))
		if dirPath == "." {
			dirPath = ""
		}

		// 检测文件编码
		detector := chardet.NewTextDetector()
		result, err := detector.DetectBest(content)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "检测文件编码失败")
		}

		// 如果是非UTF-8编码，进行转换
		if result.Charset != "UTF-8" {
			encoding, err := ianaindex.IANA.Encoding(result.Charset)
			if err == nil && encoding != nil {
				decoder := encoding.NewDecoder()
				utf8Content, err := decoder.Bytes(content)
				if err == nil {
					content = utf8Content
				}
			}
		}

		// 记录操作日志
		adminlog.WriteLog(c, "view", "browse", 0, fmt.Sprintf("查看文件：%s", path))

		rootPath := "/admin/browse"

		return c.Render("admin/file", fiber.Map{
			"Title":    "文件内容",
			"Path":     filepath.Base(relPath),
			"DirPath":  dirPath,
			"Content":  string(content),
			"RootPath": rootPath,
		}, "admin/layout")
	}

	// 记录操作日志
	adminlog.WriteLog(c, "download", "browse", 0, fmt.Sprintf("下载文件：%s", path))

	// 非文本文件，直接下载
	return c.SendFile(path)
}

// handleBrowseDelete 处理文件删除请求
func handleBrowseDelete(c *fiber.Ctx, path string) error {
	err := os.Remove(path)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "文件删除失败",
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "browse", 0, fmt.Sprintf("删除文件：%s", path))

	return c.JSON(fiber.Map{
		"message": "文件删除成功",
	})
}

// uploadByFTP 处理 FTP 上传请求
func uploadByFTP(c *fiber.Ctx, rootPath string) error {
	filePath := c.Query("file")
	fileType := c.Query("type")

	if filePath == "" || (fileType != "apk" && fileType != "zip") {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的文件路径或类型",
		})
	}

	// 解码文件路径
	decodedPath, err := url.QueryUnescape(filePath)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "无效的路径编码",
		})
	}

	fullPath := filepath.Join(rootPath, decodedPath)

	// 上传到 FTP
	if err := uploadToFTP(fullPath, fileType); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "upload", "browse", 0, fmt.Sprintf("上传文件：%s", fullPath))

	return c.JSON(fiber.Map{
		"success": true,
		"message": "文件上传成功",
	})
}

// 上传文件到 FTP
func uploadToFTP(localPath string, fileType string) error {
	// 连接 FTP
	conn, err := ftp.Dial(fmt.Sprintf("%s:%s", app.Config.FTP.Host, app.Config.FTP.Port))
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("FTP连接失败: %v", err))
		return fmt.Errorf("FTP连接失败: %v", err)
	}
	defer conn.Quit()

	username, password, err := ReadFromBinaryFile(app.Config.Server.UserDataPath)
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("读取用户数据失败: %v", err))
		return fmt.Errorf("读取用户数据失败: %v", err)
	}

	// 登录
	if err := conn.Login(username, password); err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("FTP登录失败: %v", err))
		return fmt.Errorf("FTP登录失败: %v", err)
	}

	// 根据文件类型选择上传路径
	remotePath := app.Config.FTP.APKPath
	if fileType == "zip" {
		remotePath = app.Config.FTP.ZIPPath
	}

	// 确保远程路径使用正斜杠
	remotePath = strings.ReplaceAll(remotePath, "\\", "/")
	if !strings.HasPrefix(remotePath, "/") {
		remotePath = "/" + remotePath
	}
	if !strings.HasSuffix(remotePath, "/") {
		remotePath = remotePath + "/"
	}

	// 尝试切换到目标目录
	if err := conn.ChangeDir(remotePath); err != nil {
		// 如果目录不存在，尝试创建
		if err := createRemoteDirectories(conn, remotePath); err != nil {
			writeUploadLog(localPath, fileType, false, fmt.Sprintf("创建远程目录失败: %v", err))
			return fmt.Errorf("创建远程目录失败: %v", err)
		}
	}

	// 打开本地文件
	file, err := os.Open(localPath)
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("打开文件失败: %v", err))
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 获取文件名并确保它是有效的
	fileName := filepath.Base(localPath)
	fileName = strings.TrimSpace(fileName)
	if fileName == "" || fileName == "." || fileName == ".." {
		writeUploadLog(localPath, fileType, false, "无效的文件名")
		return fmt.Errorf("无效的文件名: %s", fileName)
	}

	// 上传文件
	err = conn.Stor(fileName, file)
	if err != nil {
		writeUploadLog(localPath, fileType, false, fmt.Sprintf("上传文件失败: %v", err))
		return fmt.Errorf("上传文件失败: %v", err)
	}

	// 记录成功日志
	writeUploadLog(localPath, fileType, true, "上传成功")
	return nil
}

// createRemoteDirectories 递归创建远程目录
func createRemoteDirectories(conn *ftp.ServerConn, path string) error {
	path = strings.Trim(path, "/")
	dirs := strings.Split(path, "/")
	currentPath := "/"

	for _, dir := range dirs {
		if dir == "" {
			continue
		}
		currentPath = filepath.Join(currentPath, dir)
		currentPath = strings.ReplaceAll(currentPath, "\\", "/")

		err := conn.ChangeDir(currentPath)
		if err != nil {
			// 如果目录不存在，创建它
			if err := conn.MakeDir(currentPath); err != nil {
				return fmt.Errorf("无法创建目录 %s: %v", currentPath, err)
			}
			// 创建后切换到该目录
			if err := conn.ChangeDir(currentPath); err != nil {
				return fmt.Errorf("无法切换到目录 %s: %v", currentPath, err)
			}
		}
	}
	return nil
}

// 添加日志写入函数
func writeUploadLog(localPath, fileType string, success bool, message string) error {
	// 确保日志目录存在
	if err := os.MkdirAll(app.Config.FTP.LogDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 准备日志内容
	logContent := fmt.Sprintf("[%s] File: %s, Type: %s, Success: %v, Message: %s\n",
		time.Now().Format("2006-01-02 15:04:05"),
		localPath,
		fileType,
		success,
		message,
	)

	// 获取当前日志文件路径
	logFile := filepath.Join(app.Config.FTP.LogDir, fmt.Sprintf("ftpupload_%s.log", time.Now().Format("20060102150405")))

	// 追加写入日志
	f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %v", err)
	}
	defer f.Close()

	if _, err := f.WriteString(logContent); err != nil {
		return fmt.Errorf("写入日志失败: %v", err)
	}

	return nil
}
