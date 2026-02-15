package portfolio

import (
	"errors"

	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	pb "opensource.tbank.ru/invest/invest-go/proto"
)

func UpdatePortfolio() error {
	portfolio := GeMyPortfolio()
	if len(portfolio.Lots) == 0 {
		return errors.New("The new portfolio is empty")
	}

	return portfoliodb.UpdateLocalPortfolio(mapLotToDbLot(portfolio.Lots))
}

func FetchPositionSecurities(positions []*pb.PortfolioPosition) []security.Stock {
	figis := []string{}
	for _, position := range positions {
		figis = append(figis, position.Figi)
	}

	securities, err := security_master.GetSecuritiesFilteredByFigi(figis)
	if err != nil || len(securities) == 0 {
		logger.Log("Failed to find positions with the required figis: ", logger.ERROR)
		return []security.Stock{}
	}

	return securities
}
