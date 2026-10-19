package service

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/yuyudeqiu/chronicle/internal/model"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) error {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	err = db.AutoMigrate(&model.Task{}, &model.TaskLog{})
	if err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}

	DB = db
	return nil
}
