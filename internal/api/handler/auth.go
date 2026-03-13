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

// Login 管理员登录
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
