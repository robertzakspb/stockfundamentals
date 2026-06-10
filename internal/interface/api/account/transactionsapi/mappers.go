package transactionsapi

import (
	"errors"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
)

func mapTransactionsToDto(transactions []transaction.Transaction) []TransactionDto {
	dtos := make([]TransactionDto, len(transactions))

	for i := range transactions {
		dto := TransactionDto{
			Id:           transactions[i].Id,
			AccountId:    transactions[i].AccountId,
			Figi:         transactions[i].Figi,
			Type:         string(transactions[i].Type),
			Timestamp:    transactions[i].Timestamp,
			Side:         string(transactions[i].Side),
			Quantity:     transactions[i].Quantity,
			PricePerUnit: transactions[i].PricePerUnit,
			Currency:     transactions[i].Currency,
			Description:  transactions[i].Description,
		}
		dtos[i] = dto
	}

	return dtos
}

func mapTransactionDtosToTransactions(dtos []TransactionDto) ([]transaction.Transaction, error) {
	transactions := make([]transaction.Transaction, len(dtos))

	for i := range dtos {
		tType, found := transaction.TypeLookup[dtos[i].Type]
		if !found {
			return transactions, errors.New("Unsupported transaction type: " + dtos[i].Type)
		}
		tSide, found := transaction.OrderSideLookup[dtos[i].Side]
		if !found && tType == transaction.OrderExecution {
			return transactions, errors.New("Unsupported order side: " + dtos[i].Side)
		}
		transaction := transaction.Transaction{
			Id:           dtos[i].Id,
			AccountId:    dtos[i].AccountId,
			Figi:         dtos[i].Figi,
			Type:         tType,
			Timestamp:    dtos[i].Timestamp,
			Side:         tSide,
			Quantity:     dtos[i].Quantity,
			PricePerUnit: dtos[i].PricePerUnit,
			Currency:     dtos[i].Currency,
			Description:  dtos[i].Description,
		}
		transactions[i] = transaction
	}

	return transactions, nil
}
