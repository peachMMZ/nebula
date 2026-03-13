package asset

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"nebula/internal/storage"

	"gorm.io/gorm"
)

type AssetService struct {
	db      *gorm.DB
	storage storage.Storage
}

func NewService(db *gorm.DB, storage storage.Storage) *AssetService {
	return &AssetService{
		db:      db,
		storage: storage,
	}
}

// List 获取所有资源列表
func (s *AssetService) List() ([]Asset, error) {
	var assets []Asset
	err := s.db.Order("created_at DESC").Find(&assets).Error
	return assets, err
}

// ListByRelease 获取指定发布版本的所有资源
func (s *AssetService) ListByRelease(releaseID uint) ([]Asset, error) {
	var assets []Asset
	err := s.db.Where("release_id = ?", releaseID).Find(&assets).Error
	return assets, err
}

// ListByAppAndVersion 根据应用名称和版本标签获取资源列表 (GitHub API 风格)
func (s *AssetService) ListByAppAndVersion(appName, version string) ([]Asset, error) {
	var assets []Asset
	// 联表查询: assets -> releases -> apps
	err := s.db.Table("assets").
		Joins("JOIN releases ON assets.release_id = releases.id").
		Joins("JOIN apps ON releases.app_id = apps.id").
		Where("apps.name = ? AND releases.tag = ?", appName, version).
		Select("assets.*").
		Find(&assets).Error
	return assets, err
}

// Get 获取单个资源详情
func (s *AssetService) Get(id uint) (*Asset, error) {
	var asset Asset
	err := s.db.First(&asset, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}
	return &asset, nil
}

// GetByReleaseAndPlatform 根据发布版本、平台和架构获取资源
func (s *AssetService) GetByReleaseAndPlatform(releaseID uint, platform, arch string) (*Asset, error) {
	var asset Asset
	err := s.db.Where("release_id = ? AND platform = ? AND arch = ?", releaseID, platform, arch).
		First(&asset).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found for this platform and architecture")
		}
		return nil, err
	}
	return &asset, nil
}

// Upload 上传文件并创建资源记录
func (s *AssetService) Upload(releaseID uint, platform, arch string, file *multipart.FileHeader) (*Asset, error) {
	// 检查 Release 是否存在并获取 tag 和 app_name
	var release struct {
		ID      uint   `gorm:"column:id"`
		Tag     string `gorm:"column:tag"`
		AppID   string `gorm:"column:app_id"`
		AppName string `gorm:"column:name"`
	}
	err := s.db.Table("releases").
		Joins("JOIN apps ON releases.app_id = apps.id").
		Where("releases.id = ?", releaseID).
		Select("releases.id, releases.tag, releases.app_id, apps.name").
		First(&release).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("release not found")
		}
		return nil, err
	}

	// 检查是否已存在相同平台和架构的资源
	var existing Asset
	err = s.db.Where("release_id = ? AND platform = ? AND arch = ?", releaseID, platform, arch).
		First(&existing).Error
	if err == nil {
		return nil, errors.New("asset already exists for this platform and architecture")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// 打开上传的文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// 计算文件校验和
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, fmt.Errorf("failed to calculate checksum: %w", err)
	}
	checksum := hex.EncodeToString(hash.Sum(nil))

	// 重新打开文件用于保存（因为已经读取过了）
	src.Close()
	src, err = file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to reopen file: %w", err)
	}
	defer src.Close()

	// 构建存储路径: releases/{tag}/{platform}-{arch}/{filename}
	storagePath := fmt.Sprintf("releases/%s/%s-%s/%s", release.Tag, platform, arch, file.Filename)

	// 保存文件
	savedPath, err := s.storage.Save(storagePath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 构建 GitHub 风格的下载 URL: /api/:name/releases/download/:tag/:platform-:arch/:filename
	url := fmt.Sprintf("/api/%s/releases/download/%s/%s-%s/%s", release.AppName, release.Tag, platform, arch, file.Filename)

	// 创建数据库记录
	asset := Asset{
		ReleaseID:   releaseID,
		Platform:    platform,
		Arch:        arch,
		URL:         url,
		StoragePath: savedPath,
		Checksum:    checksum,
		Signature:   "", // TODO: 实现文件签名
	}

	if err := s.db.Create(&asset).Error; err != nil {
		// 如果数据库操作失败，删除已上传的文件
		s.storage.Delete(savedPath)
		return nil, err
	}

	return &asset, nil
}

// CreateByTag 通过应用名称和标签创建资源（上传文件并创建记录，原子操作）
func (s *AssetService) CreateByTag(appName, tag, platform, arch string, file *multipart.FileHeader) (*Asset, error) {
	// 通过 app_name 和 tag 查找 release_id
	var release struct {
		ID uint `gorm:"column:id"`
	}
	err := s.db.Table("releases").
		Joins("JOIN apps ON releases.app_id = apps.id").
		Where("apps.name = ? AND releases.tag = ?", appName, tag).
		Select("releases.id").
		First(&release).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("release not found")
		}
		return nil, err
	}

	// 调用 Upload 方法完成文件上传和记录创建
	return s.Upload(release.ID, platform, arch, file)
}

// GetStoragePath 获取资源的存储路径
func (s *AssetService) GetStoragePath(id uint) (string, error) {
	var asset Asset
	err := s.db.First(&asset, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("asset not found")
		}
		return "", err
	}

	if asset.StoragePath == "" {
		return "", errors.New("storage path not available")
	}

	// 如果是本地存储，返回完整路径
	if localStorage, ok := s.storage.(*storage.LocalStorage); ok {
		return localStorage.GetFullPath(asset.StoragePath), nil
	}

	return asset.StoragePath, nil
}

// GetStoragePathByTag 通过 app_name、tag、platformArch 和 filename 获取资源的存储路径 (GitHub API 风格)
func (s *AssetService) GetStoragePathByTag(appName, tag, platformArch, filename string) (string, error) {
	// 从 platformArch 中解析 platform 和 arch (格式: windows-amd64)
	// 但实际上我们可以直接通过 storagePath 匹配,因为存储路径已经包含了这些信息
	// 构建预期的存储路径: releases/{tag}/{platformArch}/{filename}
	expectedPath := fmt.Sprintf("releases/%s/%s/%s", tag, platformArch, filename)

	var asset Asset
	// 通过存储路径查找资源
	err := s.db.Where("storage_path = ?", expectedPath).First(&asset).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("asset not found")
		}
		return "", err
	}

	if asset.StoragePath == "" {
		return "", errors.New("storage path not available")
	}

	// 如果是本地存储，返回完整路径
	if localStorage, ok := s.storage.(*storage.LocalStorage); ok {
		return localStorage.GetFullPath(asset.StoragePath), nil
	}

	return asset.StoragePath, nil
}
