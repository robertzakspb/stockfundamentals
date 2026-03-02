package apidividend

import (
	"encoding/json"
	"fmt"
	"net/http"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	"github.com/gin-gonic/gin"
)

func StartDividendFetchingJob(c *gin.Context) {
	go appdividend.FetchAndSaveAllDividends()

	c.JSON(http.StatusOK, "Successfully started the dividend fetching job")
}

func GetAllDividends(c *gin.Context) {
	parsedFilters := ydbfilter.MapQueryFiltersToYdb(c.Request.URL.Query(), dividend.Dividend{})
	dividends, err := dbdividend.GetAllDividends(parsedFilters) //FIXME: Refactor to use the service
	dtos := convertDividendToDTO(dividends)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, dtos)
	}
}

// TODO: Deprecate (use GetAllDividends instead)
func GetUpcomingDividends(c *gin.Context) {
	upcomingDividends, err := dbdividend.GetUpcomingDividends() //FIXME: Refactor to use the service
	dtos := convertDividendToDTO(upcomingDividends)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	} else {
		c.JSON(http.StatusOK, dtos)
	}
}

func CreateNewDividendForecast(c *gin.Context) {
	bodyReader := c.Request.Body
	defer bodyReader.Close()

	bodyBytes := [10000]byte{}

	bodyReader.Read(bodyBytes[:])
	bodyString := string(bodyBytes[:])
	fmt.Println(bodyString)

	divForecast := DividendForecastDTO{}
	err := json.Unmarshal(bodyBytes[:], &divForecast)

	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
}
