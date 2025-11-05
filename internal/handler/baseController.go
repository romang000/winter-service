package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/romang000/winter-service/internal/models"
	"net/http"
)

func (c *Controller) Ping(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (c *Controller) AddWeather(ctx *gin.Context) {
	var weather models.Weather
	err := ctx.BindJSON(&weather)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if err = c.weatherService.AddWeather(weather); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	ctx.JSON(http.StatusOK, gin.H{"message": "weather added"})
}
