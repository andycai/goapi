package post

import (
	"time"

	"github.com/andycai/goapi/models"
	"github.com/andycai/goapi/modules/datacenter/db"
)

// CommandAddPost adds a new post
func CommandAddPost(post *models.Post) error {
	database := db.GetDB()
	// Check if slug already exists
	var count int64
	if err := database.Model(&models.Post{}).Where("slug = ?", post.Slug).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrPostSlugAlreadyExists
	}

	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()

	return database.Create(post).Error
}

// CommandUpdatePost updates an existing post
func CommandUpdatePost(post *models.Post) error {
	database := db.GetDB()
	// Check if post exists
	var count int64
	if err := database.Model(&models.Post{}).Where("id = ?", post.ID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPostNotFound
	}

	// Check if new slug conflicts with other posts
	if err := database.Model(&models.Post{}).Where("slug = ? AND id != ?", post.Slug, post.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrPostSlugAlreadyExists
	}

	post.UpdatedAt = time.Now()

	return database.Save(post).Error
}

// CommandDeletePost deletes a post
func CommandDeletePost(id int64) error {
	database := db.GetDB()
	// Check if post exists
	var count int64
	if err := database.Model(&models.Post{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPostNotFound
	}

	return database.Delete(&models.Post{}, id).Error
}
