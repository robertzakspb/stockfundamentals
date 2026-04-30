package appdividend

import (
	"errors"
	"strings"
	"time"

	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/dbdividend"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

func mapDividendForecastToDbModel(forecasts []dividend.DividendForecast) []dbdividend.DividendForecastDb {
	dbModels := []dbdividend.DividendForecastDb{}

	for _, forecast := range forecasts {
		dbModel := dbdividend.DividendForecastDb{
			Id:            forecast.Id,
			Figi:          forecast.Stock.Figi,
			ExpectedDPS:   forecast.ExpectedDPS,
			Currency:      forecast.Currency,
			PaymentPeriod: forecast.PaymentPeriod,
			Author:        forecast.Author,
			Comment:       forecast.Comment,
		}

		dbModels = append(dbModels, dbModel)
	}
	return dbModels
}

func mapDividendForecastDbModelToDomain(dbModels []dbdividend.DividendForecastDb) []dividend.DividendForecast {
	forecasts := []dividend.DividendForecast{}

	figis := []string{}
	for _, forecast := range forecasts {
		figis = append(figis, forecast.Stock.Figi)
	}
	securities, err := security_master.GetSecuritiesFilteredByFigi(figis)
	if err != nil {
		logger.Log("Failed to fetch securities with the target Figis", logger.ERROR)
	}

	for _, forecast := range dbModels {
		var targetStock security.Stock
		for _, stock := range securities {
			if stock.Figi != forecast.Figi {
				continue
			}
			targetStock = stock
		}
		if targetStock.Figi == "" {
			targetStock.Figi = forecast.Figi
			logger.Log("Failed to find a security with figi "+forecast.Figi, logger.ERROR)
		}

		domainStruct := dividend.DividendForecast{
			Id:            forecast.Id,
			Stock:         targetStock,
			ExpectedDPS:   forecast.ExpectedDPS,
			Currency:      forecast.Currency,
			PaymentPeriod: forecast.PaymentPeriod,
			Author:        forecast.Author,
			Comment:       forecast.Comment,
		}

		forecasts = append(forecasts, domainStruct)
	}
	return forecasts
}

func mapTinkoffDividendToDividend(tinkoffDiv *investapi.Dividend, figi string) (dividend.Dividend, error) {

	if figi == "" {
		logger.Log("Missing stock ID in the provided stock for tinkoff dividend: "+tinkoffDiv.GetDeclaredDate().String()+tinkoffDiv.GetRecordDate().String(), logger.WARNING)
		return dividend.Dividend{}, errors.New("Missing figi in the Tinkoff dividend")
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

	return dividend, nil
}

func mapDbModelToDividend(dbModelds []dbdividend.DividendDbModel) []dividend.Dividend {
	dividends := make([]dividend.Dividend, len(dbModelds))
	for i, dbModel := range dbModelds {
		newDiv := dividend.Dividend{
			Id:                dbModel.Id,
			Figi:              dbModel.Figi,
			ExpectedDPS:       float64(dbModel.ExpectedDpsTimesMillion) / 1_000_000,
			ActualDPS:         float64(dbModel.ActualDPSTimesMillion) / 1_000_000,
			Currency:          dbModel.Currency,
			AnnouncementDate:  dbModel.AnnouncementDate,
			RecordDate:        dbModel.RecordDate,
			PayoutDate:        dbModel.PayoutDate,
			PaymentPeriod:     dbModel.PaymentPeriod,
			ManagementComment: dbModel.ManagementComment,
		}
		dividends[i] = newDiv
	}

	return dividends
}
