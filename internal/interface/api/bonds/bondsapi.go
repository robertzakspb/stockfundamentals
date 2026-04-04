package bondsapi

import (
	"net/http"

	bondservice "github.com/compoundinvest/stockfundamentals/internal/application/bondservice"
	"github.com/gin-gonic/gin"
)

func StartBondAndCouponImportJob(c *gin.Context) {
	go bondservice.ImportAllBondsAndCoupons()

	c.JSON(http.StatusOK, "The bond import job has been successfully started")
}

func StartBondAccruedInterestUpdateJob(c *gin.Context) {
	go bondservice.UpdateAllBondsAci()
	c.JSON(http.StatusOK, "The bonds' accrued interest job has been successfully started")
}

func GetRussianGovernmentBondsWithFixedOrConstantCoupon(c *gin.Context) {
	bonds, err := bondservice.GetRussianGovernmentBondsWithFixedOrConstantCoupon()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, bonds)
}

func GetQuasiForeignBonds(c *gin.Context) {
	bonds, err := bondservice.GetQuasiForeignBonds()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, bonds)
}
