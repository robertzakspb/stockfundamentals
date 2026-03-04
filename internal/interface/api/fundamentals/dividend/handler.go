package apidividend

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
)

func StartDividendFetchingJob(c *gin.Context) {
	go appdividend.FetchAndSaveAllDividends()

	c.JSON(http.StatusOK, "Successfully started the dividend fetching job")
}

func GetAllDividends(c *gin.Context) {
	parsedFilters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), dividend.Dividend{})
	dividends, err := dbdividend.GetAllDividends(parsedFilters) //FIXME: Refactor to use the service
	dtos := mapDividendToDTO(dividends)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, dtos)
	}
}

// TODO: Deprecate (use GetAllDividends instead)
func GetUpcomingDividends(c *gin.Context) {
	upcomingDividends, err := dbdividend.GetUpcomingDividends() //FIXME: Refactor to use the service
	dtos := mapDividendToDTO(upcomingDividends)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, dtos)
	}
}

func CreateNewDividendForecast(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	jsonData, err := io.ReadAll(bodyReader)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		logger.Log("Failed to read the dividend forecast json from the POST payload: "+err.Error(), logger.ERROR)
		return
	}

	divForecast := DividendForecastDTO{}
	err = json.Unmarshal(jsonData, &divForecast)
	if err != nil {
		logger.Log("Failed to unmarshal the dividend forecast json in the POST payload: "+err.Error(), logger.ERROR)
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err = appdividend.SaveDividendForecast(mapDividendForecastDtoToDomain(divForecast))
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, "The dividend forecast has been successfully saved")
}

