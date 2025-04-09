package note

import (
	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// getNoteByID 获取笔记
func getNoteByID(id uint) (*models.Note, error) {
	var note models.Note
	if err := app.DB.Preload("Category").Preload("Roles").First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

// getCategoryByID 获取分类
func getCategoryByID(id uint) (*models.NoteCategory, error) {
	var category models.NoteCategory
	if err := app.DB.Preload("Roles").First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// getNotes 获取笔记列表
func getNotes() ([]models.Note, error) {
	var notes []models.Note
	if err := app.DB.Preload("Category").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// getCategories 获取分类列表
func getCategories() ([]models.NoteCategory, error) {
	var categories []models.NoteCategory
	if err := app.DB.Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// getPublicNotes 获取公开笔记列表
func getPublicNotes() ([]models.Note, error) {
	var notes []models.Note
	if err := app.DB.Where("is_public = ?", 1).Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

// getPublicCategories 获取公开分类列表
func getPublicCategories() ([]models.NoteCategory, error) {
	var categories []models.NoteCategory
	if err := app.DB.Where("is_public = ?", 1).Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// hasCategoryChildren 检查分类是否有子分类
func hasCategoryChildren(id uint) (bool, error) {
	var count int64
	if err := app.DB.Model(&models.NoteCategory{}).Where("parent_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// hasCategoryNotes 检查分类是否有笔记
func hasCategoryNotes(id uint) (bool, error) {
	var count int64
	if err := app.DB.Model(&models.Note{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// incrementNoteViewCount 增加笔记浏览次数
func incrementNoteViewCount(note *models.Note) error {
	return app.DB.Model(note).UpdateColumn("view_count", note.ViewCount+1).Error
}

// createNote 创建笔记
func createNote(note *models.Note, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(note).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(note).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// updateNote 更新笔记
func updateNote(note *models.Note, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(note).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(note).Association("Roles").Replace(roles); err != nil {
				return err
			}
		} else {
			if err := tx.Model(note).Association("Roles").Clear(); err != nil {
				return err
			}
		}

		return nil
	})
}

// deleteNote 删除笔记
func deleteNote(note *models.Note) error {
	return app.DB.Delete(note).Error
}

// createCategory 创建分类
func createCategory(category *models.NoteCategory, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(category).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(category).Association("Roles").Replace(roles); err != nil {
				return err
			}
		}

		return nil
	})
}

// updateCategory 更新分类
func updateCategory(category *models.NoteCategory, roleIDs []uint) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(category).Error; err != nil {
			return err
		}

		if len(roleIDs) > 0 {
			var roles []models.Role
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
			if err := tx.Model(category).Association("Roles").Replace(roles); err != nil {
				return err
			}
		} else {
			if err := tx.Model(category).Association("Roles").Clear(); err != nil {
				return err
			}
		}

		return nil
	})
}

// deleteCategory 删除分类
func deleteCategory(category *models.NoteCategory) error {
	return app.DB.Delete(category).Error
}
