package accountmvservice

import (
	"time"

	accountservice "github.com/compoundinvest/stockfundamentals/internal/application/account/account"

	accountmvdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/market-value"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
)

func SaveAccountMarketValueSnapshots() error {
	accounts, err := accountservice.GetAllAccounts()
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	for _, account := range accounts {
		go func() {
			accountMVs, err := CalculateAccountMarketValue(account.Id, time.Now())
			if err != nil {
				logger.Log(err.Error(), logger.ERROR)
				return
			}
			mappedMVs := mapAccountMarketValuesToDbModels(accountMVs)
			go accountmvdb.SaveMarketValue(mappedMVs)
		}()
	}
	return nil
}
