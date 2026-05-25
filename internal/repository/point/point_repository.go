package point

import (
	"time"

	"github.com/snehmatic/mindloop/models"
)

// Repository defines the interface for point transaction data access
type Repository interface {
	Create(tx models.PointTransaction) error
	GetTotalPoints() (int, error)
	GetPointSeries(start, end time.Time) ([]int, error)
	GetPointsInRange(start, end string) ([]models.PointTransaction, error)
}
