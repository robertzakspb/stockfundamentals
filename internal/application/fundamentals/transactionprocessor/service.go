package transactionprocessor

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/transaction"
)

//Main method that accepts transactions of any type and processes them
func ProcessTransactions(transactions []transaction.Transaction) error {
	for i := range transactions {
		switch transactions[i].Type
	}


	return nil
}

func ProcessExecution(execs []transaction.Transaction, accountLots []lot.Lot) error {
	for _, exec := range execs {
		if exec.IsBuyOrder() {
			newLots, err := calculateNewLotFromBuyExec(execs)
			if err != nil {
				return err
			}
			accountLots = append(accountLots, newLots...)
		}
	}

	return nil
}

func calculateNewLotFromBuyExec(exec []orderexec.Execution) ([]lot.Lot, error) {
	// newLot := lot.NewLot()



	return newLots, nil
}