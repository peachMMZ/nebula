package app

import (
	"nebula/pkg/util"

	"gorm.io/gorm"
)

type AppService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *AppService {
	return &AppService{db: db}
}

func (s *AppService) List(params map[string]any) ([]App, error) {
	var apps []App
	query := s.db

	// 只处理name和description参数，使用模糊查询
	if name, ok := params["name"]; ok {
		query = query.Where("name LIKE ?", "%"+name.(string)+"%")
	}
	if description, ok := params["description"]; ok {
		query = query.Where("description LIKE ?", "%"+description.(string)+"%")
	}

	err := query.Order("updated_at desc").Find(&apps).Error
	return apps, err
}

func (s *AppService) Get(id string) (*App, error) {
	var app App
	err := s.db.First(&app, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (s *AppService) Create(app App) error {
	// 生成唯一的短ID作为应用ID
	var appExists App
	var err error
	var id string

	for {
		// 生成短ID
		id = util.GenerateShortID()

		// 检查ID是否已存在
		err = s.db.First(&appExists, "id = ?", id).Error
		if err != nil {
			// 如果错误是记录未找到，说明ID可用
			if err == gorm.ErrRecordNotFound {
				break
			}
			// 其他错误直接返回
			return err
		}
		// 如果ID已存在，继续循环生成新ID
	}

	// 使用生成的唯一ID
	app.ID = id
	return s.db.Create(&app).Error
}

func (s *AppService) Update(id string, data map[string]any) error {
	return s.db.Model(&App{}).Where("id = ?", id).Updates(data).Error
}

func (s *AppService) Delete(id string) error {
	return s.db.Delete(&App{}, "id = ?", id).Error
}
