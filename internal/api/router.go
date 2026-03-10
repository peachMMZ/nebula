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
		authRequired := api.Group("")
		authRequired.Use(auth.JWTMiddleware(jwtService))
		{
			// 用户相关
			authRequired.GET("/auth/profile", authHandler.GetProfile)
			authRequired.POST("/auth/change-password", authHandler.ChangePassword)
		}

		// 更新检查（公开）
		api.GET("/check-update", h.CheckUpdate)

		// 应用管理（需要认证）
		appGroup := api.Group("/apps")
		appGroup.Use(auth.JWTMiddleware(jwtService))
		{
			appGroup.GET("", appHandler.List)
			appGroup.POST("", appHandler.Create)
			appGroup.GET("/:id", appHandler.Get)
			appGroup.PUT("/:id", appHandler.Update)
			appGroup.DELETE("/:id", appHandler.Delete)

			// 应用相关的版本
			appGroup.GET("/:id/releases", releaseHandler.ListByApp)
			appGroup.GET("/:id/releases/latest", releaseHandler.GetLatest)
		}

		// 版本管理（需要认证）
		releaseGroup := api.Group("/releases")
		releaseGroup.Use(auth.JWTMiddleware(jwtService))
		{
			releaseGroup.GET("", releaseHandler.List)
			releaseGroup.POST("", releaseHandler.Create)
			releaseGroup.GET("/:id", releaseHandler.Get)
			releaseGroup.PUT("/:id", releaseHandler.Update)
			releaseGroup.DELETE("/:id", releaseHandler.Delete)

			// 版本相关的资源
			releaseGroup.GET("/:id/assets", assetHandler.ListByRelease)
			releaseGroup.POST("/:id/assets/upload", assetHandler.Upload)
		}

		// 资源管理（需要认证）
		assetGroup := api.Group("/assets")
		assetGroup.Use(auth.JWTMiddleware(jwtService))
		{
			assetGroup.GET("", assetHandler.List)
			assetGroup.GET("/:id", assetHandler.Get)
			assetGroup.POST("", assetHandler.Create)
			assetGroup.PUT("/:id", assetHandler.Update)
			assetGroup.DELETE("/:id", assetHandler.Delete)
			assetGroup.GET("/:id/download", assetHandler.Download)
		}
	}
}
