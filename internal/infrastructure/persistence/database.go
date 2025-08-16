package persistence

import (
	"fmt"

	"github.com/snehmatic/mindloop/internal/domain/entities"
	"github.com/snehmatic/mindloop/internal/infrastructure/config"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	FilePath string // for SQLite
}

// Database manages database connections and operations
type Database struct {
	db     *gorm.DB
	config DatabaseConfig
}

// NewDatabase creates a new database instance
func NewDatabase(config DatabaseConfig) (*Database, error) {
	db, err := connectDB(config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &Database{
		db:     db,
		config: config,
	}, nil
}

// connectDB establishes a database connection based on the configuration
func connectDB(config DatabaseConfig) (*gorm.DB, error) {
	gormConfig := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   true,
		},
		Logger: logger.Default.LogMode(logger.Silent), // Reduce log noise
	}

	var db *gorm.DB
	var err error

	switch config.Driver {
	case "postgres":
		dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.Host, config.Port, config.User, config.Password, config.Name)
		db, err = gorm.Open(postgres.Open(dsn), gormConfig)
	case "sqlite":
		db, err = gorm.Open(sqlite.Open(config.FilePath), gormConfig)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", config.Driver)
	}

	if err != nil {
		return nil, err
	}

	return db, nil
}

// GetDB returns the database instance
func (d *Database) GetDB() *gorm.DB {
	return d.db
}

// Migrate runs database migrations
func (d *Database) Migrate() error {
	return d.db.AutoMigrate(
		&entities.Habit{},
		&entities.HabitLog{},
		&entities.Intent{},
		&entities.FocusSession{},
		&entities.JournalEntry{},
	)
}

// Close closes the database connection
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Health checks database health
func (d *Database) Health() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}

// CreateDatabaseFromConfig creates a database instance from app config
func CreateDatabaseFromConfig(appConfig *config.Config) (*Database, error) {
	var dbConfig DatabaseConfig

	switch appConfig.Mode {
	case config.Local:
		dbConfig = DatabaseConfig{
			Driver:   "sqlite",
			FilePath: "mindloop_local.db",
		}
	case config.Api, config.ByoDB:
		dbConfig = DatabaseConfig{
			Driver:   "postgres",
			Host:     appConfig.DBConfig.Host,
			Port:     appConfig.DBConfig.Port,
			User:     appConfig.DBConfig.User,
			Password: appConfig.DBConfig.Password,
			Name:     appConfig.DBConfig.Name,
		}
	default:
		return nil, fmt.Errorf("unsupported mode: %s", appConfig.Mode)
	}

	db, err := NewDatabase(dbConfig)
	if err != nil {
		return nil, err
	}

	// Run migrations
	if err := db.Migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}
