package orderprocessor

import (
	// "github.com/compoundinvest/stockfundamentals/internal/domain/entities/orderexec"
	// "github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
)

// func ProcessExecution(execs []orderexec.Execution, accountLots []lot.Lot) error {
// 	for _, exec := range execs {
// 		if exec.IsBuyOrder() {
// 			newLots, err := calculateNewLotFromBuyExec(execs)
// 			if err != nil {
// 				return err
// 			}
// 			accountLots = append(accountLots, newLots...)
// 		}
// 	}

// 	return nil
// }

// func calculateNewLotFromBuyExec(exec []orderexec.Execution) ([]lot.Lot, error) {
// 	// newLot := lot.NewLot()



// 	return newLots, nil
// }