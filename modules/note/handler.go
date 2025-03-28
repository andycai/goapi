package note

import (
	"fmt"
	"time"

	"github.com/andycai/unitool/models"
	"github.com/andycai/unitool/modules/adminlog"
	"github.com/gofiber/fiber/v2"
)

// getNoteList 处理笔记列表页面
func getNoteList(c *fiber.Ctx) error {
	return c.Render("admin/notes", fiber.Map{
		"Title": "笔记管理",
		"Scripts": []string{
			"/static/js/admin/notes.js",
		},
	}, "admin/layout")
}

// getCategoryList 处理分类列表
func getCategoryList(c *fiber.Ctx) error {
	categories, err := srv.ListCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取分类列表失败"})
	}
	return c.JSON(categories)
}

// getNoteTree 处理笔记树结构
func getNoteTree(c *fiber.Ctx) error {
	notes, err := srv.ListNotes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记列表失败"})
	}
	return c.JSON(notes)
}

// getNoteDetail 处理笔记详情
func getNoteDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	note, err := srv.GetNoteByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记失败"})
	}
	if note == nil {
		return c.Status(404).JSON(fiber.Map{"error": "笔记不存在"})
	}

	return c.JSON(note)
}

type NoteRequest struct {
	Title      string `json:"title"`
	Content    string `json:"content"`
	CategoryID uint   `json:"category_id"`
	ParentID   uint   `json:"parent_id"`
	IsPublic   uint8  `json:"is_public"`
	RoleIDs    []uint `json:"role_ids"`
}

// createNote 处理创建笔记
func createNote(c *fiber.Ctx) error {
	var req NoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	user := c.Locals("user").(*models.User)

	note := &models.Note{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		ParentID:   req.ParentID,
		IsPublic:   req.IsPublic,
		CreatedBy:  user.ID,
		UpdatedBy:  user.ID,
	}

	if err := srv.CreateNote(note, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建笔记失败"})
	}

	adminlog.CreateAdminLog(c, "create", "note", note.ID, fmt.Sprintf("创建笔记：%s", note.Title))

	return c.JSON(note)
}

// updateNote 处理更新笔记
func updateNote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	var req NoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	note, err := srv.GetNoteByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记失败"})
	}
	if note == nil {
		return c.Status(404).JSON(fiber.Map{"error": "笔记不存在"})
	}

	user := c.Locals("user").(*models.User)

	note.Title = req.Title
	note.Content = req.Content
	note.CategoryID = req.CategoryID
	note.ParentID = req.ParentID
	note.IsPublic = req.IsPublic
	note.UpdatedBy = user.ID
	note.UpdatedAt = time.Now()

	if err := srv.UpdateNote(note, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新笔记失败"})
	}

	adminlog.CreateAdminLog(c, "update", "note", note.ID, fmt.Sprintf("更新笔记：%s", note.Title))

	return c.JSON(note)
}

// deleteNote 处理删除笔记
func deleteNote(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	note, err := srv.GetNoteByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记失败"})
	}
	if note == nil {
		return c.Status(404).JSON(fiber.Map{"error": "笔记不存在"})
	}

	if err := srv.DeleteNote(note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除笔记失败"})
	}

	adminlog.CreateAdminLog(c, "delete", "note", note.ID, fmt.Sprintf("删除笔记：%s", note.Title))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

type CategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    uint   `json:"parent_id"`
	IsPublic    uint8  `json:"is_public"`
	RoleIDs     []uint `json:"role_ids"`
}

// createCategory 处理创建分类
func createCategory(c *fiber.Ctx) error {
	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	user := c.Locals("user").(*models.User)

	category := &models.NoteCategory{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsPublic:    req.IsPublic,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}

	if err := srv.CreateCategory(category, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建分类失败"})
	}

	adminlog.CreateAdminLog(c, "create", "note_category", category.ID, fmt.Sprintf("创建笔记分类：%s", category.Name))

	return c.JSON(category)
}

// updateCategory 处理更新分类
func updateCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	category, err := srv.GetCategoryByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取分类失败"})
	}
	if category == nil {
		return c.Status(404).JSON(fiber.Map{"error": "分类不存在"})
	}

	user := c.Locals("user").(*models.User)

	category.Name = req.Name
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.IsPublic = req.IsPublic
	category.UpdatedBy = user.ID
	category.UpdatedAt = time.Now()

	if err := srv.UpdateCategory(category, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新分类失败"})
	}

	adminlog.CreateAdminLog(c, "update", "note_category", category.ID, fmt.Sprintf("更新笔记分类：%s", category.Name))

	return c.JSON(category)
}

// deleteCategory 处理删除分类
func deleteCategory(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	// 检查是否有子分类
	hasChildren, err := srv.HasCategoryChildren(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "检查子分类失败"})
	}
	if hasChildren {
		return c.Status(400).JSON(fiber.Map{"error": "该分类下有子分类，无法删除"})
	}

	// 检查是否有笔记
	hasNotes, err := srv.HasCategoryNotes(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "检查分类笔记失败"})
	}
	if hasNotes {
		return c.Status(400).JSON(fiber.Map{"error": "该分类下有笔记，无法删除"})
	}

	category, err := srv.GetCategoryByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取分类失败"})
	}
	if category == nil {
		return c.Status(404).JSON(fiber.Map{"error": "分类不存在"})
	}

	if err := srv.DeleteCategory(category); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除分类失败"})
	}

	adminlog.CreateAdminLog(c, "delete", "note_category", category.ID, fmt.Sprintf("删除笔记分类：%s", category.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// getPublicNotes 处理公开笔记列表
func getPublicNotes(c *fiber.Ctx) error {
	notes, err := srv.ListPublicNotes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取公开笔记列表失败"})
	}
	return c.JSON(notes)
}

// getPublicNoteDetail 处理公开笔记详情
func getPublicNoteDetail(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	note, err := srv.GetNoteByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记失败"})
	}
	if note == nil {
		return c.Status(404).JSON(fiber.Map{"error": "笔记不存在"})
	}

	if note.IsPublic != 1 {
		return c.Status(403).JSON(fiber.Map{"error": "该笔记不是公开的"})
	}

	// 增加浏览次数
	if err := srv.IncrementNoteViewCount(note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新浏览次数失败"})
	}

	return c.JSON(note)
}

// getPublicCategories 处理公开分类列表
func getPublicCategories(c *fiber.Ctx) error {
	categories, err := srv.ListPublicCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取公开分类列表失败"})
	}
	return c.JSON(categories)
}
