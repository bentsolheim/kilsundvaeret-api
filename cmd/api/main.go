package main

import (
	"fmt"
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/api/v1/current-temp", func(c *gin.Context) {

		dataLoggerUrl := utils.GetEnvOrDefault("DATALOGGER_URL", "http://datalogger.kilsundvaeret.no")
		resp, err := http.Get(fmt.Sprintf("%s/api/v1/logger/bua/readings", dataLoggerUrl))
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
	r.Run(":9010")
}
