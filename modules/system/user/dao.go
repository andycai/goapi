package user

import "github.com/andycai/goapi/models"

func getUserList(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	db := app.DB.Model(&models.User{})
	db.Count(&total)

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		db = db.Offset(offset).Limit(limit)
	}

	if err := db.Preload("Role").Order("id desc").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
