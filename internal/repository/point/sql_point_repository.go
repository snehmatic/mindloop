package point

import (
	"time"

	"github.com/snehmatic/mindloop/models"
	"gorm.io/gorm"
)

// sqlRepository implements Repository using GORM
type sqlRepository struct {
	DB *gorm.DB
}

// NewSQLRepository creates a new SQL-based point repository
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{DB: db}
}

// Create inserts a new point transaction record
func (r *sqlRepository) Create(tx models.PointTransaction) error {
	return r.DB.Create(&tx).Error
}

// GetTotalPoints returns the lifetime total points for the user
func (r *sqlRepository) GetTotalPoints() (int, error) {
	var total int
	err := r.DB.Model(&models.PointTransaction{}).Select("COALESCE(SUM(Points), 0)").Scan(&total).Error
	if err != nil {
		return 0, err
	}
	return total, nil
}

// GetPointSeries returns daily point totals for the given range
func (r *sqlRepository) GetPointSeries(start, end time.Time) ([]int, error) {
	days := int(end.Sub(start).Hours()/24) + 1
	if days < 1 {
		days = 1
	}

	stats := make([]int, days)

	var transactions []models.PointTransaction
	if err := r.DB.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end).Find(&transactions).Error; err != nil {
		return nil, err
	}

	for _, tx := range transactions {
		txDate := tx.CreatedAt.Truncate(24 * time.Hour)
		startDate := start.Truncate(24 * time.Hour)
		diff := int(txDate.Sub(startDate).Hours() / 24)

		if diff >= 0 && diff < days {
			stats[diff] += tx.Points
		}
	}

	return stats, nil
}

// GetPointsInRange returns the point transactions within a date range
func (r *sqlRepository) GetPointsInRange(start, end string) ([]models.PointTransaction, error) {
	var history []models.PointTransaction
	query := r.DB.Model(&models.PointTransaction{})
	if start != "" && end != "" {
		query = query.Where("CreatedAt >= ? AND CreatedAt <= ?", start, end)
	}
	err := query.Order("CreatedAt ASC").Find(&history).Error
	return history, err
}
