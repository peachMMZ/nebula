package storage

import (
	"io"
)

// Storage 定义文件存储接口，支持本地存储、OSS、S3 等
type Storage interface {
	// Save 保存文件，返回存储路径
	Save(filename string, content io.Reader) (string, error)

	// Delete 删除文件
	Delete(path string) error

	// GetURL 获取文件的访问 URL
	GetURL(path string) string

	// Exists 检查文件是否存在
	Exists(path string) bool
}
