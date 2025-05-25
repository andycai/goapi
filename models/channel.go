package models

import (
	"time"

	"gorm.io/gorm"
)

// Channel represents a game channel/platform
type Channel struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null;comment:渠道名称" json:"name"`
	ServerList    string         `gorm:"type:text;comment:服务器列表" json:"server_list"`
	CDNVersion    string         `gorm:"type:varchar(50);comment:CDN版本" json:"cdn_version"`
	CDNURL        string         `gorm:"type:varchar(255);comment:CDN地址" json:"cdn_url"`
	CDNURL2       string         `gorm:"type:varchar(255);comment:CDN备用地址" json:"cdn_url2"`
	OpenPatch     string         `gorm:"type:varchar(50);comment:开放补丁" json:"open_patch"`
	LoginAPI      string         `gorm:"type:varchar(255);comment:登录API" json:"login_api"`
	LoginURL      string         `gorm:"type:varchar(255);comment:登录地址" json:"login_url"`
	PkgVersion    string         `gorm:"type:varchar(50);comment:包版本" json:"pkg_version"`
	ServerListURL string         `gorm:"type:varchar(255);comment:服务器列表地址" json:"server_list_url"`
	NoticeURL     string         `gorm:"type:varchar(255);comment:公告地址" json:"notice_url"`
	NoticeNumURL  string         `gorm:"type:varchar(255);comment:公告数量地址" json:"notice_num_url"`
	ServerGroups  []ServerGroup  `gorm:"many2many:channel_server_groups;" json:"server_groups"`
	Announcements []Announcement `gorm:"many2many:channel_announcements;" json:"announcements"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// PhysicalServer represents a physical game server
type PhysicalServer struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ServerID     uint           `gorm:"not null;comment:服务器ID" json:"server_id"`
	Name         string         `gorm:"type:varchar(100);not null;comment:服务器名称" json:"name"`
	ServerStatus uint           `gorm:"not null;comment:服务器状态" json:"server_status"`
	Available    bool           `gorm:"not null;comment:是否可用" json:"available"`
	MergeID      uint           `gorm:"comment:合服ID" json:"merge_id"`
	Online       uint           `gorm:"comment:在线人数" json:"online"`
	ServerPort   uint           `gorm:"not null;comment:服务器端口" json:"server_port"`
	ServerIP     string         `gorm:"type:varchar(50);not null;comment:服务器IP" json:"server_ip"`
	CreatedAt    time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// ServerGroup represents a group of servers
type ServerGroup struct {
	ID        uint                `gorm:"primaryKey" json:"id"`
	Name      string              `gorm:"type:varchar(100);not null;comment:分组名称" json:"name"`
	Servers   []ServerGroupServer `gorm:"foreignKey:GroupID" json:"servers"`
	CreatedAt time.Time           `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time           `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt      `gorm:"index" json:"-"`
}

// ServerGroupServer represents a server in a group with custom properties
type ServerGroupServer struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	GroupID          uint           `gorm:"not null;comment:分组ID" json:"group_id"`
	ServerID         uint           `gorm:"not null;comment:服务器ID" json:"server_id"`
	Name             string         `gorm:"type:varchar(100);not null;comment:服务器名称" json:"name"`
	ServerStatus     uint           `gorm:"not null;comment:服务器状态" json:"server_status"`
	Available        bool           `gorm:"not null;comment:是否可用" json:"available"`
	MergeID          uint           `gorm:"comment:合服ID" json:"merge_id"`
	PhysicalServerID uint           `gorm:"not null;comment:物理服务器ID" json:"physical_server_id"`
	PhysicalServer   PhysicalServer `gorm:"foreignKey:PhysicalServerID" json:"physical_server"`
	CreatedAt        time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
}

// Announcement represents a game announcement
type Announcement struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"type:varchar(200);not null;comment:公告标题" json:"title"`
	Content   string         `gorm:"type:text;not null;comment:公告内容" json:"content"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies the table name for Channel
func (Channel) TableName() string {
	return "channels"
}

// TableName specifies the table name for PhysicalServer
func (PhysicalServer) TableName() string {
	return "physical_servers"
}

// TableName specifies the table name for ServerGroup
func (ServerGroup) TableName() string {
	return "server_groups"
}

// TableName specifies the table name for ServerGroupServer
func (ServerGroupServer) TableName() string {
	return "server_group_servers"
}

// TableName specifies the table name for Announcement
func (Announcement) TableName() string {
	return "announcements"
}
