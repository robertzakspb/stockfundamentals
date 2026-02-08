package apidividend

import (
	"net/http"

	appdividend "github.com/compoundinvest/stockfundamentals/internal/application/fundamentals/dividend"
	"github.com/gin-gonic/gin"
)

func StartDividendFetchingJob(c *gin.Context) {
	go appdividend.FetchAndSaveAllDividends()

	c.JSON(http.StatusOK, "Successfully started the dividend fetching job")
}

//Refactor to use gin
// func GetAllDividends() ([]dividend.Dividend, error) {
// 	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
// 	defer cancel()

// 	config, err := config.LoadConfig()
// 	if err != nil {
// 		return []dividend.Dividend{}, err
// 	}

// 	db, err := ydb.Open(ctx, config.DB.ConnectionString)

// 	if err != nil {
// 		logger.Log(err.Error(), logger.ALERT)
// 		panic("Failed to connect to the database")
// 	}
// 	dividends, _ := dbdividend.GetAllDividends(db)

// 	return dividends, nil
// }
