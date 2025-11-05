package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/romang000/winter-service/internal/models"
	"log/slog"
)

const baseRout = "/api/v1"

type WeatherService interface {
	AddWeather(weather models.Weather) error
}

type Controller struct {
	router         *gin.Engine
	logger         *slog.Logger
	weatherService WeatherService
}

func RegisterRoutes(router *gin.Engine, logger *slog.Logger, weatherService WeatherService) Controller {
	controller := Controller{
		router:         router,
		logger:         logger,
		weatherService: weatherService,
	}
	
	api := router.Group(baseRout)
	
	api.GET("/ping", controller.Ping)
	api.POST("/weather", controller.AddWeather)
	
	return controller
}
