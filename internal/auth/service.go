package auth

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthService struct {
	db            *gorm.DB
	jwtService    *JWTService
	adminUsername string
	adminPassword string
}

func NewAuthService(db *gorm.DB, jwtService *JWTService, adminUsername, adminPassword string) *AuthService {
	return &AuthService{
		db:            db,
		jwtService:    jwtService,
		adminUsername: adminUsername,
		adminPassword: adminPassword,
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
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

// Register 用户注册
func (s *AuthService) Register(req RegisterRequest) (*User, *TokenPair, error) {
	// 检查用户名是否已存在
	var existingUser User
	err := s.db.Where("username = ?", req.Username).First(&existingUser).Error
	if err == nil {
		return nil, nil, errors.New("username already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	// 检查邮箱是否已存在
	err = s.db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, nil, errors.New("email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, err
	}

	// 加密密码
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, nil, err
	}

	// 创建用户
	user := User{
		ID:       uuid.New().String(),
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user", // 默认角色
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, nil, err
	}

	// 生成 token
	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, nil, err
	}

	return &user, tokens, nil
}

// Login 用户登录
func (s *AuthService) Login(req LoginRequest) (*User, *TokenPair, error) {
	// 1. 先检查是否为配置的管理员账号
	if req.Username == s.adminUsername && req.Password == s.adminPassword {
		// 管理员登录成功，创建虚拟用户对象
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

	// 2. 如果不是管理员，查找数据库中的用户
	var user User
	err := s.db.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid username or password")
		}
		return nil, nil, err
	}

	// 验证密码
	if !user.CheckPassword(req.Password) {
		return nil, nil, errors.New("invalid username or password")
	}

	// 生成 token
	tokens, err := s.jwtService.GenerateTokenPair(user.ID, user.Username, user.Role)
	if err != nil {
		return nil, nil, err
	}

	return &user, tokens, nil
}

// RefreshToken 刷新访问令牌
func (s *AuthService) RefreshToken(req RefreshTokenRequest) (*TokenPair, error) {
	// 解析刷新令牌获取用户 ID
	claims, err := s.jwtService.ParseToken(req.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 获取用户信息
	var user User
	err = s.db.First(&user, "id = ?", claims.UserID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	// 使用刷新令牌生成新的 token 对
	tokens, err := s.jwtService.RefreshAccessToken(req.RefreshToken, &user)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetUserByID 根据 ID 获取用户
func (s *AuthService) GetUserByID(userID string) (*User, error) {
	// 特殊处理管理员
	if userID == "admin" {
		return &User{
			ID:       "admin",
			Username: s.adminUsername,
			Email:    s.adminUsername + "@admin.local",
			Role:     "admin",
		}, nil
	}

	var user User
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID, oldPassword, newPassword string) error {
	var user User
	err := s.db.First(&user, "id = ?", userID).Error
	if err != nil {
		return errors.New("user not found")
	}

	// 验证旧密码
	if !user.CheckPassword(oldPassword) {
		return errors.New("invalid old password")
	}

	// 加密新密码
	hashedPassword, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 更新密码
	return s.db.Model(&user).Update("password", hashedPassword).Error
}
