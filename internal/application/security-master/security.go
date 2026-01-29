package security_master

import (
	"context"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/config"
	securitydb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
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

	securities := []security.Security{}
	for _, stock := range stocks {
		securities = append(securities, stock)
	}

	securitydb.SaveSecuritiesToDB(securities, db)

	return nil
}

func GetAllSecuritiesFromDB() ([]security.Stock, error) {
	return securitydb.GetAllSecuritiesFromDB()
}

func GetSecuritiesFilteredByFigi(figis []string) ([]security.Stock, error) {
	return securitydb.GetSecuritiesFilteredByFigi(figis)
}

func fetchTinkoffSecurities() []security.Stock {
	config, err := tinkoff.LoadConfig("tinkoffAPIconfig.yaml")
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []security.Stock{}
	}

	client, err := tinkoff.NewClient(context.TODO(), config, nil)
	if err != nil {
		logger.Log(err.Error(), logger.ALERT)
		return []security.Stock{}
	}

	securityService := client.NewInstrumentsServiceClient()
	if securityService == nil {
		logger.Log("Security service is nil!", logger.ALERT)
		return []security.Stock{}
	}

	securities, err := securityService.Shares(investapi.InstrumentStatus_INSTRUMENT_STATUS_BASE)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return []security.Stock{}
	}

	russianStocks := []security.Stock{}
	for _, tinkoffStock := range securities.Instruments {
		if tinkoffStock.CountryOfRisk != "RU" {
			continue //Only Russian stocks are being imported for now
		}

		russianStock := security.Stock{
			CompanyName:  tinkoffStock.Name,
			IsPublic:     true,
			Isin:         tinkoffStock.Isin,
			Figi:         tinkoffStock.Figi,
			SecurityType: mapTinkoffSecurityTypeToInternal(tinkoffStock.ShareType),
			Country:      tinkoffStock.CountryOfRisk,
			Ticker:       tinkoffStock.Ticker,
			IssueSize:    int(tinkoffStock.GetIssueSize()),
			Sector:       tinkoffStock.Sector,
			MIC:          "MISX",
		}

		russianStocks = append(russianStocks, russianStock)
	}

	return russianStocks
}

// func GetSecuritiesById(ids []string) ([]security.Security, error) {
// 	return securitydb.GetSecuritiesFilteredById(ids)
// }

func mapTinkoffSecurityTypeToInternal(shareType investapi.ShareType) security.SecurityType {
	switch shareType {
	case investapi.ShareType_SHARE_TYPE_COMMON:
		return security.OrdinaryShare
	case investapi.ShareType_SHARE_TYPE_PREFERRED:
		return security.PreferredShare
	case investapi.ShareType_SHARE_TYPE_ADR, investapi.ShareType_SHARE_TYPE_GDR:
		return security.DepositoryReceipt
	default:
		return security.Unspecified
	}
}
