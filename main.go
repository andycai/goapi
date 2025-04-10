package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/andycai/goapi/core"
	_ "github.com/andycai/goapi/docs" // 导入 swagger docs
	"github.com/andycai/goapi/lib/database"
	_ "github.com/andycai/goapi/modules"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/swagger"
	"github.com/gofiber/template/html/v2"
	"gorm.io/gorm"
)

var (
	Version   string
	BuildTime string
)

// @title UnityTool API
// @version 1.0
// @description Unity游戏开发辅助工具API文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url https://github.com/andycai/unitool
// @contact.email huayicai@gamil.com

// @license.name MIT
// @license.url https://github.com/andycai/goapi/blob/main/LICENSE

// @host localhost:3000
// @BasePath /api/v1
func main() {
	// 加载配置文件
	if err := core.LoadConfig(); err != nil {
		log.Fatalf("无法加载配置文件: %v", err)
	}

	// 初始化数据库
	dbConfig := core.GetDatabaseConfig()
	db, err := database.InitRDBMS(
		os.Stdout,
		dbConfig.Driver,
		dbConfig.DSN,
		dbConfig.MaxOpenConns,
		dbConfig.MaxIdleConns,
		dbConfig.ConnMaxLifetime,
	)
	if err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// 初始化模板引擎
	engine := html.New("./templates", ".html")

	if core.IsDevelopment() {
		engine.Reload(true) // 开发模式下启用模板重载
		engine.Debug(true)  // 开发模式下启用调试信息
	}

	// 添加模板函数
	engine.AddFunc("yield", func() string { return "" })
	engine.AddFunc("partial", func(name string, data interface{}) template.HTML {
		return template.HTML("")
	})
	// 添加 hasSuffix 函数用于检查文件扩展名
	engine.AddFunc("hasSuffix", strings.HasSuffix)
	// 添加 splitPath 函数用于分割路径，同时处理 Windows 和 Unix 风格的路径分隔符
	engine.AddFunc("splitPath", func(path string) []string {
		if path == "" {
			return []string{}
		}
		// 先将所有反斜杠转换为正斜杠
		path = strings.ReplaceAll(path, "\\", "/")
		// 分割路径
		return strings.Split(path, "/")
	})
	// 添加 sub 函数用于数字减法
	engine.AddFunc("sub", func(a, b int) int {
		return a - b
	})

	// 创建 Fiber 应用，并配置模板引擎
	fiberApp := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			// 根据错误码返回对应的错误页面
			switch code {
			case fiber.StatusBadRequest:
				return c.Status(code).Render("errors/400", fiber.Map{}, "")
			case fiber.StatusUnauthorized:
				return c.Status(code).Render("errors/401", fiber.Map{}, "")
			case fiber.StatusForbidden:
				return c.Status(code).Render("errors/403", fiber.Map{}, "")
			case fiber.StatusNotFound:
				return c.Status(code).Render("errors/404", fiber.Map{}, "")
			default:
				return c.Status(code).Render("errors/500", fiber.Map{}, "")
			}
		},
	})

	// 配置 CORS 中间件
	corsConfig := core.GetCorsConfig()
	if corsConfig.Enabled {
		fiberApp.Use(cors.New(cors.Config{
			AllowOrigins:     strings.Join(corsConfig.AllowedOrigins, ","),
			AllowMethods:     strings.Join(corsConfig.AllowedMethods, ","),
			AllowHeaders:     strings.Join(corsConfig.AllowedHeaders, ","),
			AllowCredentials: corsConfig.AllowCredentials,
			MaxAge:           int(time.Duration(corsConfig.MaxAge) * time.Hour),
		}))
	}

	// 添加 Swagger 路由
	fiberApp.Get("/swagger/*", swagger.New(swagger.Config{
		Title:        "UnityTool API",
		DeepLinking:  true,
		DocExpansion: "list",
	}))

	app := core.NewApp()
	app.Start([]*gorm.DB{db}, fiberApp)

	// 启动服务器
	fiberApp.Listen(fmt.Sprintf("%s:%d", app.Config.Server.Host, app.Config.Server.Port))
}
