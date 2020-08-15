package controller

import (
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewWeatherMetricsReportController(metricsService *service.WeatherMetricsService) WeatherMetricsReportController {
	return WeatherMetricsReportController{weatherMetricsService: metricsService}
}

type WeatherMetricsReportController struct {
	weatherMetricsService *service.WeatherMetricsService
}

func (c WeatherMetricsReportController) CurrentWeather(ctx *gin.Context) {
	if report, err := c.weatherMetricsService.CurrentWeather(); err != nil {
		ctx.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
	} else {
		ctx.JSON(http.StatusOK, rest.WrapResponse(report, nil))
	}

}
