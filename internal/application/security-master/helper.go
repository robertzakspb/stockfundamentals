package security_master

import (
	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
)

func ExtractFigisFromSecurities(securities []security.Stock) []string {
	figis := []string{}
	for _, security := range securities {
		figis = append(figis, security.GetFigi())
	}
	return figis
}

func ConvertStocksToSecurities(stocks []security.Stock) []entity.Security {
	entitySecurities := []entity.Security{}
	for i := range stocks {
		entitySecurities = append(entitySecurities, entity.Security{
			Figi:   stocks[i].Figi,
			ISIN:   stocks[i].Isin,
			Ticker: stocks[i].Ticker,
			MIC:    stocks[i].MIC,
		})
	}
	return entitySecurities
}
