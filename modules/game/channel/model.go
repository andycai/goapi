package channel

import (
	"time"

	"gorm.io/gorm"
)

// Channel represents a game channel/platform
type Channel struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"type:varchar(100);not null;comment:渠道名称" json:"name"`
	ServerList    string         `gorm:"type:text;comment:服务器列表" json:"serverList"`
	CDNVersion    string         `gorm:"type:varchar(50);comment:CDN版本" json:"cdnVersion"`
	CDNURL        string         `gorm:"type:varchar(255);comment:CDN地址" json:"cdnUrl"`
	CDNURL2       string         `gorm:"type:varchar(255);comment:CDN备用地址" json:"cdnUrl2"`
	OpenPatch     string         `gorm:"type:varchar(50);comment:开放补丁" json:"openPatch"`
	LoginAPI      string         `gorm:"type:varchar(255);comment:登录API" json:"loginApi"`
	LoginURL      string         `gorm:"type:varchar(255);comment:登录地址" json:"loginUrl"`
	PkgVersion    string         `gorm:"type:varchar(50);comment:包版本" json:"pkgVersion"`
	ServerListURL string         `gorm:"type:varchar(255);comment:服务器列表地址" json:"serverListUrl"`
	NoticeURL     string         `gorm:"type:varchar(255);comment:公告地址" json:"noticeUrl"`
	NoticeNumURL  string         `gorm:"type:varchar(255);comment:公告数量地址" json:"noticeNumUrl"`
	ServerGroups  []ServerGroup  `gorm:"many2many:channel_server_groups;" json:"serverGroups"`
	Announcements []Announcement `gorm:"many2many:channel_announcements;" json:"announcements"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// PhysicalServer represents a physical game server
type PhysicalServer struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ServerID     uint           `gorm:"not null;comment:服务器ID" json:"serverId"`
	Name         string         `gorm:"type:varchar(100);not null;comment:服务器名称" json:"name"`
	ServerStatus uint           `gorm:"not null;comment:服务器状态" json:"serverStatus"`
	Available    bool           `gorm:"not null;comment:是否可用" json:"available"`
	MergeID      uint           `gorm:"comment:合服ID" json:"mergeId"`
	Online       uint           `gorm:"comment:在线人数" json:"online"`
	ServerPort   uint           `gorm:"not null;comment:服务器端口" json:"serverPort"`
	ServerIP     string         `gorm:"type:varchar(50);not null;comment:服务器IP" json:"serverIp"`
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
	GroupID          uint           `gorm:"not null;comment:分组ID" json:"groupId"`
	ServerID         uint           `gorm:"not null;comment:服务器ID" json:"serverId"`
	Name             string         `gorm:"type:varchar(100);not null;comment:服务器名称" json:"name"`
	ServerStatus     uint           `gorm:"not null;comment:服务器状态" json:"serverStatus"`
	Available        bool           `gorm:"not null;comment:是否可用" json:"available"`
	MergeID          uint           `gorm:"comment:合服ID" json:"mergeId"`
	PhysicalServerID uint           `gorm:"not null;comment:物理服务器ID" json:"physicalServerId"`
	PhysicalServer   PhysicalServer `gorm:"foreignKey:PhysicalServerID" json:"physicalServer"`
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
