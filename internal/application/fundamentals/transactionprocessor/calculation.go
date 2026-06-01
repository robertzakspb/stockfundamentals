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
		if t.IsBuyOrder() {
			newLot, err := lot.NewLot()
		}
		if t.IsSellOrder() {

		}
		return account, lots, errors.New("Encountered an unsupported order type: " + string(t.Type))
	}

	return account, lots, nil
}
