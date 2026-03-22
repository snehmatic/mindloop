// Package db handles database connection and migrations
package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/snehmatic/mindloop/internal/config"
	"github.com/snehmatic/mindloop/internal/log"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var logger = log.Get()

// Conn establishes a connection to a PostgreSQL database
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

// GetLocalDBPath returns the absolute path to the local SQLite database
func GetLocalDBPath() string {
	localFile := "mindloop_local.db"
	if utils.FileExists(localFile) {
		return localFile
	}
	return config.GetDataDir() + "/" + localFile
}

// LocalConn establishes a connection to the local SQLite database
func LocalConn() (*gorm.DB, error) {
	dbPath := GetLocalDBPath()
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
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

	logger.Info().Msgf("Connected to local SQLite DB at %s, migrations complete!", dbPath)
	return db, nil
}

// ConnectToDb connects to the database based on the provided application configuration
func ConnectToDb(appConfig config.Config) (*gorm.DB, error) {
	logger.Debug().Msg("Connecting to DB...")
	switch appConfig.Mode {
	case config.Local:
		return LocalConn()
	case config.ByoDB:
		fallthrough // as of now, ByoDB is same as API mode
	case config.API:
		connString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			appConfig.DBConfig.Host,
			appConfig.DBConfig.Port,
			appConfig.DBConfig.User,
			appConfig.DBConfig.Password,
			appConfig.DBConfig.Name,
		)
		return Conn(connString)
	default:
		return nil, fmt.Errorf("mode selected is invalid")
	}
}

// LocalDBFileExists checks if the local SQLite database file exists
func LocalDBFileExists() bool {
	return utils.FileExists(GetLocalDBPath())
}

// MigrateDB performs automatic database schema migrations
func MigrateDB(db *gorm.DB) error {
	err := db.AutoMigrate(
		&models.Intent{},
		&models.FocusSession{},
		&models.Habit{},
		&models.HabitLog{},
		&models.JournalEntry{},
		&models.Note{},
		&models.SideQuest{},
		&models.PointTransaction{},
		&models.Routine{},
		&models.Task{},
		&models.SubTask{},
	)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to migrate DB")
		return err
	}
	return nil
}
