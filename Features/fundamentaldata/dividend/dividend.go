package dividend

import (
	"context"
	"fmt"
	"time"

	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/google/uuid"
	tinkoff "github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func fetchDividendsForAllStocks(db *ydb.Driver) []Dividend {
	stocks, err := security.FetchSecuritiesFromDB(db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	allDividends := []Dividend{}

	rateLimit := time.Second * 3
	throttle := time.Tick(rateLimit)
	for _, stock := range stocks {
		fmt.Println("Attempting to fetch a dividend at: ", time.Now())
		switch stock.GetCountry() {
		case "RU":
			dividends := fetchTinkoffDividendsFor(stock)
			allDividends = append(dividends, dividends...)
		default:
			logger.Log("No data provider may provide dividends for "+stock.GetCountry(), logger.INFORMATION)
		}
		<-throttle
	}

	return allDividends
}

func fetchTinkoffDividendsFor(stock security.Security) []Dividend {
	//TODO: Extract the file name into an environment variable
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []Dividend{}
	}

	client, err := tinkoff.NewClient(context.TODO(), config, nil)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []Dividend{}
	}

	securityService := client.NewInstrumentsServiceClient()
	parsedDividends := []Dividend{}
	earliestDividendDate := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC)
	upcomingDividendDate := time.Now().AddDate(2, 0, 0)

	tinkoffDividends, err := securityService.GetDividents(stock.GetFigi(), earliestDividendDate, upcomingDividendDate)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []Dividend{}
	}

	for _, dividend := range tinkoffDividends.GetDividends() {
		if dividend == nil || !dividendIsValid(dividend) {
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

func mapTinkoffDividendToDividend(tinkoffDiv *investapi.Dividend, stockID string) Dividend {

	if stockID == "" {
		logger.Log("Missing stock ID in the provided stock for tinkoff dividend: "+tinkoffDiv.GetDeclaredDate().String()+tinkoffDiv.GetRecordDate().String(), logger.WARNING)
	}

	return Dividend{
		Id:                uuid.New().String(),
		StockID:           stockID,
		ActualDPS:         tinkoffDiv.DividendNet.ToFloat(),
		ExpectedDPS:       0,
		Currency:          tinkoffDiv.DividendNet.GetCurrency(),
		AnnouncementDate:  time.Unix(tinkoffDiv.RecordDate.Seconds, 0),
		RecordDate:        time.Unix(tinkoffDiv.GetRecordDate().Seconds, 0),
		PayoutDate:        time.Unix(tinkoffDiv.GetPaymentDate().Seconds, 0),
		PaymentPeriod:     "", //TODO: Fix
		ManagementComment: "",
	}
}
