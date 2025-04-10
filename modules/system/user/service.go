package user

import (
	"github.com/andycai/goapi/core"
	"github.com/andycai/goapi/models"
	"github.com/gofiber/fiber/v2"
)

// GetByID 获取用户
func GetByID(id uint) *models.User {
	var vo models.User
	app.DB.Model(&vo).
		Where("id", id).
		First(&vo)

	return &vo
}

// Current 获取当前用户
func Current(c *fiber.Ctx) *models.User {
	isAuthenticated, userID := core.GetSession(c)

	if !isAuthenticated {
		return nil
	}

	return GetByID(userID)
}
