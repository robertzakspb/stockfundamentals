package portfolio

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/portfolio"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetPortfolio(c *gin.Context) {
	userPortfolio := portfolio.GeMyPortfolio().WithPLs()

	if len(userPortfolio) > 0 {
		c.JSON(http.StatusOK, userPortfolio)
	}
	c.JSON(http.StatusNoContent, "User does not have any positions")
}

func GetAccountPortfolio(c *gin.Context) {
	accountIDs := uuid.UUIDs{} //FIXME: Parse from the query parameters
	portfolio, err := portfolio.GetAccountPortfolio(accountIDs)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.JSON(http.StatusOK, portfolio)
}

func UpdatePortfolio(c *gin.Context) {
	
}
