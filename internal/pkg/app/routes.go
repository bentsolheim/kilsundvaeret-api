package app

import (
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func CreateGinEngine(config AppConfig) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/current-temp", func(c *gin.Context) {
			forwardJsonResponse(c, fmt.Sprintf("%s/api/v1/logger/bua/readings", config.DataReceiverUrl))
		})
		v1.GET("/current-debug", func(c *gin.Context) {
			forwardJsonResponse(c, fmt.Sprintf("%s/api/v1/logger/bua/debug", config.DataReceiverUrl))
		})
	}

	p := ginprom.New(
		ginprom.Namespace(""),
		ginprom.Subsystem(""),
		ginprom.Engine(r),
		ginprom.Path("/api/metrics"),
	)
	r.Use(p.Instrument())
	http.Handle("/metrics", promhttp.Handler())
	return r
}

func forwardJsonResponse(c *gin.Context, url string) {
	body, err := HttpGet(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}
