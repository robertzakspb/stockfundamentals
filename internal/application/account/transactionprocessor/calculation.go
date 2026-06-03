package transactionprocessor

import (
	"errors"
	"sort"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
	"github.com/google/uuid"
)

func recalculateLotsAndCashBalances(account account.Account, transactions []transaction.Transaction, lots []lot.Lot) (account.Account, []lot.Lot, error) {
	if account.Id == uuid.Nil || len(transactions) == 0 || len(lots) == 0 {
		return account, lots, errors.New("Invalid date was provided to the recalculate lots and balances function")
	}
	//Transactions must be processed in chronological order, reenacting the user's behavior in the OMS
	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].Timestamp.Before(transactions[j].Timestamp)
	})
	//Lots must be sorted by creation date because sales are applied according to the FIFO principle
	sort.Slice(lots, func(i, j int) bool {
		return lots[i].CreatedAt.Before(lots[j].CreatedAt)
	})

	for _, t := range transactions {
		if t.Type != transaction.OrderExecution {
			return account, lots, errors.New("Encountered an unsupported transaction type: " + string(t.Type))
		}
		if !(t.IsBuyOrder() || t.IsSellOrder()) {
			return account, lots, errors.New("Encountered an unsupported order side: " + string(t.Side))
		}
		if t.IsBuyOrder() {
			//Each buy transaction begets a separate transaction lot
			newLot, err := lot.NewLot(t.Figi, t.Quantity, t.PricePerUnit, t.Currency, t.AccountId)
			if err != nil {
				return account, lots, err
			}
			lots = append(lots, newLot)

			account.CashBalance -= t.Quantity * t.PricePerUnit
			if account.IsCashNegative() {
				return account, lots, errors.New("The cash balance of account " + account.Id.String() + " is negative as a result of applying transaction" + t.Id.String())
			}
		}
		if t.IsSellOrder() {
			//A sale might impact 1 to N lots; hence we need to loop through all lots and start closing them one by one until the transaction quantity is exhausted
			quantityToSell := t.Quantity
			targetLotIndices := lot.FindLotIndicesByFigi(lots, t.Figi)
			for i := range targetLotIndices {
				var lot *lot.Lot = &lots[i]
				if lot.IsClosed {
					continue
				}
				//Once all securities were sold, terminate the loop
				if quantityToSell == 0 {
					break
				}

				//In this simple case the entire transaction only impacts the given lot
				if quantityToSell < lot.Quantity {
					lot.Quantity -= quantityToSell
					account.CashBalance += quantityToSell * t.PricePerUnit
					quantityToSell = 0
					break
				}
				//In this scenario the given lot is not sufficient to cover the transaction and it's thus closed, and the loop proceeds to cover the next lot
				quantityToSell -= lot.Quantity                       //Reducing the quantity left to sell by the lot's quantity, as it's being closed
				account.CashBalance += lot.Quantity * t.PricePerUnit //Increasing the balance by the lot's quantity multiplied by the sale price
				lot.Quantity = 0                                     // Setting the lot's quantity to 0, as it's being closed
				lot.IsClosed = true                                  //Closing the lot
			}

			if quantityToSell != 0 {
				return account, lots, errors.New("Quantity to sell is not 0 for transaction" + t.Id.String())
			}
		}

	}

	return account, lots, nil
}
