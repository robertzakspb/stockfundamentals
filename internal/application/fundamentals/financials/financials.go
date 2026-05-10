package financials

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
)

func FetchFinancialMetrics() ([]financials.FinancialMetric, error) {
	return dbfinancials.FetchFinancialMetrics()
}
