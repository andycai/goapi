package menu

import (
	"fmt"

	"github.com/andycai/goapi/modules/system/adminlog"
	"github.com/gofiber/fiber/v2"
)

// listMenusHandler 获取菜单列表
func listMenusHandler(c *fiber.Ctx) error {
	menus, err := dao.GetMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取菜单列表失败",
		})
	}
	return c.JSON(menus)
}

// getMenuTreeHandler 获取菜单树
func getMenuTreeHandler(c *fiber.Ctx) error {
	menus, err := dao.GetMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取菜单列表失败",
		})
	}
	tree := dao.BuildMenuTree(menus, 0)
	return c.JSON(tree)
}

// createMenuHandler 创建菜单
func createMenuHandler(c *fiber.Ctx) error {
	menu := new(Menu)
	if err := c.BodyParser(menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := dao.CreateMenu(menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建菜单失败",
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "create", "menu", menu.ID, fmt.Sprintf("创建菜单：%s", menu.Name))

	return c.JSON(menu)
}

// updateMenuHandler 更新菜单
func updateMenuHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的菜单ID",
		})
	}

	menu := new(Menu)
	if err := c.BodyParser(menu); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	menu.ID = uint(id)
	if err := dao.UpdateMenu(menu); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新菜单失败",
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "update", "menu", menu.ID, fmt.Sprintf("更新菜单：%s", menu.Name))

	return c.JSON(menu)
}

// deleteMenuHandler 删除菜单
func deleteMenuHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的菜单ID",
		})
	}

	// 获取菜单信息用于日志记录
	menu, err := dao.GetMenuByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "菜单不存在",
		})
	}

	if err := dao.DeleteMenu(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除菜单失败",
		})
	}

	// 记录操作日志
	adminlog.WriteLog(c, "delete", "menu", uint(id), fmt.Sprintf("删除菜单：%s", menu.Name))

	return c.SendStatus(fiber.StatusNoContent)
}
