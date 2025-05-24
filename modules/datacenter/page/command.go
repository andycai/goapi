package page

import (
	"time"

	"github.com/andycai/goapi/models"
)

// CommandAddPage adds a new page
func CommandAddPage(page *models.Page) error {
	// Check if slug already exists
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

// CommandUpdatePage updates an existing page
func CommandUpdatePage(page *models.Page) error {
	// Check if page exists
	var count int64
	if err := app.DB.Model(&models.Page{}).Where("id = ?", page.ID).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPageNotFound
	}

	// Check if new slug conflicts with other pages
	if err := app.DB.Model(&models.Page{}).Where("slug = ? AND id != ?", page.Slug, page.ID).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return ErrPageSlugAlreadyExists
	}

	page.UpdatedAt = time.Now()

	return app.DB.Save(page).Error
}

// CommandDeletePage deletes a page
func CommandDeletePage(id int64) error {
	// Check if page exists
	var count int64
	if err := app.DB.Model(&models.Page{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return err
	}

	if count == 0 {
		return ErrPageNotFound
	}

	return app.DB.Delete(&models.Page{}, id).Error
}
