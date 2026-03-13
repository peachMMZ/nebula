package api

import (
	"nebula/internal/api/handler"
	"nebula/internal/app"
	"nebula/internal/asset"
	"nebula/internal/auth"
	"nebula/internal/release"
	"nebula/internal/storage"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(r *gin.Engine, db *gorm.DB, stor storage.Storage, jwtService *auth.JWTService, authService *auth.AuthService) {
	h := handler.New(db)
	appHandler := handler.NewAppHandler(app.NewService(db))
	releaseHandler := handler.NewReleaseHandler(release.NewService(db))
	assetHandler := handler.NewAssetHandler(asset.NewService(db, stor))
	authHandler := handler.NewAuthHandler(authService)

	api := r.Group("/api")

	{
		// 认证相关路由（公开）
		authGroup := api.Group("/auth")
		{
			// 移除注册接口 - 只允许配置的管理员登录
			// authGroup.POST("/register", authHandler.Register)
			authGroup.POST("/login", authHandler.Login)
			authGroup.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要认证的路由
		// 更新检查（公开）
		api.GET("/update/check", h.CheckUpdate)

		// 认证保护的路由
		authRequired := api.Group("")
		authRequired.Use(auth.JWTMiddleware(jwtService))

		// 应用管理（需要认证）
		{
			authRequired.GET("/apps", appHandler.List)
			authRequired.POST("/apps", appHandler.Create)
			authRequired.POST("/apps/:name", appHandler.Update)
			authRequired.DELETE("/apps/:name", appHandler.Delete)
		}

		// GitHub 风格的 releases 和 assets 路由
		{
			// 获取应用的版本列表
			authRequired.GET("/:name/releases", releaseHandler.ListByApp)
			// 获取应用的最新版本
			authRequired.GET("/:name/releases/latest", releaseHandler.GetLatest)
			// 下载指定版本的资源文件 (GitHub 风格，包含 platform 和 arch)
			authRequired.GET("/:name/releases/download/:tag/:platformArch/:filename", assetHandler.DownloadByTag)

			// 创建、更新、删除版本（需要认证）
			authRequired.POST("/:name/releases", releaseHandler.Create)
			authRequired.PUT("/:name/releases/:tag", releaseHandler.Update)
			authRequired.DELETE("/:name/releases/:tag", releaseHandler.Delete)

			// 创建资源文件（上传文件并创建记录，原子操作）
			authRequired.POST("/:name/releases/:tag/assets", assetHandler.Create)
		}
	}
}
