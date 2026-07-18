package accountmvservice

import (
	"errors"
	"sync"
	"time"

	bondportfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/bond-portfolio"
	portfolio "github.com/compoundinvest/stockfundamentals/internal/application/account/stock-portfolio"
	accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"
	"github.com/google/uuid"
)

// Uses the latest quotes to provide the latest account market value
func GetCurrentAccountMarketValue(accountId uuid.UUID, currency string) (accountmvdomain.AccountMarketValue, error) {
	if accountId == uuid.Nil {
		return accountmvdomain.AccountMarketValue{}, errors.New("Attempting to get the current market value for a nil account ID")
	}

	var bondMV, stockMV accountmvdomain.AccountMarketValue
	var err error
	var wg sync.WaitGroup

	wg.Go(func() {
		bondPortfolio, e := bondportfolio.GetAccountPositions(accountId)
		err = e //For some reason cannot assign the return error above to the err variable declared outside of the scope
		bondMV, err = CalculateBondLotsMarketValue(bondPortfolio, time.Now(), currency)
	})
	wg.Go(func() {
		stockPortfolio, e := portfolio.GetAccountPortfolio(accountId)
		err = e //For some reason cannot assign the return error above to the err variable declared outside of the scope
		mv, _, e := portfolio.CalculatePortfolioMarketValue(stockPortfolio, currency)
		err = e
		stockMV = accountmvdomain.AccountMarketValue{
			AccountId: accountId,
			Date:      time.Now(),
			Currency:  currency,
			EodValue:  mv,
		}
	})
	wg.Wait()

	if err != nil {
		return accountmvdomain.AccountMarketValue{}, err
	}

	jointMV := accountmvdomain.AccountMarketValue{
		AccountId: stockMV.AccountId,
		Date:      time.Now(),
		Currency:  stockMV.Currency,
		EodValue:  stockMV.EodValue + bondMV.EodValue,
	}
	return jointMV, nil
}
