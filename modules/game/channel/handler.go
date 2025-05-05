package channel

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Channel Handlers
func getChannelsHandler(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	channels, total, err := GetChannels(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取渠道列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"channels": channels,
		"total":    total,
	})
}

func createChannelHandler(c *fiber.Ctx) error {
	var channel Channel
	if err := c.BodyParser(&channel); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := CreateChannelWithRelations(&channel); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建渠道失败: " + err.Error(),
		})
	}

	return c.JSON(channel)
}

func updateChannelHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的渠道ID",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := UpdateChannelWithRelations(uint(id), data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新渠道失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新成功",
	})
}

func deleteChannelHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的渠道ID",
		})
	}

	if err := DeleteChannel(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除渠道失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}

// PhysicalServer Handlers
func getPhysicalServersHandler(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	servers, total, err := GetPhysicalServers(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取服务器列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"servers": servers,
		"total":   total,
	})
}

func createPhysicalServerHandler(c *fiber.Ctx) error {
	var server PhysicalServer
	if err := c.BodyParser(&server); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := CreatePhysicalServer(&server); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建服务器失败: " + err.Error(),
		})
	}

	return c.JSON(server)
}

func updatePhysicalServerHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的服务器ID",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := UpdatePhysicalServer(uint(id), data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新服务器失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新成功",
	})
}

func deletePhysicalServerHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的服务器ID",
		})
	}

	if err := DeletePhysicalServer(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除服务器失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}

// ServerGroup Handlers
func getServerGroupsHandler(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	groups, total, err := GetServerGroups(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取服务器分组列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"groups": groups,
		"total":  total,
	})
}

func createServerGroupHandler(c *fiber.Ctx) error {
	var group ServerGroup
	if err := c.BodyParser(&group); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := CreateServerGroupWithServers(&group); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建服务器分组失败: " + err.Error(),
		})
	}

	return c.JSON(group)
}

func updateServerGroupHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的分组ID",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := UpdateServerGroupWithServers(uint(id), data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新服务器分组失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新成功",
	})
}

func deleteServerGroupHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的分组ID",
		})
	}

	if err := DeleteServerGroup(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除服务器分组失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}

func getServerGroupServersHandler(c *fiber.Ctx) error {
	groupId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的分组ID",
		})
	}

	servers, err := GetServerGroupServers(uint(groupId))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取服务器分组中的服务器失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": servers,
	})
}

func addServerToGroupHandler(c *fiber.Ctx) error {
	groupId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的分组ID",
		})
	}

	serverId, err := strconv.Atoi(c.Params("serverId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的服务器ID",
		})
	}

	if err := AddServerToGroupWithValidation(uint(groupId), uint(serverId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "添加服务器到分组失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "添加成功",
	})
}

func removeServerFromGroupHandler(c *fiber.Ctx) error {
	groupId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的分组ID",
		})
	}

	serverId, err := strconv.Atoi(c.Params("serverId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的服务器ID",
		})
	}

	if err := RemoveServerFromGroup(uint(groupId), uint(serverId)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "从分组中移除服务器失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "移除成功",
	})
}

// Announcement Handlers
func getAnnouncementsHandler(c *fiber.Ctx) error {
	limit, err := strconv.Atoi(c.Query("limit", "10"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	page, err := strconv.Atoi(c.Query("page", "1"))
	if err != nil || page <= 0 {
		page = 1
	}

	announcements, total, err := GetAnnouncements(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "获取公告列表失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"announcements": announcements,
		"total":         total,
	})
}

func createAnnouncementHandler(c *fiber.Ctx) error {
	var announcement Announcement
	if err := c.BodyParser(&announcement); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := CreateAnnouncementWithValidation(&announcement); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "创建公告失败: " + err.Error(),
		})
	}

	return c.JSON(announcement)
}

func updateAnnouncementHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的公告ID",
		})
	}

	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的请求数据",
		})
	}

	if err := UpdateAnnouncementWithValidation(uint(id), data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "更新公告失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "更新成功",
	})
}

func deleteAnnouncementHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "无效的公告ID",
		})
	}

	if err := DeleteAnnouncement(uint(id)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "删除公告失败: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "删除成功",
	})
}
