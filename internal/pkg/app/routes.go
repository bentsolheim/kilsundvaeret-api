package app

import (
	"fmt"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
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
	}

	return r
}
