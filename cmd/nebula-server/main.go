package main

import (
	"log"
	"nebula/internal/api"
	"nebula/internal/auth"
	"nebula/internal/config"
	"nebula/internal/db"
	"nebula/internal/storage"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	database := db.Init(cfg.Database.DSN)

	// 初始化存储
	var stor storage.Storage
	var err error

	switch cfg.Storage.Type {
	case "local":
		stor, err = storage.NewLocalStorage(cfg.Storage.BasePath, cfg.Storage.BaseURL)
		if err != nil {
			log.Fatalf("failed to initialize storage: %v", err)
		}
	default:
		log.Fatalf("unsupported storage type: %s", cfg.Storage.Type)
	}

	// 初始化 JWT 服务
	jwtService := auth.NewJWTService(
		cfg.JWT.Secret,
		cfg.GetAccessTokenDuration(),
		cfg.GetRefreshTokenDuration(),
	)

	// 初始化认证服务（使用配置的管理员账号）
	authService := auth.NewAuthService(
		database,
		jwtService,
		cfg.Admin.Username,
		cfg.Admin.Password,
	)

	r := gin.Default()

	// 静态文件服务 - 用于提供文件下载
	r.Static("/files", cfg.Storage.BasePath)

	// 注册 API 路由
	api.RegisterRoutes(r, database, stor, jwtService, authService)

	// 生产环境：提供前端静态资源
	if cfg.Frontend.Enabled {
		r.Static("/assets", cfg.Frontend.Path+"/assets")
		r.StaticFile("/favicon.ico", cfg.Frontend.Path+"/favicon.ico")

		// 处理 SPA 路由 - 所有未匹配的路由都返回 index.html
		r.NoRoute(func(c *gin.Context) {
			c.File(cfg.Frontend.Path + "/index.html")
		})

		log.Printf("Frontend static files enabled at %s", cfg.Frontend.Path)
		log.Printf("Frontend available at http://%s", cfg.Server.Address)
	} else {
		log.Printf("Running in dev mode - frontend should be served by Vite dev server")
	}

	log.Printf("Server starting on %s (mode: %s)", cfg.Server.Address, cfg.Server.Mode)
	log.Printf("Admin username: %s", cfg.Admin.Username)
	log.Printf("API available at http://%s/api", cfg.Server.Address)
	r.Run(cfg.Server.Address)
}
