package dict

import (
	"strconv"

	"github.com/andycai/unitool/models"
	"github.com/gofiber/fiber/v2"
)

// 字典类型处理器

// listDictTypeHandler 获取字典类型列表
func listDictTypeHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	dictTypes, total, err := getDictTypeList(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取字典类型列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"dictTypes": dictTypes,
		"total":     total,
	})
}

// addDictTypeHandler 添加字典类型
func addDictTypeHandler(c *fiber.Ctx) error {
	dictType := new(models.DictType)
	if err := c.BodyParser(dictType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if dictType.Name == "" || dictType.Type == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "名称和类型不能为空",
		})
	}

	if err := addDictType(dictType); err != nil {
		if err == ErrDictTypeAlreadyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加字典类型失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加字典类型成功",
		"id":      dictType.ID,
	})
}

// editDictTypeHandler 编辑字典类型
func editDictTypeHandler(c *fiber.Ctx) error {
	dictType := new(models.DictType)
	if err := c.BodyParser(dictType); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if dictType.ID == 0 || dictType.Name == "" || dictType.Type == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID、名称和类型不能为空",
		})
	}

	if err := updateDictType(dictType); err != nil {
		if err == ErrDictTypeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		if err == ErrDictTypeAlreadyExists {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新字典类型失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新字典类型成功",
	})
}

// deleteDictTypeHandler 删除字典类型
func deleteDictTypeHandler(c *fiber.Ctx) error {
	type DeleteRequest struct {
		ID int64 `json:"id"`
	}

	var req DeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if req.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的ID",
		})
	}

	if err := deleteDictType(req.ID); err != nil {
		if err == ErrDictTypeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除字典类型失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除字典类型成功",
	})
}

// 字典数据处理器

// listDictDataHandler 获取字典数据列表
func listDictDataHandler(c *fiber.Ctx) error {
	typeCode := c.Query("type")
	if typeCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "字典类型不能为空",
		})
	}

	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	// 首先检查字典类型是否存在
	_, err = getDictTypeByType(typeCode)
	if err != nil {
		if err == ErrDictTypeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "查询字典类型失败: " + err.Error(),
		})
	}

	dictData, total, err := getDictDataList(typeCode, page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取字典数据列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"dictData": dictData,
		"total":    total,
	})
}

// addDictDataHandler 添加字典数据
func addDictDataHandler(c *fiber.Ctx) error {
	dictData := new(models.DictData)
	if err := c.BodyParser(dictData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if dictData.Type == "" || dictData.Label == "" || dictData.Value == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "类型、标签和值不能为空",
		})
	}

	if err := addDictData(dictData); err != nil {
		if err == ErrDictTypeNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加字典数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加字典数据成功",
		"id":      dictData.ID,
	})
}

// editDictDataHandler 编辑字典数据
func editDictDataHandler(c *fiber.Ctx) error {
	dictData := new(models.DictData)
	if err := c.BodyParser(dictData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证必填字段
	if dictData.ID == 0 || dictData.Type == "" || dictData.Label == "" || dictData.Value == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID、类型、标签和值不能为空",
		})
	}

	if err := updateDictData(dictData); err != nil {
		if err == ErrDictDataNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新字典数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新字典数据成功",
	})
}

// deleteDictDataHandler 删除字典数据
func deleteDictDataHandler(c *fiber.Ctx) error {
	type DeleteRequest struct {
		ID int64 `json:"id"`
	}

	var req DeleteRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if req.ID <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的ID",
		})
	}

	if err := deleteDictData(req.ID); err != nil {
		if err == ErrDictDataNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除字典数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除字典数据成功",
	})
}
