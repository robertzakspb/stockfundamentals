package transactionsapi

import (
	"time"

	accountsapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/accounts"
	portfolioapi "github.com/compoundinvest/stockfundamentals/internal/interface/api/account/stock-portfolio"
	"github.com/google/uuid"
)

type TransactionDto struct {
	Id           uuid.UUID `json:"id"`
	AccountId    uuid.UUID `json:"accountId"`
	Figi         string    `json:"figi"`
	Type         string    `json:"type"`
	Timestamp    time.Time `json:"timestamp" sql:"timestamp"`
	Side         string    `json:"side"`
	Quantity     float64   `json:"quantity"`
	PricePerUnit float64   `json:"pricePerUnit"`
	Currency     string    `json:"currency"`
	Description  string    `json:"description"`
}

type TransactionLotRelationDto struct {
	Id            string    `json:"id"`
	TransactionId string    `json:"transactionId"`
	StockLotId    string    `json:"stockLotId"`
	BondLotId     string    `json:"bondLotId"`
	Date          time.Time `json:"date"`
	Quantity      float64   `json:"quantity"`
}

// Convenience wrapper for the four encapsulated entities
type AccountsTransactionsLotsRelationsDto struct {
	AdjustedAccounts []accountsapi.AccountDto    `json:"adjustedAccounts"`
	Transactions     []TransactionDto            `json:"transactions"`
	AdjustedLots     []portfolioapi.LotDto       `json:"adjustedLots"`
	Relations        []TransactionLotRelationDto `json:"relations"`
}
