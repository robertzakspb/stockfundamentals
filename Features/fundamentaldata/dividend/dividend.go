package dividend

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	"github.com/google/uuid"
	tinkoff "github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func fetchDividendsForAllStocks(db *ydb.Driver) []Dividend {
	stocks, err := security.FetchSecuritiesFromDBWithDriver(db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	allDividends := []Dividend{}

	rateLimit := time.Second / 1
	throttle := time.Tick(rateLimit)
	for _, stock := range stocks {
		//TODO: - Delete this after testing the throttle functionality
		fmt.Println("Attempting to fetch dividends at:", time.Now(), "for", stock.CompanyName)
		switch stock.GetCountry() {
		case "RU":
			//TODO: Fix this once you figure out what the problem is!
			dividends := FetchTinkoffDividendsFor(security.Stock{Figi: "TCS00A0ZZAC4"})
			allDividends = append(allDividends, dividends...)
			logger.Log("Dividends fetched so far: "+strconv.Itoa(len(allDividends)), logger.INFORMATION)
		default:
			logger.Log("No data provider may provide dividends for "+stock.GetCompanyName(), logger.INFORMATION)
		}
		<-throttle
	}

	return allDividends
}

func FetchTinkoffDividendsFor(stock security.Security) []Dividend {
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

	fmt.Println("Record date: ", time.Unix(tinkoffDiv.GetRecordDate().GetSeconds(), 0).String())
	fmt.Println("Announcement date: ", time.Unix(tinkoffDiv.GetDeclaredDate().GetSeconds(), 0).String())
	fmt.Println("Payment date: ", time.Unix(tinkoffDiv.GetPaymentDate().GetSeconds(), 0).String())

	return Dividend{
		Id:                uuid.New().String(),
		StockID:           stockID,
		ActualDPS:         tinkoffDiv.DividendNet.ToFloat(),
		ExpectedDPS:       0,
		Currency:          tinkoffDiv.DividendNet.GetCurrency(),
		AnnouncementDate:  time.Unix(tinkoffDiv.GetDeclaredDate().GetSeconds(), 0),
		RecordDate:        time.Unix(tinkoffDiv.GetRecordDate().GetSeconds(), 0),
		PayoutDate:        time.Unix(tinkoffDiv.GetPaymentDate().GetSeconds(), 0),
		PaymentPeriod:     "", //TODO: Fix
		ManagementComment: "",
	}
}
