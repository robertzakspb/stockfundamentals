package transactionprocessor

import (
	"testing"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/account"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/account/transaction"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
	"github.com/google/uuid"
)

func Test_Recalculate_NilAccountId(t *testing.T) {
	account := account.Account{}
	lots := []lot.Lot{{}}
	transactions := []transaction.Transaction{{}}

	_, _, _, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertError(t, err)
}

func Test_Recalculate_NoTransactions(t *testing.T) {
	account := account.Account{Id: uuid.New()}
	lots := []lot.Lot{{}}
	transactions := []transaction.Transaction{}

	_, _, _, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertError(t, err)
}

func Test_Recalculate_NoLots(t *testing.T) {
	account := account.Account{Id: uuid.New()}
	lots := []lot.Lot{}
	transactions := []transaction.Transaction{{}}

	_, _, _, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertError(t, err)
}

func Test_Recalculate_SingleBuy_NegativeCashBalance(t *testing.T) {
	account := account.Account{Id: uuid.New(), CashBalance: 50}
	lots := []lot.Lot{{Figi: "figi1", Quantity: 5}}
	transactions := []transaction.Transaction{
		{Figi: "figi2", Timestamp: time.Now()},
		{Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 10, PricePerUnit: 25},
	}

	_, _, _, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertError(t, err)
}

func Test_Recalculate_SingleBuy_UnsupportedTransactionType(t *testing.T) {
	tranId1, tranId2 := uuid.New(), uuid.New()
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{{Figi: "figi1", Quantity: 5}}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi2", Timestamp: time.Now()},
		{Id: tranId2, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 10, PricePerUnit: 25},
	}

	account, lots, _, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertError(t, err)
}

func Test_Recalculate_SingleBuy_PositiveCase(t *testing.T) {
	tranId1 := uuid.New()
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{{Figi: "figi1", Quantity: 5, IsClosed: false}}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 10, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Buy, Currency: "USD"},
	}

	account, lots, relations, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertNoError(t, err)

	test.AssertEqual(t, 2, len(lots))
	test.AssertEqual(t, "figi1", lots[0].Figi)
	test.AssertEqual(t, 5, lots[0].Quantity)
	test.AssertFalse(t, lots[0].IsClosed)
	test.AssertEqual(t, "figi1", lots[1].Figi)
	test.AssertEqual(t, 10, lots[1].Quantity)
	test.AssertFalse(t, lots[1].IsClosed)

	test.AssertEqual(t, 150, account.CashBalance)

	test.AssertEqual(t, 1, len(relations))
	test.AssertEqual(t, tranId1, relations[0].TransactionId)
	test.AssertEqual(t, uuid.Nil, relations[0].BondLotId)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[0].Date))
	test.AssertEqual(t, 10, relations[0].Quantity)
}

func Test_Recalculate_MultipleBuys(t *testing.T) {
	tranId1, tranId2, tranId3 := uuid.New(), uuid.New(), uuid.New()

	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{{Figi: "figi1", Quantity: 5}}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 10, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Buy, Currency: "USD"},
		{Id: tranId2, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -6), Quantity: 2, PricePerUnit: 10, Type: transaction.OrderExecution, Side: transaction.Buy, Currency: "USD"},
		{Id: tranId3, Figi: "figi2", Timestamp: time.Now().AddDate(0, 0, -5), Quantity: 1, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Buy, Currency: "USD"},
	}

	account, lots, relations, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertNoError(t, err)

	test.AssertEqual(t, 4, len(lots))
	test.AssertEqual(t, "figi1", lots[0].Figi)
	test.AssertEqual(t, 5, lots[0].Quantity)
	test.AssertFalse(t, lots[0].IsClosed)
	test.AssertEqual(t, "figi1", lots[1].Figi)
	test.AssertEqual(t, 10, lots[1].Quantity)
	test.AssertFalse(t, lots[1].IsClosed)
	test.AssertEqual(t, "figi1", lots[2].Figi)
	test.AssertEqual(t, 2, lots[2].Quantity)
	test.AssertFalse(t, lots[2].IsClosed)
	test.AssertEqual(t, "figi2", lots[3].Figi)
	test.AssertEqual(t, 1, lots[3].Quantity)
	test.AssertFalse(t, lots[3].IsClosed)

	test.AssertEqual(t, 105, account.CashBalance)

	test.AssertEqual(t, 3, len(relations))
	test.AssertEqual(t, tranId1, relations[0].TransactionId)
	test.AssertEqual(t, tranId2, relations[1].TransactionId)
	test.AssertEqual(t, tranId3, relations[2].TransactionId)
	test.AssertEqual(t, uuid.Nil, relations[0].BondLotId)
	test.AssertEqual(t, uuid.Nil, relations[1].BondLotId)
	test.AssertEqual(t, uuid.Nil, relations[2].BondLotId)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[0].Date))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -6), relations[1].Date))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -5), relations[2].Date))
	test.AssertEqual(t, 10, relations[0].Quantity)
	test.AssertEqual(t, 2, relations[1].Quantity)
	test.AssertEqual(t, 1, relations[2].Quantity)
}

