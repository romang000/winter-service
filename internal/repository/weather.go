package repository

import (
	"database/sql"
	"github.com/romang000/winter-service/internal/models"
)

type WeatherRepository struct {
	db *sql.DB
}

func New(db *sql.DB) (*WeatherRepository, error) {
	repo := &WeatherRepository{db: db}
	
	if err := repo.createTable(); err != nil {
		return nil, err
	}
	return repo, nil
}

func (w *WeatherRepository) createTable() error {
	query := `
CREATE TABLE IF NOT EXISTS Weather (
    id BIGSERIAL PRIMARY KEY,
    temperature REAL
);`
	_, err := w.db.Exec(query)
	
	return err
}

func (w *WeatherRepository) Add(weather models.Weather) error {
	query := `INSERT INTO weather (temperature) VALUES ($1)`
	
	_, err := w.db.Exec(query, weather.Temperature)
	return err
}
