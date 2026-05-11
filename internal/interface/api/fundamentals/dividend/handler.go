package apidividend

import (
	"encoding/json"
	"io"
	"net/http"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartDividendFetchingJob(c *gin.Context) {
	go appdividend.FetchAndSaveAllDividends()

	c.JSON(http.StatusOK, shared.StringResponse{Message: "Successfully started the dividend fetching job"})
}

func GetAllDividends(c *gin.Context) {
	parsedFilters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), dividend.Dividend{})

	dividends, err := appdividend.GetFilteredDividends(parsedFilters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, shared.ResponseError{Errors: []string{err.Error()}})
		return
	}

	dtos := mapDividendToDTO(dividends)

	c.JSON(http.StatusOK, dtos)
}

func CreateNewDividendForecast(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	jsonData, err := io.ReadAll(bodyReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ResponseError{Errors: []string{err.Error()}})
		logger.Log("Failed to read the dividend forecast json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	divForecast := DividendForecastDTO{}
	err = json.Unmarshal(jsonData, &divForecast)
	if err != nil {
		logger.Log("Failed to unmarshal the dividend forecast json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, shared.ResponseError{Errors: []string{err.Error()}})
		return
	}

	err = appdividend.SaveDividendForecast(mapDividendForecastDtoToDomain(divForecast))
	if err != nil {
		c.JSON(http.StatusBadRequest, shared.ResponseError{Errors: []string{err.Error()}})
		return
	}

	c.JSON(http.StatusOK, shared.StringResponse{Message: "The dividend forecast has been successfully saved"})
}
