package user

import "github.com/andycai/goapi/models"

// 获取用户列表
func QueryUserList(page, limit int) ([]models.User, int64, error) {
	return getUserList(page, limit)
}
