package note

import (
	"fmt"
	"time"

	"github.com/andycai/unitool/models"
	"github.com/andycai/unitool/modules/system/adminlog"
	"github.com/gofiber/fiber/v2"
)

// listNotesHandler 处理笔记列表页面
func listNotesHandler(c *fiber.Ctx) error {
	return c.Render("admin/notes", fiber.Map{
		"Title": "笔记管理",
		"Scripts": []string{
			"/static/js/marked-15.0.7.js",
			"/static/js/admin/notes.js",
		},
	}, "admin/layout")
}

// listCategoriesHandler 处理分类列表
func listCategoriesHandler(c *fiber.Ctx) error {
	categories, err := getCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取分类列表失败"})
	}
	return c.JSON(categories)
}

// getNoteTreeHandler 处理笔记树结构
func getNoteTreeHandler(c *fiber.Ctx) error {
	notes, err := getNotes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记列表失败"})
	}
	return c.JSON(notes)
}

// getNoteDetailHandler 处理笔记详情
func getNoteDetailHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	note, err := getNoteByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记失败"})
	}
	if note == nil {
		return c.Status(404).JSON(fiber.Map{"error": "笔记不存在"})
	}

	return c.JSON(note)
}

// createNoteHandler 处理创建笔记
func createNoteHandler(c *fiber.Ctx) error {
	var req NoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	user := app.CurrentUser(c)

	note := &models.Note{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: req.CategoryID,
		ParentID:   req.ParentID,
		IsPublic:   req.IsPublic,
		CreatedBy:  user.ID,
		UpdatedBy:  user.ID,
	}

	if err := createNote(note, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建笔记失败"})
	}

	adminlog.WriteLog(c, "create", "note", note.ID, fmt.Sprintf("创建笔记：%s", note.Title))

	return c.JSON(note)
}

// updateNoteHandler 处理更新笔记
func updateNoteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	var req NoteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	note, err := getNoteByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记失败"})
	}
	if note == nil {
		return c.Status(404).JSON(fiber.Map{"error": "笔记不存在"})
	}

	user := app.CurrentUser(c)

	note.Title = req.Title
	note.Content = req.Content
	note.CategoryID = req.CategoryID
	note.ParentID = req.ParentID
	note.IsPublic = req.IsPublic
	note.UpdatedBy = user.ID
	note.UpdatedAt = time.Now()

	if err := updateNote(note, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新笔记失败"})
	}

	adminlog.WriteLog(c, "update", "note", note.ID, fmt.Sprintf("更新笔记：%s", note.Title))

	return c.JSON(note)
}

// deleteNoteHandler 处理删除笔记
func deleteNoteHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	note, err := getNoteByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取笔记失败"})
	}
	if note == nil {
		return c.Status(404).JSON(fiber.Map{"error": "笔记不存在"})
	}

	if err := deleteNote(note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除笔记失败"})
	}

	adminlog.WriteLog(c, "delete", "note", note.ID, fmt.Sprintf("删除笔记：%s", note.Title))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// createCategoryHandler 处理创建分类
func createCategoryHandler(c *fiber.Ctx) error {
	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	user := app.CurrentUser(c)

	category := &models.NoteCategory{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsPublic:    req.IsPublic,
		CreatedBy:   user.ID,
		UpdatedBy:   user.ID,
	}

	if err := createCategory(category, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "创建分类失败"})
	}

	adminlog.WriteLog(c, "create", "note_category", category.ID, fmt.Sprintf("创建笔记分类：%s", category.Name))

	return c.JSON(category)
}

// updateCategoryHandler 处理更新分类
func updateCategoryHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	var req CategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的请求数据"})
	}

	category, err := getCategoryByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取分类失败"})
	}
	if category == nil {
		return c.Status(404).JSON(fiber.Map{"error": "分类不存在"})
	}

	user := app.CurrentUser(c)

	category.Name = req.Name
	category.Description = req.Description
	category.ParentID = req.ParentID
	category.IsPublic = req.IsPublic
	category.UpdatedBy = user.ID
	category.UpdatedAt = time.Now()

	if err := updateCategory(category, req.RoleIDs); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新分类失败"})
	}

	adminlog.WriteLog(c, "update", "note_category", category.ID, fmt.Sprintf("更新笔记分类：%s", category.Name))

	return c.JSON(category)
}

// deleteCategoryHandler 处理删除分类
func deleteCategoryHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	// 检查是否有子分类
	hasChildren, err := hasCategoryChildren(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "检查子分类失败"})
	}
	if hasChildren {
		return c.Status(400).JSON(fiber.Map{"error": "该分类下有子分类，无法删除"})
	}

	// 检查是否有笔记
	hasNotes, err := hasCategoryNotes(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "检查分类笔记失败"})
	}
	if hasNotes {
		return c.Status(400).JSON(fiber.Map{"error": "该分类下有笔记，无法删除"})
	}

	category, err := getCategoryByID(uint(id))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取分类失败"})
	}
	if category == nil {
		return c.Status(404).JSON(fiber.Map{"error": "分类不存在"})
	}

	if err := deleteCategory(category); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "删除分类失败"})
	}

	adminlog.WriteLog(c, "delete", "note_category", category.ID, fmt.Sprintf("删除笔记分类：%s", category.Name))

	return c.JSON(fiber.Map{"message": "删除成功"})
}

// getPublicNotesHandler 处理公开笔记列表
func getPublicNotesHandler(c *fiber.Ctx) error {
	notes, err := getPublicNotes()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取公开笔记列表失败"})
	}
	return c.JSON(notes)
}

// getPublicNoteDetailHandler 处理公开笔记详情
func getPublicNoteDetailHandler(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "无效的ID"})
	}

	note, err := getNoteByID(uint(id))
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
	if err := incrementNoteViewCount(note); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "更新浏览次数失败"})
	}

	return c.JSON(note)
}

// getPublicCategoriesHandler 处理公开分类列表
func getPublicCategoriesHandler(c *fiber.Ctx) error {
	categories, err := getPublicCategories()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "获取公开分类列表失败"})
	}
	return c.JSON(categories)
}
