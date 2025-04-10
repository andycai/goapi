package browse

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func listFilesHandler(c *fiber.Ctx) error {
	path := c.Params("*")
	if path == "" {
		path = "."
	}

	// URL 解码路径
	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		return c.Status(400).SendString("Invalid path encoding")
	}

	// 获取配置的根目录的绝对路径
	rootDir := app.Config.Server.Output
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return c.Status(500).SendString("Invalid root directory configuration")
	}

	// 构建完整路径
	fullPath := filepath.Join(rootDir, decodedPath)

	// 获取完整路径的绝对路径
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return c.Status(400).SendString("Invalid path")
	}

	// 确保访问路径在根目录内
	if !strings.HasPrefix(absPath, absRootDir) {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: Path outside root directory")
	}

	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, fmt.Sprintf("File not found: %s", decodedPath))
	}

	// 如果是目录，显示目录内容
	if fileInfo.IsDir() {
		return handleBrowseDirectory(c, absPath)
	}

	// 如果是文件，显示文件内容
	return handleBrowseFile(c, absPath)
}

func deleteFileHandler(c *fiber.Ctx) error {
	path := c.Params("*")
	if path == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Path is required")
	}

	// URL 解码路径
	decodedPath, err := url.QueryUnescape(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid path encoding")
	}

	// 获取配置的根目录的绝对路径
	rootDir := app.Config.Server.Output
	absRootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Invalid root directory configuration")
	}

	// 构建完整路径
	fullPath := filepath.Join(rootDir, decodedPath)

	// 获取完整路径的绝对路径
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return c.Status(400).SendString("Invalid path")
	}

	// 确保删除路径在根目录内
	if !strings.HasPrefix(absPath, absRootDir) {
		return fiber.NewError(fiber.StatusForbidden, "Access denied: Path outside root directory")
	}

	// 检查是否是目录
	fileInfo, err := os.Stat(absPath)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("File not found")
	}
	if fileInfo.IsDir() {
		return c.Status(fiber.StatusBadRequest).SendString("Cannot delete directories")
	}

	return handleBrowseDelete(c, absPath)
}

func uploadFileHandler(c *fiber.Ctx) error {
	return uploadByFTP(c, app.Config.Server.Output)
}
