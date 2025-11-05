package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/romang000/winter-service/internal/config"
)

func NewDatabaseConnection(cfg config.DatabaseConfig) (*sql.DB, string, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.Username, cfg.Password, cfg.Name,
	)
	
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, "", fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return db, dsn, nil
}
