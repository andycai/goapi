package database

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func InitDatabase(dsn string, driver string) (*gorm.DB, error) {

	var db *gorm.DB
	var err error

	switch driver {
	case "sqlite":
		// 确保数据库目录存在
		dbDir := filepath.Dir(dsn)
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			return nil, fmt.Errorf("创建数据库目录失败: %v", err)
		}

		db, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("连接SQLite数据库失败: %v", err)
		}

	default:
		return nil, fmt.Errorf("不支持的数据库驱动: %s", driver)
	}

	return db, nil
}
