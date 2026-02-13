package portfolio

import ("github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/portfoliodb"
"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
)

func mapLotDbToLot(lotDb portfoliodb.LotDb) lot.Lot {
	return lot.Lot {
		Id: lotDb.Id,
		CreatedAt: lotDb.CreatedAt,
		UpdatedAt: lotDb.UpdatedAt,
		Quantity: lotDb.Quantity,
		PricePerUnit: lotDb.PricePerUnit,
		Currency: lotDb.Currency,
		AccountId: lotDb.AccountID,
		SecurityId: lotDb.Figi,
		CurrentPL: 0,
	}
}