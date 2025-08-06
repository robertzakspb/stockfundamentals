package api_security

import (
	"github.com/compoundinvest/stockfundamentals/internal/application/security-master"
)

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
