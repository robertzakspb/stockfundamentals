package transactionprocessor

import (
	"errors"
	"sort"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	tranlotrelation "github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/tran-lot-relation"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/google/uuid"
)

func recalculateLotsAndCashBalances(account account.Account, transactions []transaction.Transaction, lots []lot.Lot) (account.Account, []lot.Lot, []tranlotrelation.TransactionLotRelation, error) {
	relations := []tranlotrelation.TransactionLotRelation{}
	if account.Id == uuid.Nil || len(transactions) == 0 || len(lots) == 0 {
		return account, lots, relations, errors.New("Invalid date was provided to the recalculate lots and balances function")
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
			return account, lots, relations, errors.New("Encountered an unsupported transaction type: " + string(t.Type))
		}
		if !(t.IsBuyOrder() || t.IsSellOrder()) {
			return account, lots, relations, errors.New("Encountered an unsupported order side: " + string(t.Side))
		}
		if t.IsBuyOrder() {
			//Each buy transaction begets a separate position lot
			newLot, err := lot.NewLot(t.Figi, t.Quantity, t.PricePerUnit, t.Currency, t.AccountId)
			if err != nil {
				return account, lots, relations, err
			}
			lots = append(lots, newLot)

			relation, err := tranlotrelation.New(newLot.Id, uuid.Nil, t.Id, t.Timestamp, t.Quantity)
			if err != nil {
				return account, lots, relations, err

			}
			relations = append(relations, relation)

			account.CashBalance -= t.Quantity * t.PricePerUnit
			if account.IsCashNegative() {
				return account, lots, relations, err
			}
		}
		if t.IsSellOrder() {
			//A sale might impact 1 to N lots; hence we need to loop through all lots and start closing them one by one until the transaction quantity is exhausted
			quantityToSell := t.Quantity
			targetLotIndices := lot.FindLotIndicesByFigi(lots, t.Figi)
			for i := range targetLotIndices {
				var lot *lot.Lot = &lots[i]
				if lot.IsClosed {
					continue //If a lot is already closed, simply skip it
				}
				//Once all securities were sold, terminate the loop
				if quantityToSell == 0 {
					break
				}

				//In this simple case the entire transaction only impacts the given lot
				if quantityToSell < lot.Quantity {
					relation, err := tranlotrelation.New(lot.Id, uuid.Nil, t.Id, t.Timestamp, quantityToSell*-1)
					if err != nil {
						return account, lots, relations, err

					}
					relations = append(relations, relation)

					lot.Quantity -= quantityToSell //Reducing the lot's quantity by the quantity left to sell
					account.CashBalance += quantityToSell * t.PricePerUnit
					quantityToSell = 0

					break //Terminating the loop, as the transaction has been entirely processed

				} else {
					//In this scenario the given lot is not sufficient to cover the transaction and it's thus closed, and the loop proceeds to cover the next lot
					relation, err := tranlotrelation.New(lot.Id, uuid.Nil, t.Id, t.Timestamp, lot.Quantity*-1)
					if err != nil {
						return account, lots, relations, err
					}
					relations = append(relations, relation)

					quantityToSell -= lot.Quantity                       //Reducing the quantity left to sell by the lot's quantity, as it's being closed
					account.CashBalance += lot.Quantity * t.PricePerUnit //Increasing the balance by the lot's quantity multiplied by the sale price
					lot.Quantity = 0                                     // Setting the lot's quantity to 0, as it's being closed
					lot.IsClosed = true

				}
			}

			if quantityToSell != 0 {
				return account, lots, relations, errors.New("Quantity to sell is not 0 for transaction" + t.Id.String())
			}
		}

	}

	return account, lots, relations, nil
}
