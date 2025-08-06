package portfolio

import (
	"github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	security_master "github.com/compoundinvest/stockfundamentals/internal/application/security-master"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	pb "github.com/russianinvestments/invest-api-go-sdk/proto"
)

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
