package page

import (
	"errors"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 获取页面列表
func getPageList(page, limit int, status string) ([]models.Page, int64, error) {
	var pages []models.Page
	var total int64

	db := app.DB.Model(&models.Page{})
	if status != "" {
		db = db.Where("status = ?", status)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Order("id desc").Find(&pages).Error; err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// 根据ID获取页面
func getPageByID(id int64) (models.Page, error) {
	var page models.Page
	if err := app.DB.First(&page, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return page, ErrPageNotFound
		}
		return page, err
	}
	return page, nil
}

// 根据别名获取页面
func getPageBySlug(slug string) (models.Page, error) {
	var page models.Page
	if err := app.DB.Where("slug = ?", slug).First(&page).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return page, ErrPageNotFound
		}
		return page, err
	}
	return page, nil
}

// 搜索页面
func searchPages(query string, page, limit int) ([]models.Page, int64, error) {
	var pages []models.Page
	var total int64

	searchQuery := app.DB.Model(&models.Page{}).
		Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%")

	err := searchQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		searchQuery = searchQuery.Offset(offset).Limit(limit)
	}

	if err := searchQuery.Order("id desc").Find(&pages).Error; err != nil {
		return nil, 0, err
	}

	return pages, total, nil
}

// 添加页面
func addPage(page *models.Page) error {
	// 检查别名是否已存在
	var count int64
	if err := app.DB.Model(&models.Page{}).Where("slug = ?", page.Slug).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrPageSlugAlreadyExists
	}

	page.CreatedAt = time.Now()
	page.UpdatedAt = time.Now()

	return app.DB.Create(page).Error
}

// 更新页面
func updatePage(page *models.Page) error {
	// 检查页面是否存在
	var count int64
	if err := app.DB.Model(&models.Page{}).Where("id = ?", page.ID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPageNotFound
	}

	// 检查新别名是否与其他页面冲突
	if err := app.DB.Model(&models.Page{}).Where("slug = ? AND id != ?", page.Slug, page.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrPageSlugAlreadyExists
	}

	page.UpdatedAt = time.Now()

	return app.DB.Save(page).Error
}

// 删除页面
func deletePage(id int64) error {
	// 检查页面是否存在
	var count int64
	if err := app.DB.Model(&models.Page{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPageNotFound
	}

	return app.DB.Delete(&models.Page{}, id).Error
}
