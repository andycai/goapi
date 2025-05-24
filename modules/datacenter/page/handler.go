package page

import (
	"strconv"

	"github.com/andycai/goapi/models"
	"github.com/gofiber/fiber/v2"
)

// 获取页面列表
func listPageHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	status := c.Query("status", "")

	pages, total, err := QueryPageList(page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取页面列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"pages": pages,
		"total": total,
	})
}

// 获取页面详情
func getPageBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "页面别名不能为空",
		})
	}

	page, err := QueryPageBySlug(slug)
	if err != nil {
		if err == ErrPageNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取页面详情失败: " + err.Error(),
		})
	}

	return c.JSON(page)
}

// 添加页面
func addPageHandler(c *fiber.Ctx) error {
	page := new(models.Page)
	if err := c.BodyParser(page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if page.Title == "" || page.Content == "" || page.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "标题、内容和别名不能为空",
		})
	}

	if err := CommandAddPage(page); err != nil {
		if err == ErrPageSlugAlreadyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加页面失败: " + err.Error(),
		})
	}

	// 获取新创建的页面详情
	pageResp, err := QueryPageByID(page.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取页面详情失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加页面成功",
		"page":    pageResp,
	})
}

// 编辑页面
func editPageHandler(c *fiber.Ctx) error {
	page := new(models.Page)
	if err := c.BodyParser(page); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if page.ID == 0 || page.Title == "" || page.Content == "" || page.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID、标题、内容和别名不能为空",
		})
	}

	if err := CommandUpdatePage(page); err != nil {
		if err == ErrPageNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == ErrPageSlugAlreadyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新页面失败: " + err.Error(),
		})
	}

	// 获取更新后的页面详情
	pageResp, err := QueryPageByID(page.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取页面详情失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新页面成功",
		"page":    pageResp,
	})
}

// 删除页面
func deletePageHandler(c *fiber.Ctx) error {
	type DeleteRequest struct {
		ID int64 `json:"id"`
	}

	var req DeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if req.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的ID",
		})
	}

	if err := CommandDeletePage(req.ID); err != nil {
		if err == ErrPageNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除页面失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除页面成功",
	})
}

// 搜索页面
func searchPageHandler(c *fiber.Ctx) error {
	query := c.Query("q", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	pages, total, err := QuerySearchPages(query, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "搜索页面失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"pages": pages,
		"total": total,
	})
}
