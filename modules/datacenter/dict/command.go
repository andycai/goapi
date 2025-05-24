package dict

import (
	"github.com/andycai/goapi/models"
)

// CommandAddDictType adds a new dictionary type
func CommandAddDictType(dictType *models.DictType) error {
	return addDictType(dictType)
}

// CommandUpdateDictType updates an existing dictionary type
func CommandUpdateDictType(dictType *models.DictType) error {
	return updateDictType(dictType)
}

// CommandDeleteDictType deletes a dictionary type
func CommandDeleteDictType(id int64) error {
	return deleteDictType(id)
}

// CommandAddDictData adds a new dictionary data
func CommandAddDictData(dictData *models.DictData) error {
	return addDictData(dictData)
}

// CommandUpdateDictData updates an existing dictionary data
func CommandUpdateDictData(dictData *models.DictData) error {
	return updateDictData(dictData)
}

// CommandDeleteDictData deletes a dictionary data
func CommandDeleteDictData(id int64) error {
	return deleteDictData(id)
}
