package timeseries

import (
	"errors"
	"fmt"

	"github.com/compoundinvest/invest-core/quote/entity"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/quote"
	timeseriesdb "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/marketdata"
	ydbfilter "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-filter"
	ydbhelper "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/shared/ydb-helper"
)

func GetLatestLocalQuotesForAllSecurities() ([]quote.Quote, error) {
	quotes := []quote.Quote{}

	dbQuotes, err := timeseriesdb.GetLatestQuotesForAllSecurities()
	if err != nil {
		return quotes, err
	}

	mappedQuotes := mapDbQuotesToQuotes(dbQuotes)
	quotes = append(quotes, mappedQuotes...)

	return quotes, nil
}

func GetLatestLocalQuotesForFigis(figis []string) ([]entity.SimpleQuote, error) {
	quotes := []entity.SimpleQuote{}

	dbQuotes, err := timeseriesdb.GetLatestQuotesForAllSecurities()
	if err != nil {
		return quotes, err
	}

	for i := range dbQuotes {
		for j := range figis {
			if dbQuotes[i].Figi == figis[j] {
				mappedQuote := mapDbQuoteToQuote(dbQuotes[i])
				quotes = append(quotes, mappedQuote)
			}
		}
	}
	return quotes, nil
}

// Provides local bond quotes for bonds with the specified tickers in the DB
func GetLatestLocalBondQuotes(tickers []string) ([]entity.BondQuote, error) {
	quotes := []entity.BondQuote{}
	tickerFilter := ydbfilter.YdbFilter{
		YqlColumnName:  "ticker",
		Condition:      ydbfilter.Contains,
		ConditionValue: ydbhelper.ConvertStringsToYdbList(tickers),
	}

	dbQuotes, err := timeseriesdb.GetLatestBondQuotes([]ydbfilter.YdbFilter{tickerFilter})
	for i := range dbQuotes {
		quotes = append(quotes, dbQuotes[i])
	}

	if err != nil {
		return quotes, err
	}
	if len(tickers) != len(dbQuotes) {
		return quotes, errors.New(fmt.Sprint("Fetched %i quotes for %i tickers", len(dbQuotes), len(tickers)))
	}

	return quotes, nil
}
