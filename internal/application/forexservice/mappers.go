package forexservice

import forexdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/forex"

func mapFxRatesToDbModel(rates []ForexRate) []forexdb.ForexRateDb {
	dbModels := []forexdb.ForexRateDb{}
	for _, rate := range rates {
		dbModel := forexdb.ForexRateDb{
			Currency1: string(rate.Currency1),
			Currency2: string(rate.Currency2),
			Date:      rate.Date,
			Rate:      rate.Rate,
		}
		dbModels = append(dbModels, dbModel)
	}
	return dbModels
}

func mapDbModelToDomain(rates []forexdb.ForexRateDb) []ForexRate {
	domainModels := []ForexRate{}
	for _, rate := range rates {
		domain := ForexRate{
			Currency1: currencyName[rate.Currency1],
			Currency2: currencyName[rate.Currency2],
			Date:      rate.Date,
			Rate:      rate.Rate,
		}
		domainModels = append(domainModels, domain)
	}
	return domainModels
}
