package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
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
