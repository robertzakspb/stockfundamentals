package api_security

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
)

type SecurityDTO struct {
	Isin         string
	Figi         string
	CompanyName  string
	IsPublic     bool
	SecurityType string
	Country      string
	Ticker       string
	IssueSize    int
	Sector       string
	MIC          string
}

func mapStockToDto(stock security.Stock) SecurityDTO {
	return SecurityDTO{
		Isin:         stock.Isin,
		Figi:         stock.Figi,
		CompanyName:  stock.CompanyName,
		IsPublic:     stock.IsPublic,
		SecurityType: string(stock.SecurityType),
		Country:      stock.Country,
		Ticker:       stock.Ticker,
		IssueSize:    stock.IssueSize,
		Sector:       stock.Sector,
		MIC:          stock.MIC,
	}
}

// func mapDtoToStock(dto SecurityDTO) (security.Stock, error) {
// 	securityType, found := security.SecurityTypeMap[dto.SecurityType]
// 	if !found {
// 		return security.Stock{}, fmt.Errorf("unknown security type: %s", dto.SecurityType)
// 	}

// 	return security.Stock{
// 		Id:           dto.Id,
// 		Isin:         dto.Isin,
// 		Figi:         dto.Figi,
// 		CompanyName:  dto.CompanyName,
// 		IsPublic:     dto.IsPublic,
// 		SecurityType: securityType,
// 		Country:      dto.Country,
// 		Ticker:       dto.Ticker,
// 		IssueSize:    dto.IssueSize,
// 		Sector:       dto.Sector,
// 		MIC:          dto.MIC,
// 	}, nil
// }
