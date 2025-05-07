package permission

import (
	"fmt"
	"time"

	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/modules/system/adminlog"
	"github.com/gofiber/fiber/v2"
)

type CreatePermissionRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type UpdatePermissionRequest struct {
	Name        string `json:"name,omitempty"`
	Code        string `json:"code,omitempty"`
	Description string `json:"description,omitempty"`
}

// listPermissionsHandler 获取权限列表
func listPermissionsHandler(c *fiber.Ctx) error {
	// 获取分页参数
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 计算偏移量
	offset := (page - 1) * pageSize

	// 获取总记录数
	var total int64
	if err := app.DB.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取权限总数失败"})
	}

	// 获取分页数据
	var permissions []models.Permission
	if err := app.DB.Offset(offset).Limit(pageSize).Find(&permissions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取权限列表失败"})
	}

	return c.JSON(fiber.Map{
		"items": permissions,
		"total": total,
	})
}

func listAllPermissionsHandler(c *fiber.Ctx) error {
	// 获取总记录数
	var total int64
	if err := app.DB.Model(&models.Permission{}).Count(&total).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取权限总数失败"})
	}

	var permissions []models.Permission
	if err := app.DB.Find(&permissions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取权限列表失败"})
	}

	return c.JSON(fiber.Map{
		"items": permissions,
		"total": total,
	})
}

// createPermissionHandler 创建权限
func createPermissionHandler(c *fiber.Ctx) error {
	var req CreatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	// 检查权限编码是否已存在
	var count int64
	app.DB.Model(&models.Permission{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "权限编码已存在"})
	}

	permission := models.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := app.DB.Create(&permission).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建权限失败"})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "create", "permission", permission.ID, fmt.Sprintf("创建权限：%s", permission.Name))

	return c.JSON(permission)
}

// updatePermissionHandler 更新权限
func updatePermissionHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdatePermissionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	var permission models.Permission
	if err := app.DB.First(&permission, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "权限不存在"})
	}

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Code != "" {
		// 检查新的权限编码是否已存在
		var count int64
		app.DB.Model(&models.Permission{}).Where("code = ? AND id != ?", req.Code, id).Count(&count)
		if count > 0 {
			return c.Status(400).JSON(fiber.Map{"error": "权限编码已存在"})
		}
		updates["code"] = req.Code
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}

	if err := app.DB.Model(&permission).Updates(updates).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新权限失败"})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "update", "permission", permission.ID, fmt.Sprintf("更新权限：%s", permission.Name))

	return c.JSON(permission)
}

// deletePermissionHandler 删除权限
func deletePermissionHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	// 检查权限是否被角色使用
	var count int64
	if err := app.DB.Table("role_permissions").Where("permission_id = ?", id).Count(&count).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "检查权限使用状态失败"})
	}

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "该权限正在被角色使用，无法删除"})
	}

	var permission models.Permission
	if err := app.DB.First(&permission, id).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "权限不存在"})
	}

	if err := app.DB.Delete(&permission).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除权限失败"})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "permission", permission.ID, fmt.Sprintf("删除权限：%s", permission.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}
