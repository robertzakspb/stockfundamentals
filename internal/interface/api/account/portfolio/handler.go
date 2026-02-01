package portfolio

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/portfolio"
	"github.com/gin-gonic/gin"
)

func GetPortfolio(c *gin.Context) {
	userPortfolio := portfolio.GeMyPortfolio().WithPLs()
	if len(userPortfolio) > 0 {
		c.JSON(http.StatusOK, userPortfolio)
	}

	c.JSON(http.StatusNoContent, "User does not have any positions")
}
