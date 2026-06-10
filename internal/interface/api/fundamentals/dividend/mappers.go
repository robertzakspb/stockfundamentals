package apidividend

import (
	"sort"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
)

func mapDividendToDTO(dividends []dividend.Dividend) []DividendDTO {
	dtos := []DividendDTO{}
	for _, dividend := range dividends {
		dto := DividendDTO{
			Figi:              dividend.Id.String(),
			ActualDPS:         dividend.ActualDPS,
			ExpectedDPS:       dividend.ExpectedDPS,
			Currency:          dividend.Currency,
			AnnouncementDate:  dividend.AnnouncementDate,
			RecordDate:        dividend.RecordDate,
			PayoutDate:        dividend.PayoutDate,
			PaymentPeriod:     dividend.PaymentPeriod,
			Type:              dividend.Type,
			Regularity:        dividend.Regularity,
			ManagementComment: dividend.ManagementComment,
		}
		dtos = append(dtos, dto)
	}

	return dtos
}

func mapDividendForecastDtoToDomain(dto DividendForecastDTO) dividend.DividendForecast {
	return dividend.DividendForecast{
		Stock:              security.Stock{Figi: dto.Figi},
		ExpectedDPS:        dto.ExpectedDPS,
		Currency:           dto.Currency,
		PaymentPeriod:      dto.PaymentPeriod,
		Author:             dto.Author,
		Comment:            dto.Comment,
		ExpectedPayoutDate: dto.ExpectedPayoutDate,
	}
}

func mapDividendForecastDomainToDto(domains []dividend.DividendForecast) []DividendForecastDTO {
	dtos := []DividendForecastDTO{}
	for _, domain := range domains {

		dto := DividendForecastDTO{
			Figi:               domain.Stock.Figi,
			Ticker:             domain.Stock.Ticker,
			ExpectedDPS:        domain.ExpectedDPS,
			Currency:           domain.Currency,
			PaymentPeriod:      domain.PaymentPeriod,
			Author:             domain.Author,
			Comment:            domain.Comment,
			Yield:              domain.Yield,
			ExpectedPayoutDate: domain.ExpectedPayoutDate,
		}
		dtos = append(dtos, dto)
	}
	sort.Slice(dtos, func(i, j int) bool {
		if dtos[i].ExpectedPayoutDate.Before(dtos[j].ExpectedPayoutDate) {
			return true
		} else {
			return false
		}
	})

	return dtos
}

func mapSecurityDivForecastToDto(forecasts []dividend.SecurityDivForecast) []SecurityDivForecastDto {
	dtos := []SecurityDivForecastDto{}

	for _, forecast := range forecasts {
		dto := SecurityDivForecastDto{
			Figi:             forecast.Figi,
			Forecasts:        mapDividendForecastDomainToDto(forecast.Forecasts),
			CumulativeReturn: forecast.CumulativeReturn(),
		}
		dtos = append(dtos, dto)
	}

	return dtos
}
