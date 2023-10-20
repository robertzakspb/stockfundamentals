package company

import (
	"database/sql"

	"github.com/gin-gonic/gin"
)

func GetCompanyFromDB(ctx *gin.Context) {
	db := ctx.MustGet("DB").(*sql.DB)

	ticker := "AAPL"
	security := db.QueryRow("SELECT * FROM company WHERE ordinary_share_ticker = $1 OR preferred_share_ticker = $1", ticker)

	var company Company
	security.Scan(&company.ID, &company.Name, &company.SecurityType, &company.Country, &company.IsPublic, &company.OrdinaryShareCount, &company.OrdinaryShareTicker, &company.PreferredShareCount, &company.PreferredShareTicker)
}
