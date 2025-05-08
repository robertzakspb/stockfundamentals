package security

import (
	"context"
	"time"

	"github.com/compoundinvest/stockfundamentals/infrastructure/config"
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	tinkoff "github.com/russianinvestments/invest-api-go-sdk/investgo"
	investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
	"github.com/ydb-platform/ydb-go-sdk/v3"
)

func FetchAndSaveSecurities() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()

	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	db, err := ydb.Open(ctx, config.DB.ConnectionString)

	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		panic("Failed to connect to the database")
	}
	stocks := fetchTinkoffSecurities()
	if len(stocks) == 0 {
		logger.Log("Fetched 0 securities from Tinkoff API, this is unexpected", logger.ERROR)
	}

	securities := []Security{}
	for _, stock := range stocks {
		securities = append(securities, stock)
	}

	saveSecuritiesToDB(securities, db)

	return nil
}

func fetchTinkoffSecurities() []Stock {

	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []Stock{}
	}

	client, err := tinkoff.NewClient(context.TODO(), config, nil)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []Stock{}
	}

	securityService := client.NewInstrumentsServiceClient()
	if securityService == nil {
		logger.Log("Security service is nil!", logger.ALERT)
		return []Stock{}
	}

	securities, err := securityService.Shares(investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []Stock{}
	}

	russianStocks := []Stock{}
	for _, tinkoffStock := range securities.Instruments {
		if tinkoffStock.CountryOfRisk != "RU" {
			continue //Only Russian stocks are being imported for now
		}

		russianStock := Stock{
			Id:           "",
			CompanyName:  tinkoffStock.Name,
			IsPublic:     true,
			Isin:         tinkoffStock.Isin,
			Figi:         tinkoffStock.Figi,
			SecurityType: mapTinkoffSecurityTypeToInternal(tinkoffStock.ShareType),
			Country:      tinkoffStock.CountryOfRisk,
			Ticker:       tinkoffStock.Ticker,
			IssueSize:    int(tinkoffStock.GetIssueSize()),
			Sector:       tinkoffStock.Sector,
		}

		russianStocks = append(russianStocks, russianStock)
	}

	return russianStocks
}

func mapTinkoffSecurityTypeToInternal(shareType investapi.ShareType) SecurityType {
	switch shareType {
	case investapi.ShareType_SHARE_TYPE_COMMON:
		return OrdinaryShare
	case investapi.ShareType_SHARE_TYPE_PREFERRED:
		return PreferredShare
	case investapi.ShareType_SHARE_TYPE_ADR, investapi.ShareType_SHARE_TYPE_GDR:
		return DepositoryReceipt
	default:
		return Unspecified
	}
}
