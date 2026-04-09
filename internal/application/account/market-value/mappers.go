package accountmv

import (
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/dividend-calendar/market-value"
	accountmvdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/market-value"
)

func mapAccountMarketValueToDbModel(mv accountmvdomain.AccountMarketValue) accountmvdb.AccountMarketValueDB {
	return accountmvdb.AccountMarketValueDB{
		AccountId: mv.AccountId,
		Currency:  mv.Currency,
		Date:      mv.Date,
		EodValue:  mv.EodValue,
	}
}

func mapAccountMarketValueDbModelToDomain(mv accountmvdb.AccountMarketValueDB) accountmvdomain.AccountMarketValue{
	return accountmvdomain.AccountMarketValue{
		AccountId: mv.AccountId,
		Currency:  mv.Currency,
		Date:      mv.Date,
		EodValue:  mv.EodValue,
	}
}
