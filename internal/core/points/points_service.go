package points

import (
	"time"

	"github.com/snehmatic/mindloop/internal/repository/point"
	"github.com/snehmatic/mindloop/models"
)

// Service awards points and exposes query methods.
// Implemented by pointsService (SQL-backed).
type Service interface {
	AwardPoints(activityType models.PointCategory, activityID uint, points int) (bool, error)
	GetTotalPoints() (int, error)
	GetPointsInRange(start, end string) ([]models.PointTransaction, error)
	GetPointSeries(start, end time.Time) ([]int, error)
}

var _ Service = (*pointsService)(nil)

type pointsService struct {
	repo point.Repository
}

func NewService(repo point.Repository) *pointsService {
	return &pointsService{repo: repo}
}

func (s *pointsService) AwardPoints(activityType models.PointCategory, activityID uint, points int) (bool, error) {
	currentTotal, err := s.repo.GetTotalPoints()
	if err != nil {
		return false, err
	}

	transaction := models.PointTransaction{
		ActivityType: activityType,
		ActivityID:   activityID,
		Points:       points,
	}
	if err := s.repo.Create(transaction); err != nil {
		return false, err
	}

	newTotal := currentTotal + points

	currentMilestone := currentTotal / MilestoneInterval
	newMilestone := newTotal / MilestoneInterval

	return newMilestone > currentMilestone, nil
}

func (s *pointsService) GetTotalPoints() (int, error) {
	return s.repo.GetTotalPoints()
}

func (s *pointsService) GetPointsInRange(start, end string) ([]models.PointTransaction, error) {
	return s.repo.GetPointsInRange(start, end)
}

func (s *pointsService) GetPointSeries(start, end time.Time) ([]int, error) {
	return s.repo.GetPointSeries(start, end)
}
