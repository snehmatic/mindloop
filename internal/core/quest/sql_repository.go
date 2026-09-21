package quest

import (
	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) Repository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *SQLRepository) CreateQuest(quest *models.SideQuest) error {
	return r.db.Create(quest).Error
}

func (r *SQLRepository) ListQuests() ([]models.SideQuest, error) {
	var quests []models.SideQuest
	result := r.db.Order("CreatedAt DESC").Find(&quests)
	return quests, result.Error
}

func (r *SQLRepository) GetActiveQuest() (*models.SideQuest, error) {
	var activeQuests []models.SideQuest
	if err := r.db.Where("status = ?", models.StatusActive).Limit(1).Find(&activeQuests).Error; err != nil {
		return nil, err
	}
	if len(activeQuests) == 0 {
		return nil, nil
	}
	return &activeQuests[0], nil
}

func (r *SQLRepository) GetQuest(id uint) (*models.SideQuest, error) {
	var quest models.SideQuest
	if err := r.db.First(&quest, id).Error; err != nil {
		return nil, err
	}
	return &quest, nil
}

func (r *SQLRepository) UpdateQuest(quest *models.SideQuest) error {
	return r.db.Save(quest).Error
}

func (r *SQLRepository) DeleteQuest(id uint) error {
	return r.db.Delete(&models.SideQuest{}, id).Error
}
