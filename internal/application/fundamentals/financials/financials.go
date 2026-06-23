package financialsservice

import (
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func ExecuteFundamentalsJob() error {
	securities, err := security_master.GetAllSecuritiesFromDB()
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	for _, security := range securities {
		//Currently no fundamentals are available for non-RU countries
		if security.Country != "RU" {
			continue
		}

	}
}

func FetchFundamentalsForSecurity(security security.Stock) []financials.FinancialMetric {
	
}
