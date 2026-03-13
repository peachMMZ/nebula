package auth

import (
	"errors"
)

type AuthService struct {
	jwtService    *JWTService
	adminUsername string
	adminPassword string
}

func NewAuthService(jwtService *JWTService, adminUsername, adminPassword string) *AuthService {
	return &AuthService{
		jwtService:    jwtService,
		adminUsername: adminUsername,
		adminPassword: adminPassword,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Login 管理员登录
func (s *AuthService) Login(req LoginRequest) (*User, *TokenPair, error) {
	// 验证管理员账号
	if req.Username != s.adminUsername || req.Password != s.adminPassword {
		return nil, nil, errors.New("invalid username or password")
	}

	// 创建管理员用户对象
	adminUser := &User{
		ID:       "admin",
		Username: s.adminUsername,
		Email:    s.adminUsername + "@admin.local",
		Role:     "admin",
	}

	// 生成 token
	tokens, err := s.jwtService.GenerateTokenPair(adminUser.ID, adminUser.Username, adminUser.Role)
	if err != nil {
		return nil, nil, err
	}

	return adminUser, tokens, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(req RefreshTokenRequest) (*TokenPair, error) {
	// 解析刷新令牌获取用户信息
	claims, err := s.jwtService.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 只支持管理员用户
	if claims.UserID != "admin" {
		return nil, errors.New("user not found")
	}

	// 创建管理员用户对象
	adminUser := &User{
		ID:       "admin",
		Username: s.adminUsername,
		Email:    s.adminUsername + "@admin.local",
		Role:     "admin",
	}

	// 使用刷新令牌生成新的 token 对
	tokens, err := s.jwtService.RefreshAccessToken(req.RefreshToken, adminUser)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}
