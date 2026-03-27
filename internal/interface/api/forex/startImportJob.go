package forexapi

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/forexservice"
	"github.com/gin-gonic/gin"
)

func StartForexImportJob(c *gin.Context) {
	go forexservice.StartFxRateImportJob()

	c.JSON(http.StatusOK, "The forex import job has been successfully started")
}
