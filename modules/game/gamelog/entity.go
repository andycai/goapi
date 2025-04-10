package gamelog

import "github.com/andycai/goapi/models"

// LogReq 日志请求结构体
type LogReq struct {
	AppID    string           `json:"app_id"`
	Package  string           `json:"package"`
	RoleName string           `json:"role_name"`
	Device   string           `json:"device"`
	Logs     []models.GameLog `json:"list"`
}
