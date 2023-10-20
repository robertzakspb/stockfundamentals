package fundamentals

import (
	"fmt"
	"github.com/compoundinvest/stockfundamentals/Utilities/dateutils"
	"time"
)

type PolygonFundamentalDataDTO struct {
	Ticker  string
	Results []struct {
		CompanyName      string      `json:"company_name"`
		ReportStartDate  PolygonDate `json:"start_date"`
		ReportEndDate    PolygonDate `json:"end_date"`
		ReportFilingDate PolygonDate `json:"filing_date"`
		Financials       struct {
			IncomeStatement struct {
				Revenue    FinancialMetric `json:"revenues"`
				DilutedEPS FinancialMetric `json:"diluted_earnings_per_share"`
			} `json:"income_statement"`
		} `json:"financials"`
	} `json:"results"`
}

// Converts an instance of PolygonFundamentalDataDTO to an instance of StockFundamentals
func (apiResponse PolygonFundamentalDataDTO) AsStockFundamentals() (fundamendals StockFundamentals) {
	eps := StockFinancialResult{Metric: EPS}
	revenue := StockFinancialResult{Metric: Revenue}

	for i := 0; i < len(apiResponse.Results); i++ {
		revenue.Values = append(revenue.Values, apiResponse.Results[i].Financials.IncomeStatement.Revenue)
		revenue.Values[i].Year = apiResponse.Results[i].ReportEndDate.Time

		eps.Values = append(eps.Values, apiResponse.Results[i].Financials.IncomeStatement.DilutedEPS)
		eps.Values[i].Year = apiResponse.Results[i].ReportEndDate.Time
	}

	fundamentals := StockFundamentals{
		Ticker:      apiResponse.Ticker,
		CompanyName: apiResponse.Results[0].CompanyName,
		Financials:  []StockFinancialResult{eps, revenue},
	}

	return fundamentals
}

// Ancillary struct used to decode datetimes returned in the responses of Polygon API
type PolygonDate struct {
	time.Time //Polygon API supplies dates in the ISO format: "2009-09-26"
}

func (t *PolygonDate) UnmarshalJSON(bytes []byte) error {
	date, err := dateutils.ParseISODate(bytes)
	if err != nil {
		fmt.Println(err)
	}

	t.Time = date
	return nil
}

type FinancialMetric struct {
	Label string  `json:"label"`
	Value float64 `json:"value,omitempty"`
	Unit  string  `json:"unit"`
	Year  time.Time
}

//FULL DTO:

