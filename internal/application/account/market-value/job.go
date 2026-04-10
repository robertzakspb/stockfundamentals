package accountmvservice

import (
	"time"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"

	accountmvdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func SaveAccountMarketValueSnapshots() {
	accounts, err := accountservice.GetAllAccounts()
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}

	for _, account := range accounts {
		go func() {
			accountMV, err := CalculateAccountMarketValue(account.Id, time.Now())
			if err != nil {
				logger.Log(err.Error(), logger.ERROR)
			}
			mappedMV := mapAccountMarketValueToDbModel(accountMV)
			accountmvdb.SaveMarketValue([]accountmvdb.AccountMarketValueDB{mappedMV})
		}()
	}
}
