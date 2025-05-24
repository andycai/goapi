package post

import (
	"errors"

	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/modules/datacenter/db"
	"gorm.io/gorm"
)

var (
	ErrPostNotFound          = errors.New("文章不存在")
	ErrPostSlugAlreadyExists = errors.New("文章别名已存在")
)

// 初始化服务
func initService() {
	// 暂无需要初始化的服务逻辑
}

// 获取文章列表
func getPostList(page, limit int, status string) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	database := db.GetDB()
	db := database.Model(&models.Post{})
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
func getPostByID(id int64) (models.Post, error) {
	var post models.Post
	database := db.GetDB()
	if err := database.First(&post, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, ErrPostNotFound
		}
		return post, err
	}
	return post, nil
}

// 根据别名获取文章
func getPostBySlug(slug string) (models.Post, error) {
	var post models.Post
	database := db.GetDB()
	if err := database.Where("slug = ?", slug).First(&post).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return post, ErrPostNotFound
		}
		return post, err
	}
	return post, nil
}

// 搜索文章
func searchPosts(query string, page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	database := db.GetDB()
	searchQuery := database.Model(&models.Post{}).
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
