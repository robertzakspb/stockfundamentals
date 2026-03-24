package bondsapi

import (
	"net/http"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bonds"
	"github.com/gin-gonic/gin"
)

func StartBondAndCouponImportJob(c *gin.Context) {
	go bondservice.ImportAllBondsAndCoupons()

	c.JSON(http.StatusOK, "The bond import job has been successfully started")
}

func UpdateAllBondsAci(c *gin.Context) {
	err := bondservice.UpdateAllBondsAci()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, "The bonds' ACI has been successfully updated")
}
