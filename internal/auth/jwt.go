package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 声明
type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair 访问令牌和刷新令牌
type TokenPair struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    int64  `json:"expiresIn"` // 秒
}

// JWTService JWT 服务
type JWTService struct {
	secret               string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

// NewJWTService 创建 JWT 服务
func NewJWTService(secret string, accessTokenDuration, refreshTokenDuration time.Duration) *JWTService {
	return &JWTService{
		secret:               secret,
		accessTokenDuration:  accessTokenDuration,
		refreshTokenDuration: refreshTokenDuration,
	}
}

// GenerateTokenPair 生成访问令牌和刷新令牌
func (s *JWTService) GenerateTokenPair(userID, username, role string) (*TokenPair, error) {
	// 生成访问令牌
	accessToken, err := s.generateToken(userID, username, role, s.accessTokenDuration)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌（更长的过期时间，不包含敏感信息）
	refreshToken, err := s.generateToken(userID, "", "", s.refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.accessTokenDuration.Seconds()),
	}, nil
}

// generateToken 生成 token
func (s *JWTService) generateToken(userID, username, role string, duration time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.secret))
}

// ParseToken 解析 token
func (s *JWTService) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// ValidateToken 验证 token 是否有效
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	claims, err := s.ParseToken(tokenString)
	if err != nil {
		return nil, err
	}

	// 检查是否过期
	if time.Now().Unix() > claims.ExpiresAt.Unix() {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

// RefreshAccessToken 使用刷新令牌生成新的访问令牌
func (s *JWTService) RefreshAccessToken(refreshToken string, user *User) (*TokenPair, error) {
	// 验证刷新令牌
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// 确保刷新令牌属于该用户
	if claims.UserID != user.ID {
		return nil, errors.New("token does not match user")
	}

	// 生成新的 token 对
	return s.GenerateTokenPair(user.ID, user.Username, user.Role)
}
