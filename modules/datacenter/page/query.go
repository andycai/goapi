package page

import (
	"errors"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

var (
	ErrPageNotFound          = errors.New("页面不存在")
	ErrPageSlugAlreadyExists = errors.New("页面别名已存在")
)

// 获取页面列表
func QueryPageList(page, limit int, status string) ([]models.Page, int64, error) {
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
func QueryPageByID(id int64) (models.Page, error) {
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
func QueryPageBySlug(slug string) (models.Page, error) {
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
func QuerySearchPages(query string, page, limit int) ([]models.Page, int64, error) {
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
