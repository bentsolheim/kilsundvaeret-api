package controller

import (
	"github.com/bentsolheim/go-app-utils/rest"
	"github.com/bentsolheim/go-app-utils/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ForwardJsonResponse(c *gin.Context, url string) {
	body, err := utils.HttpGet(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, rest.WrapResponse(nil, err))
		return
	}
	c.Data(http.StatusOK, "application/json", body)
}
