package accountmvservice

import (
	"time"

	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/google/uuid"
)

func GetAccountStockMarketValueGroupedByCurrency(accountId uuid.UUID, date time.Time) (map[string]accountmvdomain.AccountMarketValue, error) {
	accountPortfolio, err := portfolio.GetAccountPortfolio(accountId)
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}
	if len(accountPortfolio.Lots) == 0 {
		return map[string]accountmvdomain.AccountMarketValue{}, nil
	}

	accountPortfolio.Lots, err = portfolio.PopulateLotSecurities(accountPortfolio.Lots)
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}
	if len(accountPortfolio.Lots) == 0 {
		return map[string]accountmvdomain.AccountMarketValue{}, nil
	}

	stockPortfolioMarketValue, currency, err := portfolio.CalculatePortfolioMarketValue(accountPortfolio, accountPortfolio.Lots[0].Currency)
	if err != nil {
		return map[string]accountmvdomain.AccountMarketValue{}, err
	}

	mv := map[string]accountmvdomain.AccountMarketValue{currency: {
		AccountId: accountId,
		Date:      date,
		Currency:  currency,
		EodValue:  stockPortfolioMarketValue,
	}}
	return mv, nil
}
