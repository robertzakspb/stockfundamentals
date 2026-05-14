package financialsservice

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/fundamentals/financials"
	dbfinancials "github.com/compoundinvest/stockfundamentals/internal/infrastructure/db/fundamentals/financials"
	"github.com/compoundinvest/stockfundamentals/internal/test"
	"github.com/google/uuid"
)

func Test_MapFinancialMetricsModelToDbModels(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	stockId1, stockId2 := uuid.New(), uuid.New()
	name1, name2 := "testName1", "testName2"

	metrics := []financials.FinancialMetric{
		{
			Id:              id1,
			StockId:         stockId1,
			Name:            name1,
			ReportingPeriod: financials.H2,
			Year:            2025,
			Value:           1_000_000,
			Currency:        "RUB",
		},
		{
			Id:              id2,
			StockId:         stockId2,
			Name:            name2,
			ReportingPeriod: financials.Q1,
			Year:            2026,
			Value:           2_000_000,
			Currency:        "USD",
		},
	}

	dbModels := MapFinancialMetricsModelToDbModels(metrics)

	test.AssertEqual(t, 2, len(dbModels))

	test.AssertEqual(t, id1, dbModels[0].Id)
	test.AssertEqual(t, stockId1, dbModels[0].StockId)
	test.AssertEqual(t, name1, dbModels[0].Name)
	test.AssertEqual(t, string(financials.H2), dbModels[0].ReportingPeriod)
	test.AssertEqual(t, 2025, dbModels[0].Year)
	test.AssertEqual(t, 1_000_000, dbModels[0].Value)
	test.AssertEqual(t, "RUB", dbModels[0].Currency)

	test.AssertEqual(t, id2, dbModels[1].Id)
	test.AssertEqual(t, stockId2, dbModels[1].StockId)
	test.AssertEqual(t, name2, dbModels[1].Name)
	test.AssertEqual(t, string(financials.Q1), dbModels[1].ReportingPeriod)
	test.AssertEqual(t, 2026, dbModels[1].Year)
	test.AssertEqual(t, 2_000_000, dbModels[1].Value)
	test.AssertEqual(t, "USD", dbModels[1].Currency)
}

func Test_mapYdbMetricsToMetrics_Basic(t *testing.T) {
	id1, id2 := uuid.New(), uuid.New()
	stockId1, stockId2 := uuid.New(), uuid.New()
	name1, name2 := "testName1", "testName2"

	dbMetrics := []dbfinancials.FinancialMetricDbModel{
		{
			Id:              id1,
			StockId:         stockId1,
			Name:            name1,
			ReportingPeriod: "H2",
			Year:            2025,
			Value:           1_000_000,
			Currency:        "RUB",
		},
		{
			Id:              id2,
			StockId:         stockId2,
			Name:            name2,
			ReportingPeriod: "Q1",
			Year:            2026,
			Value:           2_000_000,
			Currency:        "USD",
		},
	}

	mappedMetrics := mapYdbMetricsToMetrics(dbMetrics)

	test.AssertEqual(t, 2, len(mappedMetrics))

	test.AssertEqual(t, id1, mappedMetrics[0].Id)
	test.AssertEqual(t, stockId1, mappedMetrics[0].StockId)
	test.AssertEqual(t, name1, mappedMetrics[0].Name)
	test.AssertEqual(t, financials.H2, mappedMetrics[0].ReportingPeriod)
	test.AssertEqual(t, 2025, mappedMetrics[0].Year)
	test.AssertEqual(t, 1_000_000, mappedMetrics[0].Value)
	test.AssertEqual(t, "RUB", mappedMetrics[0].Currency)

	test.AssertEqual(t, id2, mappedMetrics[1].Id)
	test.AssertEqual(t, stockId2, mappedMetrics[1].StockId)
	test.AssertEqual(t, name2, mappedMetrics[1].Name)
	test.AssertEqual(t, financials.Q1, mappedMetrics[1].ReportingPeriod)
	test.AssertEqual(t, 2026, mappedMetrics[1].Year)
	test.AssertEqual(t, 2_000_000, mappedMetrics[1].Value)
	test.AssertEqual(t, "USD", mappedMetrics[1].Currency)
}
