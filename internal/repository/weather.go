package repository

import (
	"database/sql"
	"github.com/romang000/winter-service/internal/models"
)

type WeatherRepository struct {
	db *sql.DB
}

func New(db *sql.DB) (*WeatherRepository, error) {
	return &WeatherRepository{db: db}, nil
}

func (w *WeatherRepository) Add(weather models.Weather) error {
	query := `INSERT INTO weather (temperature) VALUES ($1)`
	
	_, err := w.db.Exec(query, weather.Temperature)
	return err
}
