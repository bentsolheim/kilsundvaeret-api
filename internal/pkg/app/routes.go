package app

import (
	"fmt"
	"github.com/Depado/ginprom"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

func CreateGinEngine(config AppConfig) *gin.Engine {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/current-temp", func(c *gin.Context) {

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/logger/bua/readings", config.DataLoggerUrl))
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}
			c.String(http.StatusOK, string(body))
		})
		v1.GET("/current-debug", func(c *gin.Context) {

			resp, err := http.Get(fmt.Sprintf("%s/api/v1/logger/bua/debug", config.DataLoggerUrl))
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
				return
			}

			c.String(http.StatusOK, string(body))
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
