package db

import (
	"nebula/internal/app"
	"nebula/internal/asset"
	"nebula/internal/auth"
	"nebula/internal/release"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func Init(dsn string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(
		&app.App{},
		&release.Release{},
		&asset.Asset{},
		&auth.User{},
	)

	return db
}
