package appdividend

import (
	"context"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"

	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

func fetchTinkoffDividendsFor(stock security.Security) []dividend.Dividend {
	//TODO: Extract the file name into an environment variable
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []dividend.Dividend{}
	}

	client, err := tinkoff.NewClient(context.TODO(), config, nil)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []dividend.Dividend{}
	}

	securityService := client.NewInstrumentsServiceClient()
	parsedDividends := []dividend.Dividend{}
	earliestDividendDate := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
	upcomingDividendDate := time.Now().AddDate(2, 0, 0)

	tinkoffDividends, err := securityService.GetDividents(stock.GetFigi(), earliestDividendDate, upcomingDividendDate)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []dividend.Dividend{}
	}

	for _, dividend := range tinkoffDividends.GetDividends() {
		if dividend == nil || !dividendIsValid(dividend) {
			logger.Log("The provided dividend is invalid: "+dividend.String(), logger.ERROR)
			continue
		}

		parsedDividends = append(parsedDividends, mapTinkoffDividendToDividend(dividend, stock.GetId()))
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

func mapTinkoffDividendToDividend(tinkoffDiv *investapi.Dividend, figi string) dividend.Dividend {

	if figi == "" {
		logger.Log("Missing stock ID in the provided stock for tinkoff dividend: "+tinkoffDiv.GetDeclaredDate().String()+tinkoffDiv.GetRecordDate().String(), logger.WARNING)
	}

	dividend := dividend.Dividend{
		Id:                uuid.New(),
		Figi:              figi,
		ActualDPS:         tinkoffDiv.DividendNet.ToFloat(),
		ExpectedDPS:       0,
		Currency:          strings.ToUpper(tinkoffDiv.DividendNet.GetCurrency()),
		AnnouncementDate:  time.Unix(tinkoffDiv.GetDeclaredDate().GetSeconds(), 0),
		RecordDate:        time.Unix(tinkoffDiv.GetRecordDate().GetSeconds(), 0),
		PayoutDate:        time.Unix(tinkoffDiv.GetPaymentDate().GetSeconds(), 0),
		PaymentPeriod:     "", //TODO: Fix
		ManagementComment: "",
	}

	return dividend
}
