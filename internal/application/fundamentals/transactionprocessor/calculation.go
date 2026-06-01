package transactionprocessor

import (
	"errors"
	"sort"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
)

func recalculateLotsAndCashBalances(account account.Account, transactions []transaction.Transaction, lots []lot.Lot) (account.Account, []lot.Lot, error) {
	//Transactions must be processed in chronological order, reenacting the user's behavior in the OMS
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Timestamp.After(transactions[j].Timestamp)
	})
	//Lots must be sorted by creation date because sales are applied according to the FIFO principle
	sort.Slice(lots, func(i, j int) bool {
		return lots[i].CreatedAt.After(lots[j].CreatedAt)
	})
	//TODO: Check the sorting logic!!

	for _, t := range transactions {
		if !(t.IsBuyOrder() || t.IsSellOrder()) {
			return account, lots, errors.New("Encountered an unsupported transaction type: " + string(t.Type))
		}
		if t.IsBuyOrder() {
			newLot, err := lot.NewLot(t.Figi, t.Quantity, t.PricePerUnit, t.Currency, t.AccountId)
			if err != nil {
				return account, lots, err
			}
			lots = append(lots, newLot)
		}
		if t.IsSellOrder() {
			//TODO:
		}

	}

	return account, lots, nil
}
