package post

import (
	"errors"
	"time"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 获取文章列表
func getPostList(page, limit int, status string) ([]models.Post, int64, error) {
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
func getPostByID(id int64) (models.Post, error) {
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
func getPostBySlug(slug string) (models.Post, error) {
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
func searchPosts(query string, page, limit int) ([]models.Post, int64, error) {
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

// 添加文章
func addPost(post *models.Post) error {
	// 检查别名是否已存在
	var count int64
	if err := app.DB.Model(&models.Post{}).Where("slug = ?", post.Slug).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrPostSlugAlreadyExists
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	return app.DB.Create(post).Error
}

// 更新文章
func updatePost(post *models.Post) error {
	// 检查文章是否存在
	var count int64
	if err := app.DB.Model(&models.Post{}).Where("id = ?", post.ID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPostNotFound
	}

	// 检查新别名是否与其他文章冲突
	if err := app.DB.Model(&models.Post{}).Where("slug = ? AND id != ?", post.Slug, post.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrPostSlugAlreadyExists
	}

	post.UpdatedAt = time.Now()

	return app.DB.Save(post).Error
}

// 删除文章
func deletePost(id int64) error {
	// 检查文章是否存在
	var count int64
	if err := app.DB.Model(&models.Post{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPostNotFound
	}

	return app.DB.Delete(&models.Post{}, id).Error
}
