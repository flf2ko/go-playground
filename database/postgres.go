package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/flf2ko/playground/go-api-sample/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	conn *gorm.DB
}

func NewDB() (*DB, error) {
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "jsonapi")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Taipei",
		host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto migrate the schema
	if err := db.AutoMigrate(&models.JSONRecord{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL with GORM")
	return &DB{conn: db}, nil
}

func (db *DB) SaveJSONRecord(ctx context.Context, url, content string) (*models.JSONRecord, error) {
	record := &models.JSONRecord{
		URL:     url,
		Content: content,
	}

	if err := db.conn.WithContext(ctx).Create(record).Error; err != nil {
		return nil, fmt.Errorf("failed to save JSON record: %w", err)
	}

	return record, nil
}

func (db *DB) GetJSONRecords(ctx context.Context, limit int) ([]models.JSONRecord, error) {
	if limit <= 0 {
		limit = 10
	}

	var records []models.JSONRecord
	if err := db.conn.WithContext(ctx).Order("created_at DESC").Limit(limit).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("failed to get JSON records: %w", err)
	}

	return records, nil
}

func (db *DB) Close() error {
	if db.conn != nil {
		sqlDB, err := db.conn.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
