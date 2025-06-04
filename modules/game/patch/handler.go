package patch

import (
	"strconv"

	"github.com/andycai/goapi/pkg/utility/path"
	"github.com/gofiber/fiber/v2"
)

// saveConfigHandler 保存补丁配置
func saveConfigHandler(c *fiber.Ctx) error {
	config := new(PatchConfig)
	if err := c.BodyParser(config); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 验证路径安全性
	if !path.IsValid(config.PatchPath) || !path.IsValid(config.ZipPath) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的目录路径",
		})
	}

	if err := updateConfig(config); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "保存配置失败",
		})
	}

	return c.JSON(fiber.Map{
		"message": "配置保存成功",
	})
}

// getConfigHandler 获取补丁配置
func getConfigHandler(c *fiber.Ctx) error {
	config := getConfig()
	return c.JSON(config)
}

// generatePatchHandler 生成补丁包
func generatePatchHandler(c *fiber.Ctx) error {
	type GenerateRequest struct {
		OldVersion  string `json:"old_version"`
		NewVersion  string `json:"new_version"`
		Description string `json:"description"`
		Branch      string `json:"branch"`
		Platform    string `json:"platform"`
	}

	var req GenerateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 生成补丁包
	record, err := GeneratePatch(req.OldVersion, req.NewVersion, req.Description, req.Branch, req.Platform)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "生成补丁包失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "补丁包生成成功",
		"record":  record,
	})
}

// listPatchRecordsHandler 获取补丁记录列表
func listPatchRecordsHandler(c *fiber.Ctx) error {
	// 获取分页参数
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	records, totalCount, err := GetPatchRecords(limit, page)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取补丁记录失败",
		})
	}

	return c.JSON(fiber.Map{
		"records": records,
		"total":   totalCount,
	})
}

// applyPatchHandler 应用补丁包
func applyPatchHandler(c *fiber.Ctx) error {
	type ApplyRequest struct {
		RecordID uint `json:"record_id"`
	}

	var req ApplyRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	// 应用补丁包
	if err := ApplyPatch(req.RecordID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "应用补丁包失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "补丁包应用成功",
	})
}
