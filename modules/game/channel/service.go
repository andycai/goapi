package channel

import (
	"errors"

	"gorm.io/gorm"
)

// initService 初始化服务
func initService() {
	// 初始化服务逻辑
}

// Channel Service operations
func CreateChannelWithRelations(channel *Channel) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(channel).Error; err != nil {
			return err
		}

		// 更新关联关系
		if len(channel.ServerGroups) > 0 {
			if err := tx.Model(channel).Association("ServerGroups").Replace(channel.ServerGroups); err != nil {
				return err
			}
		}

		if len(channel.Announcements) > 0 {
			if err := tx.Model(channel).Association("Announcements").Replace(channel.Announcements); err != nil {
				return err
			}
		}

		return nil
	})
}

func UpdateChannelWithRelations(id uint, data map[string]interface{}) error {
	channel, err := GetChannelByID(id)
	if err != nil {
		return err
	}

	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(channel).Updates(data).Error; err != nil {
			return err
		}

		// 更新关联关系
		if serverGroups, ok := data["server_groups"]; ok {
			if err := tx.Model(channel).Association("ServerGroups").Replace(serverGroups); err != nil {
				return err
			}
		}

		if announcements, ok := data["announcements"]; ok {
			if err := tx.Model(channel).Association("Announcements").Replace(announcements); err != nil {
				return err
			}
		}

		return nil
	})
}

// ServerGroup Service operations
func CreateServerGroupWithServers(group *ServerGroup) error {
	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(group).Error; err != nil {
			return err
		}

		// 添加服务器到分组
		for _, server := range group.Servers {
			server.GroupID = group.ID
			if err := tx.Create(&server).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func UpdateServerGroupWithServers(id uint, data map[string]interface{}) error {
	group, err := GetServerGroupByID(id)
	if err != nil {
		return err
	}

	return app.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(group).Updates(data).Error; err != nil {
			return err
		}

		// 更新服务器列表
		if servers, ok := data["servers"]; ok {
			// 删除旧的服务器关联
			if err := tx.Where("group_id = ?", id).Delete(&ServerGroupServer{}).Error; err != nil {
				return err
			}

			// 添加新的服务器关联
			for _, server := range servers.([]ServerGroupServer) {
				server.GroupID = id
				if err := tx.Create(&server).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// ServerGroupServer Service operations
// func AddServerToGroupWithValidation(groupID uint, server *ServerGroupServer) error {
// 	// 验证物理服务器是否存在
// 	if _, err := GetPhysicalServerByID(server.PhysicalServerID); err != nil {
// 		return errors.New("物理服务器不存在")
// 	}

// 	// 验证服务器是否已经在分组中
// 	var count int64
// 	if err := app.DB.Model(&ServerGroupServer{}).
// 		Where("group_id = ? AND physical_server_id = ?", groupID, server.PhysicalServerID).
// 		Count(&count).Error; err != nil {
// 		return err
// 	}

// 	if count > 0 {
// 		return errors.New("服务器已存在于该分组中")
// 	}

// 	return AddServerToGroup(groupID, server.PhysicalServerID)
// }

// ServerGroupServer Service operations
func AddServerToGroupWithValidation(groupID uint, physicalServerId uint) error {
	// 验证物理服务器是否存在
	physicalServer, err := GetPhysicalServerByID(physicalServerId)
	if err != nil {
		return errors.New("物理服务器不存在")
	}

	// 验证服务器是否已经在分组中
	var count int64
	if err := app.DB.Model(&ServerGroupServer{}).
		Where("group_id = ? AND physical_server_id = ?", groupID, physicalServerId).
		Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("服务器已存在于该分组中")
	}

	return AddServerToGroup(groupID, physicalServer)
}

// Announcement Service operations
func CreateAnnouncementWithValidation(announcement *Announcement) error {
	if announcement.Title == "" {
		return errors.New("公告标题不能为空")
	}

	if announcement.Content == "" {
		return errors.New("公告内容不能为空")
	}

	return CreateAnnouncement(announcement)
}

func UpdateAnnouncementWithValidation(id uint, data map[string]interface{}) error {
	if title, ok := data["title"]; ok && title == "" {
		return errors.New("公告标题不能为空")
	}

	if content, ok := data["content"]; ok && content == "" {
		return errors.New("公告内容不能为空")
	}

	return UpdateAnnouncement(id, data)
}