func Test_Recalculate_SingleSell_Negative_ExceededAvailableQuantity(t *testing.T) {
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{{Figi: "figi1", Quantity: 5}}
	transactions := []transaction.Transaction{
		{Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 10, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
	}

	account, lots, _, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertError(t, err)
}

func Test_Recalculate_SingleSell_OneLotIsClosed(t *testing.T) {
	tranId1, lotId1 := uuid.New(), uuid.New()
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{
		{Id: lotId1, Figi: "figi1", Quantity: 5, IsClosed: false},
		{Figi: "figi1", Quantity: 10, IsClosed: false},
	}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 5, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
	}

	account, lots, relations, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 2, len(lots))
	test.AssertEqual(t, "figi1", lots[0].Figi)
	test.AssertEqual(t, 0, lots[0].Quantity)
	test.AssertTrue(t, lots[0].IsClosed)
	test.AssertEqual(t, "figi1", lots[1].Figi)
	test.AssertEqual(t, 10, lots[1].Quantity)
	test.AssertFalse(t, lots[1].IsClosed)

	test.AssertEqual(t, 525, account.CashBalance)

	test.AssertEqual(t, 1, len(relations))
	test.AssertEqual(t, tranId1, relations[0].TransactionId)
	test.AssertEqual(t, lotId1, relations[0].StockLotId)
	test.AssertEqual(t, uuid.Nil, relations[0].BondLotId)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[0].Date))
	test.AssertEqual(t, -5, relations[0].Quantity)
}

func Test_Recalculate_SingleSell_MultipleLotsAreClosed(t *testing.T) {
	tranId1, lotId1, lotId2 := uuid.New(), uuid.New(), uuid.New()
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{
		{Id: lotId1, Figi: "figi1", Quantity: 5, IsClosed: false},
		{Id: lotId2, Figi: "figi1", Quantity: 10, IsClosed: false},
	}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 15, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
	}

	account, lots, relations, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 2, len(lots))
	test.AssertEqual(t, "figi1", lots[0].Figi)
	test.AssertEqual(t, 0, lots[0].Quantity)
	test.AssertTrue(t, lots[0].IsClosed)
	test.AssertEqual(t, "figi1", lots[1].Figi)
	test.AssertEqual(t, 0, lots[1].Quantity)
	test.AssertTrue(t, lots[1].IsClosed)

	test.AssertEqual(t, 775, account.CashBalance)

	test.AssertEqual(t, 2, len(relations))
	test.AssertEqual(t, tranId1, relations[0].TransactionId)
	test.AssertEqual(t, tranId1, relations[1].TransactionId)
	test.AssertEqual(t, lotId1, relations[0].StockLotId)
	test.AssertEqual(t, lotId2, relations[1].StockLotId)
	test.AssertEqual(t, uuid.Nil, relations[0].BondLotId)
	test.AssertEqual(t, uuid.Nil, relations[1].BondLotId)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[0].Date))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[1].Date))
	test.AssertEqual(t, -5, relations[0].Quantity)
	test.AssertEqual(t, -10, relations[1].Quantity)
}

func Test_Recalculate_SingleSell_OneLostIsClosedAndAnotherPartiallyClosed(t *testing.T) {
	tranId1, lotId1, lotId2 := uuid.New(), uuid.New(), uuid.New()
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{
		{Id: lotId1, Figi: "figi1", Quantity: 5, IsClosed: false},
		{Id: lotId2, Figi: "figi1", Quantity: 10, IsClosed: false},
	}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 7, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
	}

	account, lots, relations, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertNoError(t, err)
	test.AssertEqual(t, 2, len(lots))
	test.AssertEqual(t, "figi1", lots[0].Figi)
	test.AssertEqual(t, 0, lots[0].Quantity)
	test.AssertTrue(t, lots[0].IsClosed)
	test.AssertEqual(t, "figi1", lots[1].Figi)
	test.AssertEqual(t, 8, lots[1].Quantity)
	test.AssertFalse(t, lots[1].IsClosed)

	test.AssertEqual(t, 575, account.CashBalance)

	test.AssertEqual(t, 2, len(relations))
	test.AssertEqual(t, tranId1, relations[0].TransactionId)
	test.AssertEqual(t, tranId1, relations[1].TransactionId)
	test.AssertEqual(t, lotId1, relations[0].StockLotId)
	test.AssertEqual(t, lotId2, relations[1].StockLotId)
	test.AssertEqual(t, uuid.Nil, relations[0].BondLotId)
	test.AssertEqual(t, uuid.Nil, relations[1].BondLotId)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[0].Date))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[1].Date))
	test.AssertEqual(t, -5, relations[0].Quantity)
	test.AssertEqual(t, -2, relations[1].Quantity)
}

