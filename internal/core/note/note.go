package note

import (
	"errors"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{DB: db}
}

func (s *Service) CreateNote(title, content, labels string) (*models.Note, error) {
	if title == "" && content == "" {
		return nil, errors.New("note must have a title or content")
	}
	note := &models.Note{
		Title:   title,
		Content: content,
		Labels:  labels,
	}
	if err := s.DB.Create(note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

func (s *Service) ListNotes() ([]models.Note, error) {
	var notes []models.Note
	if err := s.DB.Order("UpdatedAt desc").Find(&notes).Error; err != nil {
		return nil, err
	}
	return notes, nil
}

func (s *Service) GetNote(id int) (*models.Note, error) {
	var note models.Note
	if err := s.DB.First(&note, id).Error; err != nil {
		return nil, err
	}
	return &note, nil
}

func (s *Service) UpdateNote(id int, title, content, labels string) (*models.Note, error) {
	note, err := s.GetNote(id)
	if err != nil {
		return nil, err
	}
	note.Title = title
	note.Content = content
	note.Labels = labels
	if err := s.DB.Save(note).Error; err != nil {
		return nil, err
	}
	return note, nil
}

func (s *Service) DeleteNote(id int) error {
	return s.DB.Delete(&models.Note{}, id).Error
}

func (s *Service) DeleteAll() error {
	return s.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Note{}).Error
}
