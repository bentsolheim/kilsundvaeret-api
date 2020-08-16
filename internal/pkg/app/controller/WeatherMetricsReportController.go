package controller

import (
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/bentsolheim/kilsundvaeret-api/internal/pkg/app/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type WeatherReportResource struct {
	Metrics []WeatherMetricResource
}

type WeatherMetricResource struct {
	Type        service.MetricType
	CreatedDate int64
	Value       float64
}

func toResource(wr *service.WeatherReport) WeatherReportResource {
	var metrics []WeatherMetricResource
	for _, m := range wr.Metrics {
		metrics = append(
			metrics,
			WeatherMetricResource{
				Type:        m.Type,
				CreatedDate: m.CreatedDate.Unix() * 1000,
				Value:       m.Value,
			},
		)
	}
	return WeatherReportResource{Metrics: metrics}
}

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
		ctx.JSON(http.StatusOK, rest.WrapResponse([]WeatherReportResource{toResource(report)}, nil))
	}

}