func Test_Recalculate_MultipleSells_MultipleLotsAreClosed(t *testing.T) {
	tranId1, tranId2, lotId1, lotId2 := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{
		{Id: lotId1, Figi: "figi1", Quantity: 5, IsClosed: false},
		{Id: lotId2, Figi: "figi1", Quantity: 10, IsClosed: false},
	}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 5, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
		{Id: tranId2, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -6), Quantity: 10, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
	}

	account, lots, relations, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertNoError(t, err)

	test.AssertEqual(t, 2, len(lots))
	test.AssertEqual(t, "figi1", lots[0].Figi)
	test.AssertEqual(t, 0, lots[0].Quantity)
	test.AssertTrue(t, lots[0].IsClosed)
	test.AssertEqual(t, "figi1", lots[1].Figi)
	test.AssertEqual(t, 0, lots[1].Quantity)
	test.AssertTrue(t, lots[1].IsClosed)

	test.AssertEqual(t, 775, account.CashBalance)

	test.AssertEqual(t, 2, len(relations))
	test.AssertEqual(t, tranId1, relations[0].TransactionId)
	test.AssertEqual(t, tranId2, relations[1].TransactionId)
	test.AssertEqual(t, lotId1, relations[0].StockLotId)
	test.AssertEqual(t, lotId2, relations[1].StockLotId)
	test.AssertEqual(t, uuid.Nil, relations[0].BondLotId)
	test.AssertEqual(t, uuid.Nil, relations[1].BondLotId)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[0].Date))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -6), relations[1].Date))
	test.AssertEqual(t, -5, relations[0].Quantity)
	test.AssertEqual(t, -10, relations[1].Quantity)
}

func Test_Recalculate_SingleSell_OneFullCloseAndOnePartialClose(t *testing.T) {
	tranId1, tranId2, lotId1, lotId2 := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	account := account.Account{Id: uuid.New(), CashBalance: 400}
	lots := []lot.Lot{
		{Id: lotId1, Figi: "figi1", Quantity: 5, IsClosed: false},
		{Id: lotId2, Figi: "figi1", Quantity: 10, IsClosed: false},
	}
	transactions := []transaction.Transaction{
		{Id: tranId1, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -7), Quantity: 5, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
		{Id: tranId2, Figi: "figi1", Timestamp: time.Now().AddDate(0, 0, -6), Quantity: 3, PricePerUnit: 25, Type: transaction.OrderExecution, Side: transaction.Sell, Currency: "USD"},
	}

	account, lots, relations, err := recalculateLotsAndCashBalances(account, transactions, lots)

	test.AssertNoError(t, err)

	test.AssertEqual(t, 2, len(lots))
	test.AssertEqual(t, "figi1", lots[0].Figi)
	test.AssertEqual(t, 0, lots[0].Quantity)
	test.AssertTrue(t, lots[0].IsClosed)
	test.AssertEqual(t, "figi1", lots[1].Figi)
	test.AssertEqual(t, 7, lots[1].Quantity)
	test.AssertFalse(t, lots[1].IsClosed)

	test.AssertEqual(t, 600, account.CashBalance)

	test.AssertEqual(t, 2, len(relations))
	test.AssertEqual(t, tranId1, relations[0].TransactionId)
	test.AssertEqual(t, tranId2, relations[1].TransactionId)
	test.AssertEqual(t, lotId1, relations[0].StockLotId)
	test.AssertEqual(t, lotId2, relations[1].StockLotId)
	test.AssertEqual(t, uuid.Nil, relations[0].BondLotId)
	test.AssertEqual(t, uuid.Nil, relations[1].BondLotId)
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -7), relations[0].Date))
	test.AssertTrue(t, timehelpers.AreEqualDates(time.Now().AddDate(0, 0, -6), relations[1].Date))
	test.AssertEqual(t, -5, relations[0].Quantity)
	test.AssertEqual(t, -3, relations[1].Quantity)
}
