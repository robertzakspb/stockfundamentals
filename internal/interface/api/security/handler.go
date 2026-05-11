package api_security

import (
	"net/http"

	"github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/interface/shared"
	"github.com/gin-gonic/gin"
)

func StartSecurityMasterImportJob(c *gin.Context) {
	go security_master.FetchAndSaveSecurities()

	c.JSON(http.StatusOK, shared.StringResponse{Message: "Successfully executed the security import job"})
}
