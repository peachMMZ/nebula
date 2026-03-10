package handler

import (
	"nebula/internal/api/response"
	"nebula/internal/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *auth.AuthService
}

func NewAuthHandler(service *auth.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register 用户注册
// POST /api/auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req auth.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	user, tokens, err := h.service.Register(req)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	response.Ok(c, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// Login 用户登录
// POST /api/auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	user, tokens, err := h.service.Login(req)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}

	response.Ok(c, gin.H{
		"user":   user,
		"tokens": tokens,
	})
}

// RefreshToken 刷新访问令牌
// POST /api/auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	tokens, err := h.service.RefreshToken(req)
	if err != nil {
		response.Fail(c, 401, err.Error())
		return
	}

	response.Ok(c, tokens)
}

// GetProfile 获取当前用户信息
// GET /api/auth/profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	response.Ok(c, user)
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

// ChangePassword 修改密码
// POST /api/auth/change-password
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := auth.GetCurrentUserID(c)
	if !exists {
		response.Fail(c, 401, "unauthorized")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailBadRequest(c, err.Error())
		return
	}

	err := h.service.ChangePassword(userID, req.OldPassword, req.NewPassword)
	if err != nil {
		response.FailServer(c, err.Error())
		return
	}

	response.OkMsg(c, "password changed successfully")
}
