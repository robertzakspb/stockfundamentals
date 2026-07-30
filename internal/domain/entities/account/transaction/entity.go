package transaction

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id           uuid.UUID
	AccountId    uuid.UUID
	Figi         string
	Type         Type
	Timestamp    time.Time
	Side         OrderSide
	Quantity     float64
	PricePerUnit float64
	Currency     string
	Description  string
}

type OrderSide string

const (
	Buy  OrderSide = "BUY"
	Sell OrderSide = "SELL"
	None OrderSide = "NONE"
)

var OrderSideStringValue = map[OrderSide]string{
	Buy:  "BUY",
	Sell: "SELL",
	None: "NONE",
}
var OrderSideLookup = map[string]OrderSide{
	"BUY":  Buy,
	"SELL": Sell,
	"NONE": None,
}

type Type string

const (
	OrderExecution        Type = "ORDER_EXECUTION"
	Deposit               Type = "DEPOSIT"
	Withdrawal            Type = "WITHDRAWAL"
	CouponTax             Type = "COUPON_TAX"
	Tax                   Type = "TAX"
	DividendTax           Type = "DIVIDEND_TAX"
	BondRedemption        Type = "BOND_REDEMPTION"
	PartialBondRedemption Type = "PARTIAL_BOND_REDEMPTION"
	BrokerFee             Type = "BROKER_FEE"
	Dividend              Type = "DIVIDEND"
	CouponPayment         Type = "COUPON_PAYMENT"
	Unknown               Type = "UNKNOWN"
)

var TypeLookup = map[string]Type{
	"ORDER_EXECUTION":         OrderExecution,
	"DEPOSIT":                 Deposit,
	"WITHDRAWAL":              Withdrawal,
	"COUPON_TAX":              CouponTax,
	"TAX":                     Tax,
	"DIVIDEND_TAX":            DividendTax,
	"UNKNOWN":                 Unknown,
	"BOND_REDEMPTION":         BondRedemption,
	"BROKER_FEE":              BrokerFee,
	"DIVIDEND":                Dividend,
	"COUPON_PAYMENT":          CouponPayment,
	"PARTIAL_BOND_REDEMPTION": PartialBondRedemption,
}

func New(accountId uuid.UUID, figi string, timestamp time.Time, quantity, price float64, description, side, transactionType, currency string) (Transaction, error) {
	if quantity < 0 {
		return Transaction{}, fmt.Errorf("quantity is smaller than 0: %f", quantity)
	}
	if price < 0 {
		return Transaction{}, fmt.Errorf("price is smaller than 0: %f", price)
	}
	if len(description) > 1000 {
		return Transaction{}, fmt.Errorf("description cannot contain more than 1000 characters")
	}
	if currency == "" {
		return Transaction{}, fmt.Errorf("Invalid currency")
	}
	orderSide, found := OrderSideLookup[strings.ToUpper(side)]
	if !found {
		return Transaction{}, errors.New("Invalid order side: " + side)
	}
	parsedType, found := TypeLookup[strings.ToUpper(transactionType)]
	if !found {
		return Transaction{}, errors.New("Invalid transaction type: " + transactionType)
	}

	return Transaction{
		Id:           uuid.New(),
		AccountId:    accountId,
		Figi:         figi,
		Timestamp:    timestamp,
		Quantity:     quantity,
		PricePerUnit: price,
		Side:         orderSide,
		Description:  description,
		Type:         parsedType,
	}, nil
}
