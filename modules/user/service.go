package user

import (
	"github.com/andycai/unitool/core"
	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

var srv *UserService

// UserService 用户服务
type UserService struct{}

// initService 初始化服务
func initService() {
	srv = &UserService{}
}

// GetByID 获取用户
func (s *UserService) GetByID(id uint) *models.User {
	var vo models.User
	app.DB.Model(&vo).
		Where("id", id).
		First(&vo)

	return &vo
}

// Current 获取当前用户
func (s *UserService) Current(c *fiber.Ctx) *models.User {
	isAuthenticated, userID := core.GetSession(c)

	if !isAuthenticated {
		return nil
	}

	return s.GetByID(userID)
}

// GetService 获取服务实例
func GetService() *UserService {
	return srv
}
