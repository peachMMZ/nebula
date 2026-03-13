package auth

// User 用户模型 (仅用于管理员虚拟用户)
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"` // admin
}
