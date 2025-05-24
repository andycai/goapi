package parameter

import (
	"errors"

	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 获取参数列表
func getParameters(limit, page int, search string) ([]models.Parameter, int64, error) {
	var parameters []models.Parameter
	var total int64
	offset := (page - 1) * limit

	query := app.DB.Model(&models.Parameter{})

	// 如果有搜索条件
	if search != "" {
		query = query.Where("name LIKE ? OR type LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	if err := query.Limit(limit).Offset(offset).Order("id DESC").Find(&parameters).Error; err != nil {
		return nil, 0, err
	}

	return parameters, total, nil
}

// 获取单个参数
func getParameter(id uint) (models.Parameter, error) {
	var param models.Parameter
	if err := app.DB.First(&param, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return param, errors.New("参数不存在")
		}
		return param, err
	}
	return param, nil
}

// 创建参数
func createParameter(param *models.Parameter) error {
	return app.DB.Create(param).Error
}

// 更新参数
func updateParameter(param *models.Parameter) error {
	return app.DB.Save(param).Error
}

// 删除参数
func deleteParameter(id uint) error {
	// 查找参数
	var param models.Parameter
	if err := app.DB.First(&param, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("参数不存在")
		}
		return err
	}

	return app.DB.Delete(&param).Error
}
