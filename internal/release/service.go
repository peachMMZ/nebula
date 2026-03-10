package release

import (
	"errors"

	"gorm.io/gorm"
)

type ReleaseService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *ReleaseService {
	return &ReleaseService{db: db}
}

// List 获取所有版本列表
func (s *ReleaseService) List() ([]Release, error) {
	var releases []Release
	err := s.db.Order("created_at DESC").Find(&releases).Error
	return releases, err
}

// ListByApp 获取指定应用的所有版本
func (s *ReleaseService) ListByApp(appID string) ([]Release, error) {
	var releases []Release
	err := s.db.Where("app_id = ?", appID).Order("created_at DESC").Find(&releases).Error
	return releases, err
}

// ListByAppAndChannel 获取指定应用和渠道的版本
func (s *ReleaseService) ListByAppAndChannel(appID, channel string) ([]Release, error) {
	var releases []Release
	query := s.db.Where("app_id = ?", appID)
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	err := query.Order("created_at DESC").Find(&releases).Error
	return releases, err
}

// Get 获取单个版本详情
func (s *ReleaseService) Get(id uint) (*Release, error) {
	var release Release
	err := s.db.First(&release, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("release not found")
		}
		return nil, err
	}
	return &release, nil
}

// GetLatest 获取应用的最新版本
func (s *ReleaseService) GetLatest(appID string, channel string) (*Release, error) {
	var release Release
	query := s.db.Where("app_id = ?", appID)
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	err := query.Order("pub_date DESC").First(&release).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no release found for this app")
		}
		return nil, err
	}
	return &release, nil
}

// Create 创建新版本
func (s *ReleaseService) Create(release Release) error {
	// 检查应用是否存在
	var count int64
	err := s.db.Table("apps").Where("id = ?", release.AppID).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("app not found")
	}

	// 检查版本是否已存在
	err = s.db.Where("app_id = ? AND version = ?", release.AppID, release.Version).First(&Release{}).Error
	if err == nil {
		return errors.New("version already exists for this app")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.Create(&release).Error
}

// Update 更新版本信息
func (s *ReleaseService) Update(id uint, data map[string]any) error {
	// 检查版本是否存在
	var release Release
	err := s.db.First(&release, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("release not found")
		}
		return err
	}

	// 如果更新版本号，检查是否重复
	if newVersion, ok := data["version"].(string); ok {
		if newVersion != release.Version {
			err := s.db.Where("app_id = ? AND version = ? AND id != ?", release.AppID, newVersion, id).
				First(&Release{}).Error
			if err == nil {
				return errors.New("version already exists for this app")
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}
	}

	return s.db.Model(&Release{}).Where("id = ?", id).Updates(data).Error
}

// Delete 删除版本
func (s *ReleaseService) Delete(id uint) error {
	// 检查版本是否存在
	var release Release
	err := s.db.First(&release, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("release not found")
		}
		return err
	}

	// TODO: 同时删除相关的 Assets
	// 可以通过数据库外键级联删除，或者在这里手动删除

	return s.db.Delete(&Release{}, id).Error
}
