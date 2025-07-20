package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
	"github.com/snehmatic/mindloop/config"
	"github.com/snehmatic/mindloop/internal/utils"
	"github.com/snehmatic/mindloop/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Conn(connString string) (*gorm.DB, error) {
	log.Debug().Msg("Function Conn called")
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
	})
	if err != nil {
		return db, err
	}

	err = db.AutoMigrate() // db model structs go here
	if err != nil {
		return db, err
	}

	log.Info().Msg("Connected to DB, migrations complete!")
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

	err = db.AutoMigrate(models.Habit{}) // db model structs go here
	if err != nil {
		return db, err
	}

	fmt.Println("Connected to local SQLite DB, migrations complete!")
	return db, nil
}

func ConnectToDb(appConfig config.Config) (*gorm.DB, error) {
	switch appConfig.Mode {
	case models.Local:
		return LocalConn()
	case models.Api:
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
