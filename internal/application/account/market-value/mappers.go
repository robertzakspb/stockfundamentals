package accountmvservice

import (
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	accountmvdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/market-value"
)

func mapAccountMarketValuesToDbModels(mvs []accountmvdomain.AccountMarketValue) []accountmvdb.AccountMarketValueDB {
	dbModels := make([]accountmvdb.AccountMarketValueDB, len(mvs))

	for i := range mvs {
		dbModel := accountmvdb.AccountMarketValueDB{
			AccountId: mvs[i].AccountId,
			Currency:  mvs[i].Currency,
			Date:      mvs[i].Date,
			EodValue:  mvs[i].EodValue,
		}
		dbModels[i] = dbModel
	}
	return dbModels
}

func mapAccountMarketValueDbModelToDomain(mv accountmvdb.AccountMarketValueDB) accountmvdomain.AccountMarketValue {
	return accountmvdomain.AccountMarketValue{
		AccountId: mv.AccountId,
		Currency:  mv.Currency,
		Date:      mv.Date,
		EodValue:  mv.EodValue,
	}
}
