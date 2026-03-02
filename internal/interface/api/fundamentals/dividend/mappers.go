package apidividend

import 	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/dividend"

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
		Figi: dto.Figi,
		ExpectedDPS: dto.ExpectedDPS,
		Currency: dto.Currency,
		PaymentPeriod: dto.PaymentPeriod,
		Author: dto.Author,
		Comment: dto.Comment,
	}
}
