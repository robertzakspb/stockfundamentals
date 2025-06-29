package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/compoundinvest/stockfundamentals/dataseed"
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/dividend"
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/financials"
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	// "github.com/compoundinvest/stockfundamentals/features/portfolio"
)

func main() {
	// portfolio.GeMyPortfolio().PrintAllPositions()
	dataseed.InitialSeed()
	fetchExternalData()
}

func fetchExternalData() {
	security.FetchAndSaveSecurities()
	dividend.FetchAndSaveAllDividends()
	// timeseries.FetchAndSaveHistoricalQuotes()
}

func convertStockYqltoCSV() {
	//Starting from securities
	securities, _ := security.FetchSecuritiesFromDB()
	serbianStocks := []security.Stock{}
	for _, security := range securities {
		if security.Country != "RS" {
			continue
		}
		serbianStocks = append(serbianStocks, security)
	}

	file, _ := os.Create("securitiesSeed.csv")
	file.WriteString("id|company_name|is_public|isin|security_type|country_iso2|ticker|issue_size|sector\n")
	for _, security := range serbianStocks {
		const separator = "|"
		file.WriteString(security.Id.String() + separator)
		file.WriteString(security.CompanyName + separator)
		file.WriteString(strconv.FormatBool(security.IsPublic) + separator)
		file.WriteString(security.Isin + separator)
		file.WriteString(string(security.SecurityType) + separator)
		file.WriteString(security.Country + separator)
		file.WriteString(security.Ticker + separator)
		file.WriteString(strconv.Itoa(security.IssueSize) + separator)
		file.WriteString(security.Sector)
		file.WriteString("\n")
	}
	file.Close()
	//TODO: Keep going!
}

func convertdividendYqlToCsv() {
	securities, _ := security.FetchSecuritiesFromDB()
	serbianStocks := []security.Stock{}
	for _, security := range securities {
		if security.Country != "RS" {
			continue
		}
		serbianStocks = append(serbianStocks, security)
	}

	dividends, err := dividend.GetAllDividends()
	if err != nil {
		panic("No dividends here!")
	}

	serbianDividends := []dividend.Dividend{}
	for _, div := range dividends {
		for _, stock := range serbianStocks {
			if div.StockID == stock.Id {
				serbianDividends = append(serbianDividends, div)
			} else {
				fmt.Println("Not equal: ", div.StockID.String(), stock.Id.String())
			}
		}
	}

	file, _ := os.Create("dividendseeddd.csv")
	file.WriteString("id|stock_id|actual_DPS|expected_DPS|currency|record_date|payout_date|payment_period|management_comment\n")

	for _, dividend := range serbianDividends {
		const separator = "|"
		file.WriteString(dividend.Id.String() + separator)
		file.WriteString(dividend.StockID.String() + separator)
		file.WriteString(strconv.FormatFloat(dividend.ActualDPS, 'f', -1, 64) + separator)
		file.WriteString(strconv.FormatFloat(dividend.ExpectedDPS, 'f', -1, 64) + separator)
		file.WriteString(dividend.Currency + separator)
		file.WriteString(dividend.RecordDate.Format("2006-01-02") + separator)
		file.WriteString(dividend.PayoutDate.Format("2006-01-02") + separator)
		file.WriteString(dividend.PaymentPeriod + separator)
		file.WriteString(dividend.ManagementComment)
		file.WriteString("\n")
	}
}

func convertStockAndRevenueYqlToCSV() {
	metrics, err := financials.FetchFinancialsFromDB()
	if err != nil {
		fmt.Println(err)
	}

	file, _ := os.Create("revenue-income-seed.csv")
	file.WriteString("id|stock_id|metric|reporting_period|year|metric_value|metric_currency\n")

	//TODO: CSV generation logic goes here, yay!
	for _, metric := range metrics {
		const separator = "|"
		file.WriteString(metric.Id.String() + separator)
		file.WriteString(metric.StockId.String() + separator)
		file.WriteString(metric.Name + separator)
		file.WriteString(string(metric.Period) + separator)
		file.WriteString(strconv.Itoa(metric.Year) + separator)
		file.WriteString(strconv.Itoa(metric.Value) + separator)
		file.WriteString(metric.Currency)
		file.WriteString("\n")
	}
}
