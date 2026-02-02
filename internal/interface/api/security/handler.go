package api_security

import (
	"net/http"

	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/gin-gonic/gin"
)

func StartSecurityMasterImportJob(c *gin.Context) {
	go security_master.FetchAndSaveSecurities()

	c.JSON(http.StatusOK, "Successfully started the security import job")
}

func FetchSecuritiesFromDB() ([]SecurityDTO, error) {
	securities, err := security_master.GetAllSecuritiesFromDB()
	if err != nil {
		return []SecurityDTO{}, err
	}

	dtos := []SecurityDTO{}
	for _, security := range securities {
		dtos = append(dtos, mapStockToDto(security))
	}

	return dtos, nil
}
