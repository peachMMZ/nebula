package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage 本地文件系统存储实现
type LocalStorage struct {
	basePath string // 文件存储的根目录
	baseURL  string // 文件访问的 URL 前缀
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(basePath, baseURL string) (*LocalStorage, error) {
	// 确保存储目录存在
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &LocalStorage{
		basePath: basePath,
		baseURL:  baseURL,
	}, nil
}

// Save 保存文件到本地
func (s *LocalStorage) Save(filename string, content io.Reader) (string, error) {
	// 构建完整的文件路径
	fullPath := filepath.Join(s.basePath, filename)

	// 确保目标目录存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}

	// 创建文件
	file, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// 复制内容
	if _, err := io.Copy(file, content); err != nil {
		// 如果保存失败，尝试删除文件
		os.Remove(fullPath)
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// 返回相对路径（存储在数据库中）
	return filename, nil
}

// Delete 删除本地文件
func (s *LocalStorage) Delete(path string) error {
	fullPath := filepath.Join(s.basePath, path)
	if err := os.Remove(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

// GetURL 获取文件的访问 URL
func (s *LocalStorage) GetURL(path string) string {
	if path == "" {
		return ""
	}
	// 将路径中的反斜杠替换为正斜杠（Windows 兼容）
	path = filepath.ToSlash(path)
	return s.baseURL + "/" + path
}

// Exists 检查文件是否存在
func (s *LocalStorage) Exists(path string) bool {
	fullPath := filepath.Join(s.basePath, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

// GetFullPath 获取文件的完整本地路径
func (s *LocalStorage) GetFullPath(path string) string {
	return filepath.Join(s.basePath, path)
}