// type PolygonFundamentalDataDTO struct {
// 	Results []struct {
// 		Financials struct {
// 			// CashFlowStatement struct {
// 			// 	NetCashFlowContinuing                        PolygonFinancialMetricDTO `json:"net_cash_flow_continuing"`
// 			// 	NetCashFlowFromFinancingActivities           PolygonFinancialMetricDTO `json:"net_cash_flow_from_financing_activities"`
// 			// 	NetCashFlowFromOperatingActivitiesContinuing PolygonFinancialMetricDTO `json:"net_cash_flow_from_operating_activities_continuing"`
// 			// 	NetCashFlowFromInvestingActivitiesContinuing PolygonFinancialMetricDTO `json:"net_cash_flow_from_investing_activities_continuing"`
// 			// 	NetCashFlowFromInvestingActivities           PolygonFinancialMetricDTO `json:"net_cash_flow_from_investing_activities"`
// 			// 	NetCashFlowFromOperatingActivities           PolygonFinancialMetricDTO `json:"net_cash_flow_from_operating_activities"`
// 			// 	NetCashFlow                                  PolygonFinancialMetricDTO `json:"net_cash_flow"`
// 			// 	NetCashFlowFromFinancingActivitiesContinuing PolygonFinancialMetricDTO `json:"net_cash_flow_from_financing_activities_continuing"`
// 			// } `json:"cash_flow_statement"`
// 			// ComprehensiveIncome struct {
// 			// 	ComprehensiveIncomeLossAttributableToParent                 PolygonFinancialMetricDTO `json:"comprehensive_income_loss_attributable_to_parent"`
// 			// 	ComprehensiveIncomeLoss                                     PolygonFinancialMetricDTO `json:"comprehensive_income_loss"`
// 			// 	ComprehensiveIncomeLossAttributableToNoncontrollingInterest PolygonFinancialMetricDTO `json:"comprehensive_income_loss_attributable_to_noncontrolling_interest"`
// 			// 	OtherComprehensiveIncomeLoss                                PolygonFinancialMetricDTO `json:"other_comprehensive_income_loss"`
// 			// } `json:"comprehensive_income"`
// 			IncomeStatement struct {
// 				// GrossProfit                                                         PolygonFinancialMetricDTO `json:"gross_profit"`
// 				// IncomeLossBeforeEquityMethodInvestments                             PolygonFinancialMetricDTO `json:"income_loss_before_equity_method_investments"`
// 				Revenue    PolygonFinancialMetricDTO `json:"revenues"`
// 				DilutedEPS PolygonFinancialMetricDTO `json:"diluted_earnings_per_share"`
// 				// PreferredStockDividendsAndOtherAdjustments                          PolygonFinancialMetricDTO `json:"preferred_stock_dividends_and_other_adjustments"`
// 				// OperatingExpenses                                                   PolygonFinancialMetricDTO `json:"operating_expenses"`
// 				// IncomeTaxExpenseBenefitDeferred                                     PolygonFinancialMetricDTO `json:"income_tax_expense_benefit_deferred"`
// 				// OperatingIncomeLoss                                                 PolygonFinancialMetricDTO `json:"operating_income_loss"`
// 				// BenefitsCostsExpenses                                               PolygonFinancialMetricDTO `json:"benefits_costs_expenses"`
// 				// CostOfRevenue                                                       PolygonFinancialMetricDTO `json:"cost_of_revenue"`
// 				// IncomeLossFromContinuingOperationsBeforeTax                         PolygonFinancialMetricDTO `json:"income_loss_from_continuing_operations_before_tax"`
// 				// NonoperatingIncomeLoss                                              PolygonFinancialMetricDTO `json:"nonoperating_income_loss"`
// 				// NetIncomeLossAttributableToParent                                   PolygonFinancialMetricDTO `json:"net_income_loss_attributable_to_parent"`
// 				// NetIncomeLoss                                                       PolygonFinancialMetricDTO `json:"net_income_loss"`
// 				// NetIncomeLossAvailableToCommonStockholdersBasic                     PolygonFinancialMetricDTO `json:"net_income_loss_available_to_common_stockholders_basic"`
// 				// BasicEarningsPerShare                                               PolygonFinancialMetricDTO `json:"basic_earnings_per_share"`
// 				// ParticipatingSecuritiesDistributedAndUndistributedEarningsLossBasic PolygonFinancialMetricDTO `json:"participating_securities_distributed_and_undistributed_earnings_loss_basic"`
// 				// IncomeLossFromContinuingOperationsAfterTax                          PolygonFinancialMetricDTO `json:"income_loss_from_continuing_operations_after_tax"`
// 				// CostsAndExpenses                                                    PolygonFinancialMetricDTO `json:"costs_and_expenses"`
// 				// NetIncomeLossAttributableToNoncontrollingInterest                   PolygonFinancialMetricDTO `json:"net_income_loss_attributable_to_noncontrolling_interest"`
// 				// IncomeTaxExpenseBenefit                                             PolygonFinancialMetricDTO `json:"income_tax_expense_benefit"`
// 			} `json:"income_statement"`
// 			// BalanceSheet struct {
// 			// 	NoncurrentLiabilities                      PolygonFinancialMetricDTO `json:"noncurrent_liabilities"`
// 			// 	EquityAttributableToNoncontrollingInterest PolygonFinancialMetricDTO `json:"equity_attributable_to_noncontrolling_interest"`
// 			// 	Assets                                     PolygonFinancialMetricDTO `json:"assets"`
// 			// 	CurrentLiabilities                         PolygonFinancialMetricDTO `json:"current_liabilities"`
// 			// 	Equity                                     PolygonFinancialMetricDTO `json:"equity"`
// 			// 	CurrentAssets                              PolygonFinancialMetricDTO `json:"current_assets"`
// 			// 	EquityAttributableToParent                 PolygonFinancialMetricDTO `json:"equity_attributable_to_parent"`
// 			// 	LiabilitiesAndEquity                       PolygonFinancialMetricDTO `json:"liabilities_and_equity"`
// 			// 	Liabilities                                PolygonFinancialMetricDTO `json:"liabilities"`
// 			// 	NoncurrentAssets                           PolygonFinancialMetricDTO `json:"noncurrent_assets"`
// 			// } `json:"balance_sheet"`
// 		} `json:"financials"`
// 		StartDate           string `json:"start_date"`
// 		EndDate             string `json:"end_date"`
// 		FilingDate          string `json:"filing_date"`
// 		Cik                 string `json:"cik"`
// 		CompanyName         string `json:"company_name"`
// 		FiscalPeriod        string `json:"fiscal_period"`
// 		FiscalYear          string `json:"fiscal_year"`
// 		SourceFilingURL     string `json:"source_filing_url"`
// 		SourceFilingFileURL string `json:"source_filing_file_url"`
// 	} `json:"results"`
// 	Status    string `json:"status"`
// 	RequestID string `json:"request_id"`
// 	Count     int    `json:"count"`
// }
