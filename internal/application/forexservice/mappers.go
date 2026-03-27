package forexservice

import forexdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/forex"

func mapFxRatesToDbModel(rates[] ForexRate) []forexdb.ForexRateDb {
	dbModels := []forexdb.ForexRateDb{}
	for _, rate := range rates {
		dbModel := forexdb.ForexRateDb{
			Currency1: rate.Currency1,
			Currency2: rate.Currency2,
			Date: rate.Date,
			Rate: rate.Rate,
		}
		dbModels = append(dbModels, dbModel)
	}
	return dbModels
}