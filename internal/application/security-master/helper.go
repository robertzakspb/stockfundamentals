package security_master

import "github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"

func ExtractFigisFromSecurities(securities []security.Stock) []string {
	figis := []string{}
	for _, security := range securities {
		figis = append(figis, security.GetFigi())
	}
	return figis
}
