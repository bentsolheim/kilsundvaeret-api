package controller

import (
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewWeatherMetricsController(metricsService *service.WeatherMetricsService) *WeatherMetricsController {
	return &WeatherMetricsController{weatherMetricsService: metricsService}
}

type WeatherMetricsController struct {
	weatherMetricsService *service.WeatherMetricsService
}

func (c WeatherMetricsController) GetAll(ctx *gin.Context) {
	metricType := ctx.Query("type")
	loggerName := ctx.Query("loggerName")
	if metrics, err := c.weatherMetricsService.FindAllFiltered(service.MetricsFilter{
		Type:       &metricType,
		LoggerName: &loggerName,
	}); err != nil {
		ctx.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
	} else {
		ctx.JSON(http.StatusOK, rest.WrapResponse(metrics, nil))
	}
}
