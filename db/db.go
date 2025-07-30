package db

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var logger = log.Get()

func Conn(connString string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		return db, err
	}

	err = MigrateDB(db)
	if err != nil {
		return db, err
	}

	logger.Info().Msg("Connected to DB, migrations complete!")
	return db, nil
}

func LocalConn() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("mindloop_local.db"), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		return db, err
	}

	err = MigrateDB(db)
	if err != nil {
		return db, err
	}

	logger.Info().Msg("Connected to local SQLite DB, migrations complete!")
	return db, nil
}

func ConnectToDb(appConfig config.Config) (*gorm.DB, error) {
	logger.Debug().Msg("Connecting to DB...")
	switch appConfig.Mode {
	case config.Local:
		return LocalConn()
	case config.ByoDB:
		fallthrough // as of now, ByoDB is same as Api mode
	case config.Api:
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			appConfig.DBConfig.Host,
			appConfig.DBConfig.Port,
			appConfig.DBConfig.User,
			appConfig.DBConfig.Password,
			appConfig.DBConfig.Name,
		)
		return Conn(connString)
	default:
		return nil, fmt.Errorf("Mode selected is invalid!")
	}
}

func LocalDBFileExists() bool {
	return utils.FileExists("mindloop_local.db")
}

func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Intent{},
		&models.FocusSession{},
		&models.Habit{},
		&models.HabitLog{},
		&models.JournalEntry{},
	)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to migrate DB")
		return err
	}
	return nil
}
