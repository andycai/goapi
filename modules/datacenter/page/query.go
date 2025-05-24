package page

import (
	"github.com/andycai/goapi/models"
)

// 获取页面列表
func QueryPageList(page, limit int, status string) ([]models.Page, int64, error) {
	return getPageList(page, limit, status)
}

// 根据ID获取页面
func QueryPageByID(id int64) (models.Page, error) {
	return getPageByID(id)
}

// 根据别名获取页面
func QueryPageBySlug(slug string) (models.Page, error) {
	return getPageBySlug(slug)
}

// 搜索页面
func QuerySearchPages(query string, page, limit int) ([]models.Page, int64, error) {
	return searchPages(query, page, limit)
}
