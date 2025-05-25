package channel

import "github.com/andycai/goapi/models"

// Channel DAO operations
func GetChannels(page, limit int) ([]models.Channel, int64, error) {
	var channels []models.Channel
	var total int64

	if err := app.DB.Model(&models.Channel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Preload("ServerGroups").Preload("Announcements").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&channels).Error; err != nil {
		return nil, 0, err
	}

	return channels, total, nil
}

func CreateChannel(channel *models.Channel) error {
	return app.DB.Create(channel).Error
}

func GetChannelByID(id uint) (*models.Channel, error) {
	var channel models.Channel
	if err := app.DB.Preload("ServerGroups").Preload("Announcements").
		First(&channel, id).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func UpdateChannel(id uint, data map[string]interface{}) error {
	return app.DB.Model(&models.Channel{}).Where("id = ?", id).Updates(data).Error
}

func DeleteChannel(id uint) error {
	return app.DB.Delete(&models.Channel{}, id).Error
}

// PhysicalServer DAO operations
func GetPhysicalServers(page, limit int) ([]models.PhysicalServer, int64, error) {
	var servers []models.PhysicalServer
	var total int64

	if err := app.DB.Model(&models.PhysicalServer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&servers).Error; err != nil {
		return nil, 0, err
	}

	return servers, total, nil
}

func CreatePhysicalServer(server *models.PhysicalServer) error {
	return app.DB.Create(server).Error
}

func GetPhysicalServerByID(id uint) (*models.PhysicalServer, error) {
	var server models.PhysicalServer
	if err := app.DB.First(&server, id).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

func UpdatePhysicalServer(id uint, data map[string]interface{}) error {
	return app.DB.Model(&models.PhysicalServer{}).Where("id = ?", id).Updates(data).Error
}

func DeletePhysicalServer(id uint) error {
	return app.DB.Delete(&models.PhysicalServer{}, id).Error
}

// ServerGroup DAO operations
func GetServerGroups(page, limit int) ([]models.ServerGroup, int64, error) {
	var groups []models.ServerGroup
	var total int64

	if err := app.DB.Model(&models.ServerGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Preload("Servers.PhysicalServer").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func CreateServerGroup(group *models.ServerGroup) error {
	return app.DB.Create(group).Error
}

func GetServerGroupByID(id uint) (*models.ServerGroup, error) {
	var group models.ServerGroup
	if err := app.DB.Preload("Servers.PhysicalServer").
		First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func UpdateServerGroup(id uint, data map[string]interface{}) error {
	return app.DB.Model(&models.ServerGroup{}).Where("id = ?", id).Updates(data).Error
}

func DeleteServerGroup(id uint) error {
	return app.DB.Delete(&models.ServerGroup{}, id).Error
}

// ServerGroupServer DAO operations
func GetServerGroupServers(groupId uint) ([]*models.ServerGroupServer, error) {
	var servers []*models.ServerGroupServer
	if err := app.DB.Where("group_id = ?", groupId).Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

func AddServerToGroup(groupId uint, server *models.PhysicalServer) error {
	serverGroupServer := &models.ServerGroupServer{
		GroupID:          groupId,
		ServerID:         server.ServerID,
		Name:             server.Name,
		ServerStatus:     server.ServerStatus,
		Available:        server.Available,
		MergeID:          server.MergeID,
		PhysicalServerID: server.ID,
		PhysicalServer:   *server,
	}
	return app.DB.Create(serverGroupServer).Error
}

func RemoveServerFromGroup(groupId, serverId uint) error {
	return app.DB.Where("group_id = ? AND server_id = ?", groupId, serverId).Delete(&models.ServerGroupServer{}).Error
}

// Announcement DAO operations
func GetAnnouncements(page, limit int) ([]models.Announcement, int64, error) {
	var announcements []models.Announcement
	var total int64

	if err := app.DB.Model(&models.Announcement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&announcements).Error; err != nil {
		return nil, 0, err
	}

	return announcements, total, nil
}

func CreateAnnouncement(announcement *models.Announcement) error {
	return app.DB.Create(announcement).Error
}

func GetAnnouncementByID(id uint) (*models.Announcement, error) {
	var announcement models.Announcement
	if err := app.DB.First(&announcement, id).Error; err != nil {
		return nil, err
	}
	return &announcement, nil
}

func UpdateAnnouncement(id uint, data map[string]interface{}) error {
	return app.DB.Model(&models.Announcement{}).Where("id = ?", id).Updates(data).Error
}

func DeleteAnnouncement(id uint) error {
	return app.DB.Delete(&models.Announcement{}, id).Error
}
