package aumapi

import (
	"net/http"
	"strings"

	accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func GetAssetsUnderManagement(c *gin.Context) {
	currencies, err := shared.GetFromQueryParams("currencies", c.Request.URL.Query())
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
	}

	aum, err := accountmvservice.GetTotalAssetsUnderManagement(strings.Split(currencies, ",")...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ErrorResponse{Errors: []string{err.Error()}})
	}

	aumDto := mapAumsToDtos(aum)

	c.JSON(http.StatusOK, aumDto)
}
