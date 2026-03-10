package app

import "gorm.io/gorm"

type AppService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *AppService {
	return &AppService{db: db}
}

func (s *AppService) List() ([]App, error) {
	var apps []App
	err := s.db.Find(&apps).Error
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
	return s.db.Create(&app).Error
}

func (s *AppService) Update(id string, data map[string]any) error {
	return s.db.Model(&App{}).Where("id = ?", id).Updates(data).Error
}

func (s *AppService) Delete(id string) error {
	return s.db.Delete(&App{}, "id = ?", id).Error
}
