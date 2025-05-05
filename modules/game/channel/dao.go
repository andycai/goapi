package channel

import (
	"log"
	"time"

	"github.com/andycai/goapi/enum"
	"github.com/andycai/goapi/models"
	"gorm.io/gorm"
)

// 数据访问层

func autoMigrate() error {
	return app.DB.AutoMigrate(
		&Channel{},
		&PhysicalServer{},
		&ServerGroup{},
		&ServerGroupServer{},
		&Announcement{},
	)
}

// 初始化数据
func initData() error {
	if err := initMenus(); err != nil {
		return err
	}

	if err := initPermissions(); err != nil {
		return err
	}

	return nil
}

func initMenus() error {
	// 检查是否已初始化
	if app.IsInitializedModule("channel:menu") {
		log.Println("[渠道模块]菜单数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建渠道菜单
		channelMenus := []*models.Menu{
			{
				MenuID:     3002,
				ParentID:   enum.MenuIdGame,
				Name:       "渠道管理",
				Path:       "/admin/channel",
				Icon:       "channel",
				Sort:       2,
				Permission: "channel:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				MenuID:     3012,
				ParentID:   enum.MenuIdGame,
				Name:       "物理服务器",
				Path:       "/admin/physical_servers",
				Icon:       "physical_server",
				Sort:       12,
				Permission: "server:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				MenuID:     3013,
				ParentID:   enum.MenuIdGame,
				Name:       "服务器分组",
				Path:       "/admin/server_groups",
				Icon:       "server_group",
				Sort:       13,
				Permission: "server:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
			{
				MenuID:     3014,
				ParentID:   enum.MenuIdGame,
				Name:       "公告管理",
				Path:       "/admin/announcement",
				Icon:       "announcement",
				Sort:       14,
				Permission: "announcement:view",
				IsShow:     true,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			},
		}

		if err := tx.CreateInBatches(channelMenus, len(channelMenus)).Error; err != nil {
			return err
		}

		// 标记菜单已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "channel:menu",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

func initPermissions() error {
	// 检查是否已初始化
	if app.IsInitializedModule("channel:permission") {
		log.Println("[渠道模块]权限数据已初始化，跳过")
		return nil
	}

	// 开始事务
	return app.DB.Transaction(func(tx *gorm.DB) error {
		// 创建渠道相关权限
		permissions := []models.Permission{
			{
				Name:        "渠道查看",
				Code:        "channel:view",
				Description: "查看渠道列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "渠道管理",
				Code:        "channel:manage",
				Description: "管理渠道（创建、编辑等）",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器查看",
				Code:        "server:view",
				Description: "查看服务器列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "服务器管理",
				Code:        "server:manage",
				Description: "管理服务器（创建、编辑等）",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "公告查看",
				Code:        "announcement:view",
				Description: "查看公告列表",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
			{
				Name:        "公告管理",
				Code:        "announcement:manage",
				Description: "管理公告（创建、编辑等）",
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		}

		if err := tx.Create(&permissions).Error; err != nil {
			return err
		}

		// 标记模块已初始化
		if err := tx.Create(&models.ModuleInit{
			Module:      "channel:permission",
			Initialized: 1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// Channel DAO operations
func GetChannels(page, limit int) ([]Channel, int64, error) {
	var channels []Channel
	var total int64

	if err := app.DB.Model(&Channel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Preload("ServerGroups").Preload("Announcements").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&channels).Error; err != nil {
		return nil, 0, err
	}

	return channels, total, nil
}

func CreateChannel(channel *Channel) error {
	return app.DB.Create(channel).Error
}

func GetChannelByID(id uint) (*Channel, error) {
	var channel Channel
	if err := app.DB.Preload("ServerGroups").Preload("Announcements").
		First(&channel, id).Error; err != nil {
		return nil, err
	}
	return &channel, nil
}

func UpdateChannel(id uint, data map[string]interface{}) error {
	return app.DB.Model(&Channel{}).Where("id = ?", id).Updates(data).Error
}

func DeleteChannel(id uint) error {
	return app.DB.Delete(&Channel{}, id).Error
}

// PhysicalServer DAO operations
func GetPhysicalServers(page, limit int) ([]PhysicalServer, int64, error) {
	var servers []PhysicalServer
	var total int64

	if err := app.DB.Model(&PhysicalServer{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&servers).Error; err != nil {
		return nil, 0, err
	}

	return servers, total, nil
}

func CreatePhysicalServer(server *PhysicalServer) error {
	return app.DB.Create(server).Error
}

func GetPhysicalServerByID(id uint) (*PhysicalServer, error) {
	var server PhysicalServer
	if err := app.DB.First(&server, id).Error; err != nil {
		return nil, err
	}
	return &server, nil
}

func UpdatePhysicalServer(id uint, data map[string]interface{}) error {
	return app.DB.Model(&PhysicalServer{}).Where("id = ?", id).Updates(data).Error
}

func DeletePhysicalServer(id uint) error {
	return app.DB.Delete(&PhysicalServer{}, id).Error
}

// ServerGroup DAO operations
func GetServerGroups(page, limit int) ([]ServerGroup, int64, error) {
	var groups []ServerGroup
	var total int64

	if err := app.DB.Model(&ServerGroup{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Preload("Servers.PhysicalServer").
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&groups).Error; err != nil {
		return nil, 0, err
	}

	return groups, total, nil
}

func CreateServerGroup(group *ServerGroup) error {
	return app.DB.Create(group).Error
}

func GetServerGroupByID(id uint) (*ServerGroup, error) {
	var group ServerGroup
	if err := app.DB.Preload("Servers.PhysicalServer").
		First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

func UpdateServerGroup(id uint, data map[string]interface{}) error {
	return app.DB.Model(&ServerGroup{}).Where("id = ?", id).Updates(data).Error
}

func DeleteServerGroup(id uint) error {
	return app.DB.Delete(&ServerGroup{}, id).Error
}

// ServerGroupServer DAO operations
func GetServerGroupServers(groupId uint) ([]*ServerGroupServer, error) {
	var servers []*ServerGroupServer
	if err := app.DB.Where("group_id = ?", groupId).Find(&servers).Error; err != nil {
		return nil, err
	}
	return servers, nil
}

func AddServerToGroup(groupId uint, server *PhysicalServer) error {
	serverGroupServer := &ServerGroupServer{
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
	return app.DB.Where("group_id = ? AND server_id = ?", groupId, serverId).Delete(&ServerGroupServer{}).Error
}

// Announcement DAO operations
func GetAnnouncements(page, limit int) ([]Announcement, int64, error) {
	var announcements []Announcement
	var total int64

	if err := app.DB.Model(&Announcement{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := app.DB.Order("created_at DESC").Offset(offset).Limit(limit).Find(&announcements).Error; err != nil {
		return nil, 0, err
	}

	return announcements, total, nil
}

func CreateAnnouncement(announcement *Announcement) error {
	return app.DB.Create(announcement).Error
}

func GetAnnouncementByID(id uint) (*Announcement, error) {
	var announcement Announcement
	if err := app.DB.First(&announcement, id).Error; err != nil {
		return nil, err
	}
	return &announcement, nil
}

func UpdateAnnouncement(id uint, data map[string]interface{}) error {
	return app.DB.Model(&Announcement{}).Where("id = ?", id).Updates(data).Error
}

func DeleteAnnouncement(id uint) error {
	return app.DB.Delete(&Announcement{}, id).Error
}
