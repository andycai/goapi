package dict

import (
	"github.com/andycai/goapi/models"
)

// 获取字典类型列表
func QueryDictTypeList(page, limit int) ([]models.DictType, int64, error) {
	return getDictTypeList(page, limit)
}

// 根据类型编码获取字典类型
func QueryDictTypeByType(typeCode string) (models.DictType, error) {
	return getDictTypeByType(typeCode)
}

// 根据字典类型获取字典数据列表
func QueryDictDataList(typeCode string, page, limit int) ([]models.DictData, int64, error) {
	return getDictDataList(typeCode, page, limit)
}

// 获取所有字典数据（不分页）
func QueryAllDictData(typeCode string) ([]models.DictData, error) {
	return getAllDictData(typeCode)
}
