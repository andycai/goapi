package user

import (
	"fmt"
	"strconv"
	"time"

	"github.com/andycai/goapi/internal"
	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/modules/system/adminlog"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// listUsersHandler 获取用户列表
func listUsersHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	users, total, err := QueryUserList(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(internal.Response{
			Code:    1,
			Message: "获取用户列表失败: " + err.Error(),
		})
	}

	// var users []models.User
	// if err := app.DB.Preload("Role").Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
	// 	return c.Status(500).JSON(fiber.Map{"error": "获取用户列表失败"})
	// }

	var response = internal.Response{
		Code: 0,
		Data: fiber.Map{
			"users": users,
			"total": total,
		},
	}

	return c.JSON(response)
}

// createUserHandler 创建用户
func createUserHandler(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 检查用户名是否已存在
	var count int64
	app.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "用户名已存在"})
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "密码加密失败"})
	}

	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		Nickname:  req.Nickname,
		RoleID:    req.RoleID,
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := app.DB.Create(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建用户失败"})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "create", "user", user.ID, fmt.Sprintf("创建用户：%s", user.Username))

	return c.JSON(user)
}

// updateUserHandler 更新用户
func updateUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.RoleID != 0 {
		updates["role_id"] = req.RoleID
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "密码加密失败"})
		}
		updates["password"] = string(hashedPassword)
	}

	var user models.User
	if err := app.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "用户不存在"})
	}

	if err := app.DB.Model(&user).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新用户失败"})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "update", "user", user.ID, fmt.Sprintf("更新用户：%s", user.Username))

	return c.JSON(user)
}

// deleteUserHandler 删除用户
func deleteUserHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if err := app.DB.First(&user, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "用户不存在",
		})
	}

	if err := app.DB.Delete(&user).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "删除用户失败",
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "user", user.ID, fmt.Sprintf("删除用户：%s", user.Username))

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}
