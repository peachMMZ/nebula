package release

import (
	"errors"
	"nebula/types"
	"time"

	"gorm.io/gorm"
)

type ReleaseService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *ReleaseService {
	return &ReleaseService{db: db}
}

// List 根据 app_name 获取版本列表 (必须提供 app_name)
func (s *ReleaseService) List(appName string, channel string) ([]Release, error) {
	if appName == "" {
		return nil, errors.New("app_name is required")
	}

	// 先通过 app_name 查找 app_id
	var app struct {
		ID string `gorm:"column:id"`
	}
	err := s.db.Table("apps").Select("id").Where("name = ?", appName).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("app not found")
		}
		return nil, err
	}

	// 查询该应用的 releases
	var releases []Release
	query := s.db.Where("app_id = ?", app.ID)
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	err = query.Order("created_at DESC").Find(&releases).Error
	return releases, err
}

// ListByApp 获取指定应用的所有版本 (通过 app_name)
func (s *ReleaseService) ListByApp(appName string) ([]Release, error) {
	return s.List(appName, "")
}

// ListByAppAndChannel 获取指定应用和渠道的版本 (通过 app_name)
func (s *ReleaseService) ListByAppAndChannel(appName, channel string) ([]Release, error) {
	return s.List(appName, channel)
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

// GetLatest 获取应用的最新版本 (通过 app_name)
func (s *ReleaseService) GetLatest(appName string, channel string) (*Release, error) {
	if appName == "" {
		return nil, errors.New("app_name is required")
	}

	// 先通过 app_name 查找 app_id
	var app struct {
		ID string `gorm:"column:id"`
	}
	err := s.db.Table("apps").Select("id").Where("name = ?", appName).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("app not found")
		}
		return nil, err
	}

	// 查询最新版本
	var release Release
	query := s.db.Where("app_id = ?", app.ID)
	if channel != "" {
		query = query.Where("channel = ?", channel)
	}
	err = query.Order("pub_date DESC").First(&release).Error
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

// CreateByAppName 通过应用名称创建新版本
func (s *ReleaseService) CreateByAppName(appName string, release Release) error {
	// 通过 app_name 查找 app_id
	var app struct {
		ID string `gorm:"column:id"`
	}
	err := s.db.Table("apps").Select("id").Where("name = ?", appName).First(&app).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("app not found")
		}
		return err
	}

	// 设置 app_id
	release.AppID = app.ID
	// 设置pub_date
	release.PubDate = types.JSONTime(time.Now())

	// 检查 tag 是否已存在
	err = s.db.Where("app_id = ? AND tag = ?", release.AppID, release.Tag).First(&Release{}).Error
	if err == nil {
		return errors.New("tag already exists for this app")
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

// UpdateByTag 通过应用名称和标签更新版本信息
func (s *ReleaseService) UpdateByTag(appName, tag string, release Release) error {
	// 通过 app_name 和 tag 查找 release
	var existingRelease Release
	err := s.db.Table("releases").
		Joins("JOIN apps ON releases.app_id = apps.id").
		Where("apps.name = ? AND releases.tag = ?", appName, tag).
		Select("releases.*").
		First(&existingRelease).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("release not found")
		}
		return err
	}

	// 准备更新数据
	data := map[string]any{
		"version":  release.Version,
		"notes":    release.Notes,
		"channel":  release.Channel,
		"pub_date": release.PubDate,
	}

	// 如果要更新 tag，检查是否重复
	if release.Tag != "" && release.Tag != existingRelease.Tag {
		err := s.db.Where("app_id = ? AND tag = ? AND id != ?", existingRelease.AppID, release.Tag, existingRelease.ID).
			First(&Release{}).Error
		if err == nil {
			return errors.New("tag already exists for this app")
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		data["tag"] = release.Tag
	}

	return s.db.Model(&Release{}).Where("id = ?", existingRelease.ID).Updates(data).Error
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

// DeleteByTag 通过应用名称和标签删除版本
func (s *ReleaseService) DeleteByTag(appName, tag string) error {
	// 通过 app_name 和 tag 查找 release
	var release Release
	err := s.db.Table("releases").
		Joins("JOIN apps ON releases.app_id = apps.id").
		Where("apps.name = ? AND releases.tag = ?", appName, tag).
		Select("releases.*").
		First(&release).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("release not found")
		}
		return err
	}

	// TODO: 同时删除相关的 Assets
	// 可以通过数据库外键级联删除，或者在这里手动删除

	return s.db.Delete(&Release{}, release.ID).Error
}
