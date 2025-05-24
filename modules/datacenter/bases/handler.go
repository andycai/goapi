package bases

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// listEntityHandler 获取实体列表
func listEntityHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	search := c.Query("search", "")

	entities, total, err := QueryEntities(limit, page, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取实体列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"entities": entities,
		"total":    total,
	})
}

// getEntityHandler 获取单个实体
func getEntityHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的实体ID",
		})
	}

	entity, err := QueryEntity(uint(id))
	if err != nil {
		if err == ErrEntityNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取实体失败: " + err.Error(),
		})
	}

	return c.JSON(entity)
}

// addEntityHandler 添加实体
func addEntityHandler(c *fiber.Ctx) error {
	req := new(EntityRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "未登录",
		})
	}

	entity, err := CommandCreateEntity(req, currentUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加实体失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加实体成功",
		"entity":  entity,
	})
}

// editEntityHandler 编辑实体
func editEntityHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的实体ID",
		})
	}

	req := new(EntityRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "未登录",
		})
	}

	entity, err := CommandUpdateEntity(uint(id), req, currentUser.ID)
	if err != nil {
		if err == ErrEntityNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新实体失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新实体成功",
		"entity":  entity,
	})
}

// deleteEntityHandler 删除实体
func deleteEntityHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的实体ID",
		})
	}

	if err := CommandDeleteEntity(uint(id)); err != nil {
		if err == ErrEntityNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除实体失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除实体成功",
	})
}

// listFieldHandler 获取字段列表
func listFieldHandler(c *fiber.Ctx) error {
	// 获取分页参数
	entityID, err := strconv.ParseUint(c.Query("entity_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的实体ID",
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	search := c.Query("search", "")

	fields, total, err := QueryFields(uint(entityID), limit, page, search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取字段列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"fields": fields,
		"total":  total,
	})
}

// getFieldHandler 获取单个字段
func getFieldHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的字段ID",
		})
	}

	field, err := QueryField(uint(id))
	if err != nil {
		if err == ErrFieldNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取字段失败: " + err.Error(),
		})
	}

	return c.JSON(field)
}

// addFieldHandler 添加字段
func addFieldHandler(c *fiber.Ctx) error {
	req := new(FieldRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "未登录",
		})
	}

	field, err := CommandCreateField(req, currentUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加字段失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加字段成功",
		"field":   field,
	})
}

// editFieldHandler 编辑字段
func editFieldHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的字段ID",
		})
	}

	req := new(FieldRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "未登录",
		})
	}

	field, err := CommandUpdateField(uint(id), req, currentUser.ID)
	if err != nil {
		if err == ErrFieldNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新字段失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新字段成功",
		"field":   field,
	})
}

// deleteFieldHandler 删除字段
func deleteFieldHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的字段ID",
		})
	}

	if err := CommandDeleteField(uint(id)); err != nil {
		if err == ErrFieldNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除字段失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除字段成功",
	})
}

// listEntityDataHandler 获取实体数据列表
func listEntityDataHandler(c *fiber.Ctx) error {
	// 获取分页参数
	entityID, err := strconv.ParseUint(c.Query("entity_id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的实体ID",
		})
	}

	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	data, total, err := QueryEntityData(uint(entityID), limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取实体数据列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data":  data,
		"total": total,
	})
}

// getEntityDataHandler 获取单个实体数据
func getEntityDataHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的数据ID",
		})
	}

	data, err := QueryEntityDataByID(uint(id))
	if err != nil {
		if err == ErrDataNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取实体数据失败: " + err.Error(),
		})
	}

	return c.JSON(data)
}

// addEntityDataHandler 添加实体数据
func addEntityDataHandler(c *fiber.Ctx) error {
	req := new(EntityDataRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "未登录",
		})
	}

	data, err := CommandCreateEntityData(req, currentUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加实体数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加实体数据成功",
		"data":    data,
	})
}

// editEntityDataHandler 编辑实体数据
func editEntityDataHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的数据ID",
		})
	}

	req := new(EntityDataRequest)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 获取当前用户ID
	currentUser := app.CurrentUser(c)
	if currentUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "未登录",
		})
	}

	data, err := CommandUpdateEntityData(uint(id), req, currentUser.ID)
	if err != nil {
		if err == ErrDataNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新实体数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新实体数据成功",
		"data":    data,
	})
}

// deleteEntityDataHandler 删除实体数据
func deleteEntityDataHandler(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的数据ID",
		})
	}

	if err := CommandDeleteEntityData(uint(id)); err != nil {
		if err == ErrDataNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除实体数据失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除实体数据成功",
	})
}
