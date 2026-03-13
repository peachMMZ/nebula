package db

import (
	"nebula/internal/app"
	"nebula/internal/asset"
	"nebula/internal/release"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Init(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// 自动迁移数据库表
	db.AutoMigrate(
		&app.App{},
		&release.Release{},
		&asset.Asset{},
	)

	return db
}
