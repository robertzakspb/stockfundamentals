package accountreturnapi

import accountmvdomain "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/market-value"

func mapAccountMarketValueToDto(mv accountmvdomain.AccountMarketValue) AccountMarketValueDto {
	return AccountMarketValueDto{
		AccountId: mv.AccountId,
		Currency:  mv.Currency,
		Date:      mv.Date,
		EodValue:  mv.EodValue,
	}
}

func mapDomainToDto(domain accountmvdomain.Return) AccountReturnDto {
	return AccountReturnDto{
		AccountId:                domain.AccountId.String(),
		Currency:                 domain.Currency,
		AbsoluteReturn:           domain.AbsoluteReturn,
		AbsoluteReturnPercentage: domain.AbsoluteReturnPercentage,
		AnnualizedReturn:         domain.AnnualizedReturnPercentage,
		StartDate:                domain.StartDate,
		EndDate:                  domain.EndDate,
	}
}
