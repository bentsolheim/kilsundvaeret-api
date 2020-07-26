package main

import (
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/api/v1/current-temp", func(c *gin.Context) {
		resp, err := http.Get("http://datalogger.kilsundvaeret.no/api/v1/logger/bua/readings")
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
