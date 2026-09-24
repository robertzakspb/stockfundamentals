package aumapi

import accountmvservice "github.com/compoundinvest/stockfundamentals/internal/application/account/market-value"

func mapAumsToDtos(aums []accountmvservice.AUM) []AumDto {
	dtos := []AumDto{}
	for _, aum := range aums {
		dto := AumDto{
			TotalAum: aum.TotalAum,
			Currency: aum.Currency,
			Date:     aum.Date,
		}
		dtos = append(dtos, dto)
	}
	return dtos
}
