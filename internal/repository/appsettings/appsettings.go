package appsettings

import "github.com/snehmatic/mindloop/models"

// Repository provides data-access operations for AppSetting records.
type Repository interface {
	GetSetting(key string) (*models.AppSetting, error)
	SaveSetting(setting *models.AppSetting) error
}
