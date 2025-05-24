package post

import (
	"strconv"

	"github.com/andycai/goapi/models"
	"github.com/gofiber/fiber/v2"
)

// 获取文章列表
func listPostHandler(c *fiber.Ctx) error {
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

	posts, total, err := QueryPostList(page, limit, status)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取文章列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"posts": posts,
		"total": total,
	})
}

// 获取文章详情
func getPostBySlugHandler(c *fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "文章别名不能为空",
		})
	}

	post, err := QueryPostBySlug(slug)
	if err != nil {
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取文章详情失败: " + err.Error(),
		})
	}

	return c.JSON(post)
}

// 添加文章
func addPostHandler(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if post.Title == "" || post.Content == "" || post.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "标题、内容和别名不能为空",
		})
	}

	if err := CommandAddPost(post); err != nil {
		if err == ErrPostSlugAlreadyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加文章失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加文章成功",
		"id":      post.ID,
	})
}

// 编辑文章
func editPostHandler(c *fiber.Ctx) error {
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if post.ID == 0 || post.Title == "" || post.Content == "" || post.Slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID、标题、内容和别名不能为空",
		})
	}

	if err := CommandUpdatePost(post); err != nil {
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == ErrPostSlugAlreadyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新文章失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新文章成功",
	})
}

// 删除文章
func deletePostHandler(c *fiber.Ctx) error {
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

	if err := CommandDeletePost(req.ID); err != nil {
		if err == ErrPostNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除文章失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除文章成功",
	})
}

// 搜索文章
func searchPostHandler(c *fiber.Ctx) error {
	query := c.Query("q", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	posts, total, err := QuerySearchPosts(query, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "搜索文章失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"posts": posts,
		"total": total,
	})
}
