package lot

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetLots(ctx *gin.Context) {
	db := ctx.MustGet("DB").(*sql.DB)

	lots, err := GetLotsFromDB(db)
	if err != nil {
		fmt.Println("failed to fetch the user's portfolio: ", err)
	}

	ctx.JSON(http.StatusOK, lots)
}
