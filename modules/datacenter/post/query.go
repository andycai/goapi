package post

import (
	"errors"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound          = errors.New("文章不存在")
	ErrPostSlugAlreadyExists = errors.New("文章别名已存在")
)

// 获取文章列表
func QueryPostList(page, limit int, status string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	db := app.DB.Model(&models.Post{})
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

	if err := db.Order("id desc").Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

// 根据ID获取文章
func QueryPostByID(id int64) (models.Post, error) {
	var post models.Post
	if err := app.DB.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, ErrPostNotFound
		}
		return post, err
	}
	return post, nil
}

// 根据别名获取文章
func QueryPostBySlug(slug string) (models.Post, error) {
	var post models.Post
	if err := app.DB.Where("slug = ?", slug).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, ErrPostNotFound
		}
		return post, err
	}
	return post, nil
}

// 搜索文章
func QuerySearchPosts(query string, page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	searchQuery := app.DB.Model(&models.Post{}).
		Where("title LIKE ? OR content LIKE ?", "%"+query+"%", "%"+query+"%")

	err := searchQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		searchQuery = searchQuery.Offset(offset).Limit(limit)
	}

	if err := searchQuery.Order("id desc").Find(&posts).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}
