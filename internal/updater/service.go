package updater

import (
	"errors"
	"nebula/internal/asset"
	"nebula/internal/release"
	"nebula/pkg/util"

	"gorm.io/gorm"
)

type CheckRequest struct {
	App      string
	Version  string
	Platform string
	Arch     string
}

type CheckResponse struct {
	Update   bool   `json:"update"`
	Version  string `json:"version,omitempty"`
	Notes    string `json:"notes,omitempty"`
	URL      string `json:"url,omitempty"`
	Checksum string `json:"checksum,omitempty"`
}

func CheckUpdate(db *gorm.DB, req CheckRequest) (*CheckResponse, error) {
	// 参数验证
	if req.App == "" {
		return nil, errors.New("app is required")
	}
	if req.Version == "" {
		return nil, errors.New("version is required")
	}
	if req.Platform == "" {
		return nil, errors.New("platform is required")
	}
	if req.Arch == "" {
		return nil, errors.New("arch is required")
	}

	// 查找最新版本（按发布日期排序）
	var latest release.Release
	err := db.Where("app_id = ?", req.App).
		Order("pub_date DESC").
		First(&latest).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no release found for this app")
		}
		return nil, err
	}

	// 比较版本号
	if !util.IsNewerVersion(req.Version, latest.Version) {
		return &CheckResponse{Update: false}, nil
	}

	// 查找对应平台和架构的资源
	var ast asset.Asset
	err = db.Where("release_id = ? AND platform = ? AND arch = ?",
		latest.ID, req.Platform, req.Arch).First(&ast).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no asset found for this platform and architecture")
		}
		return nil, err
	}

	// 验证资源 URL 是否存在
	if ast.URL == "" {
		return nil, errors.New("asset URL is empty")
	}

	return &CheckResponse{
		Update:   true,
		Version:  latest.Version,
		Notes:    latest.Notes,
		URL:      ast.URL,
		Checksum: ast.Checksum,
	}, nil
}
