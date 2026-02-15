package appdividend

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	securitydb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	tinkoff "opensource.tbank.ru/invest/invest-go/investgo"
	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

func FetchAndSaveAllDividends() error {
	dividends := FetchDividendsForAllStocks()

	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		return errors.New("Unable to fetch dividends due to internal configuration issues")
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return errors.New("Unable to fetch dividends due to internal database issues")
	}

	err = dbdividend.SaveDividendsToDB(dividends, db)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return errors.New(err.Error())
	}

	return nil
}

func FetchDividendsForAllStocks() []dividend.Dividend {
	stocks, err := securitydb.GetAllSecuritiesFromDB()
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
	}

	allDividends := []dividend.Dividend{}

	rateLimit := time.Second / 2
	throttle := time.Tick(rateLimit)

	for _, stock := range stocks {
		switch stock.GetCountry() {
		case "RU":
			dividends := fetchTinkoffDividendsFor(stock)
			allDividends = append(allDividends, dividends...)
			logger.Log(fmt.Sprintf("Fetched %d dividends for %s", len(dividends), stock.CompanyName), logger.INFORMATION)
		default:
			logger.Log("No data provider may provide dividends for "+stock.GetCompanyName(), logger.INFORMATION)
			continue
		}
		<-throttle
	}

	return allDividends
}

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
