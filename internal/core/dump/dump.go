package dump

import (
	"fmt"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateDump(content string) (*models.BrainDump, error) {
	if content == "" {
		return nil, fmt.Errorf("content cannot be empty")
	}
	bd := &models.BrainDump{
		Content: content,
	}
	if err := s.db.Create(bd).Error; err != nil {
		return nil, err
	}
	return bd, nil
}
