package appdividend

import (
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"

	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"

	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

func fetchTinkoffDividendsFor(securityService *tinkoff.InstrumentsServiceClient, stock security.Security) []dividend.Dividend {

	earliestDividendDate := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
	upcomingDividendDate := time.Now().AddDate(2, 0, 0)

	tinkoffDividends, err := securityService.GetDividents(stock.GetFigi(), earliestDividendDate, upcomingDividendDate)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []dividend.Dividend{}
	}

	parsedDividends := []dividend.Dividend{}

	for _, dividend := range tinkoffDividends.GetDividends() {
		if dividend == nil || !dividendIsValid(dividend) {
			logger.Log("The provided dividend is invalid: "+dividend.String(), logger.ERROR)
			continue
		}
		mappedDiv, err := mapTinkoffDividendToDividend(dividend, stock.GetId())
		if err != nil {
			continue
		}
		parsedDividends = append(parsedDividends, mappedDiv)
	}

	return parsedDividends
}

func dividendIsValid(dividend *investapi.Dividend) bool {
	dividendIsValid := true

	if dividend.GetDividendNet() == nil ||
		dividend.DividendNet.GetCurrency() == "" ||
		dividend.DeclaredDate == nil ||
		dividend.RecordDate == nil {
		dividendIsValid = false
	}

	return dividendIsValid
}
