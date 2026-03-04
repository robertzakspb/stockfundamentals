package apidividend

import (
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
			ManagementComment: dividend.ManagementComment,
		}
		dtos = append(dtos, dto)
	}

	return dtos
}

func mapDividendForecastDtoToDomain(dto DividendForecastDTO) dividend.DividendForecast {
	return dividend.DividendForecast{
		Stock:         security.Stock{Figi: dto.Figi},
		ExpectedDPS:   dto.ExpectedDPS,
		Currency:      dto.Currency,
		PaymentPeriod: dto.PaymentPeriod,
		Author:        dto.Author,
		Comment:       dto.Comment,
	}
}

func mapDividendForecastDomainToDto(domains []dividend.DividendForecast) []DividendForecastDTO {
	dtos := []DividendForecastDTO{}
	for _, domain := range domains {

		dto := DividendForecastDTO{
			Figi:          domain.Stock.Figi,
			ExpectedDPS:   domain.ExpectedDPS,
			Currency:      domain.Currency,
			PaymentPeriod: domain.PaymentPeriod,
			Author:        domain.Author,
			Comment:       domain.Comment,
			Yield:         domain.Yield,
		}
		dtos = append(dtos, dto)
	}
	return dtos
}
