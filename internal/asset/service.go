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
	// 检查 Release 是否存在
	var count int64
	err := s.db.Table("releases").Where("id = ?", releaseID).Count(&count).Error
	if err != nil {
		return nil, err
	}
	if count == 0 {
		return nil, errors.New("release not found")
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

	// 构建存储路径: releases/{releaseID}/{platform}-{arch}/{filename}
	storagePath := fmt.Sprintf("releases/%d/%s-%s/%s", releaseID, platform, arch, file.Filename)

	// 保存文件
	savedPath, err := s.storage.Save(storagePath, src)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	// 获取文件访问 URL
	url := s.storage.GetURL(savedPath)

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

// Create 创建资源记录（用于外部URL）
func (s *AssetService) Create(asset Asset) error {
	// 检查 Release 是否存在
	var count int64
	err := s.db.Table("releases").Where("id = ?", asset.ReleaseID).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("release not found")
	}

	// 检查是否已存在相同平台和架构的资源
	err = s.db.Where("release_id = ? AND platform = ? AND arch = ?", asset.ReleaseID, asset.Platform, asset.Arch).
		First(&Asset{}).Error
	if err == nil {
		return errors.New("asset already exists for this platform and architecture")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.Create(&asset).Error
}

// Update 更新资源信息（不包括文件）
func (s *AssetService) Update(id uint, data map[string]any) error {
	// 检查资源是否存在
	var asset Asset
	err := s.db.First(&asset, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("asset not found")
		}
		return err
	}

	// 如果更新平台或架构，检查是否会造成重复
	newPlatform, hasPlatform := data["platform"].(string)
	newArch, hasArch := data["arch"].(string)

	if hasPlatform || hasArch {
		checkPlatform := asset.Platform
		checkArch := asset.Arch
		if hasPlatform {
			checkPlatform = newPlatform
		}
		if hasArch {
			checkArch = newArch
		}

		err := s.db.Where("release_id = ? AND platform = ? AND arch = ? AND id != ?",
			asset.ReleaseID, checkPlatform, checkArch, id).
			First(&Asset{}).Error
		if err == nil {
			return errors.New("asset already exists for this platform and architecture")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return s.db.Model(&Asset{}).Where("id = ?", id).Updates(data).Error
}

// Delete 删除资源（包括文件）
func (s *AssetService) Delete(id uint) error {
	// 获取资源信息
	var asset Asset
	err := s.db.First(&asset, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("asset not found")
		}
		return err
	}

	// 删除存储的文件
	if asset.StoragePath != "" {
		if err := s.storage.Delete(asset.StoragePath); err != nil {
			// 记录错误但继续删除数据库记录
			// TODO: 添加日志
		}
	}

	// 删除数据库记录
	return s.db.Delete(&Asset{}, id).Error
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
