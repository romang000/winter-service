package service

import (
	"github.com/romang000/winter-service/internal/models"
	"log/slog"
)

type WeatherRepository interface {
	Add(weather models.Weather) error
}

type WeatherService struct {
	weatherRepository WeatherRepository
	logger            *slog.Logger
}

func NewWeatherService(repository WeatherRepository, logger *slog.Logger) *WeatherService {
	return &WeatherService{
		weatherRepository: repository,
		logger:            logger,
	}
}

func (w *WeatherService) AddWeather(weather models.Weather) error {
	return w.weatherRepository.Add(weather)
}
