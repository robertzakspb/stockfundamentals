package transactionprocessor

import (
	"errors"

	"github.com/compoundinvest/stockfundamentals/internal/application/shared"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/account/transactionsdb"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	investapi "opensource.tbank.ru/invest/invest-go/proto"
)

func mapTransactionsToDbModel(transactions []transaction.Transaction) []transactionsdb.TransactionDbModel {
	dbModels := make([]transactionsdb.TransactionDbModel, len(transactions))

	for i := range transactions {
		dbModel := transactionsdb.TransactionDbModel{
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
		dbModels[i] = dbModel
	}

	return dbModels
}

func mapTransactionDbModelsToTransactions(dbModels []transactionsdb.TransactionDbModel) ([]transaction.Transaction, error) {
	transactions := make([]transaction.Transaction, len(dbModels))

	for i := range dbModels {
		tType, found := transaction.TypeLookup[dbModels[i].Type]
		if !found {
			return transactions, errors.New("Unable to map the provided value – " + dbModels[i].Type + "- to a transaction type")
		}
		side, found := transaction.OrderSideLookup[dbModels[i].Side]
		if !found {
			return transactions, errors.New("Unable to map the provided value – " + dbModels[i].Type + "- to a transaction side")
		}

		transaction := transaction.Transaction{
			Id:           dbModels[i].Id,
			AccountId:    dbModels[i].AccountId,
			Figi:         dbModels[i].Figi,
			Type:         tType,
			Timestamp:    dbModels[i].Timestamp,
			Side:         side,
			Quantity:     dbModels[i].Quantity,
			PricePerUnit: dbModels[i].PricePerUnit,
			Currency:     dbModels[i].Currency,
			Description:  dbModels[i].Description,
		}
		transactions[i] = transaction
	}

	return transactions, nil
}

func mapTBankTransactionsToTransaction(tTransactions []*investapi.OperationItem) []transaction.Transaction {
	transactions := []transaction.Transaction{}

	for i := range tTransactions {
		var tType transaction.Type
		switch tTransactions[i].Type {
		case investapi.OperationType_OPERATION_TYPE_INPUT:
			tType = transaction.Deposit
		case investapi.OperationType_OPERATION_TYPE_UNSPECIFIED:
			tType = transaction.Unknown
		case investapi.OperationType_OPERATION_TYPE_BOND_TAX:
			tType = transaction.CouponTax
		case investapi.OperationType_OPERATION_TYPE_TAX:
			tType = transaction.Tax
		case investapi.OperationType_OPERATION_TYPE_BOND_REPAYMENT_FULL:
			tType = transaction.BondRedemption
		case investapi.OperationType_OPERATION_TYPE_DIVIDEND_TAX:
			tType = transaction.DividendTax
		case investapi.OperationType_OPERATION_TYPE_OUTPUT:
			tType = transaction.Withdrawal
		case investapi.OperationType_OPERATION_TYPE_BOND_REPAYMENT:
			tType = transaction.PartialBondRedemption
		case investapi.OperationType_OPERATION_TYPE_BUY, investapi.OperationType_OPERATION_TYPE_SELL:
			tType = transaction.OrderExecution
		case investapi.OperationType_OPERATION_TYPE_BROKER_FEE:
			tType = transaction.BrokerFee
		case investapi.OperationType_OPERATION_TYPE_DIVIDEND, investapi.OperationType_OPERATION_TYPE_DIV_EXT:
			tType = transaction.Dividend
		case investapi.OperationType_OPERATION_TYPE_COUPON:
			tType = transaction.CouponPayment
		default:
			logger.Log("Unexpected transaction type from T API: "+tTransactions[i].Type.String(), logger.WARNING)
			tType = transaction.Unknown
		}

		var side transaction.OrderSide //TODO: Test this logic
		switch tTransactions[i].Type {
		case investapi.OperationType_OPERATION_TYPE_SELL:
			side = transaction.Sell
		case investapi.OperationType_OPERATION_TYPE_BUY:
			side = transaction.Buy
		default:
			side = transaction.None
		}

		price := 0.0
		switch tTransactions[i].Type {
		case investapi.OperationType_OPERATION_TYPE_BROKER_FEE, investapi.OperationType_OPERATION_TYPE_INPUT, investapi.OperationType_OPERATION_TYPE_DIVIDEND, investapi.OperationType_OPERATION_TYPE_DIV_EXT, investapi.OperationType_OPERATION_TYPE_DIVIDEND_TAX:
			price = tTransactions[i].Payment.ToFloat()
		case investapi.OperationType_OPERATION_TYPE_BUY, investapi.OperationType_OPERATION_TYPE_SELL:
			price = tTransactions[i].Price.ToFloat()
		default:
			price = tTransactions[i].Price.ToFloat()
		}

		//The description should include both the security's name (if present) and the description
		description := tTransactions[i].Description
		if tTransactions[i].Name != "" {
			description = tTransactions[i].Name + ". " + description
		}

		transaction := transaction.Transaction{
			Id:           uuid.New(),
			AccountId:    uuid.MustParse(shared.TINKOFF_IIS_ACCOUNT_ID), //Quick and dirty solution
			Figi:         tTransactions[i].Figi,
			Type:         tType,
			Timestamp:    tTransactions[i].Date.AsTime(),
			PricePerUnit: price,
			Side:         side,
			Description:  description,
			Currency:     tTransactions[i].Price.Currency,
			Quantity:     float64(tTransactions[i].Quantity),
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}
