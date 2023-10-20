// Package compoundinterest contains methods for calculation of:
// 1. Compound growth rate of capital over a defined period at a defined growth rate
// 2. Calculating the compound annualized growth rate of a portfolio
package compoundinterest

import "fmt"

// Compounds the specified capital during the specified number of years at the specified rate, considering the annual contributions
func Compound(capital float64, numberOfYears int, annualGrowthRate float64, annualContribution float64) float64 {
	currency := "₽"

	fmt.Printf("Capital: %v%.3f. Annual growth rate: %v\n", currency, capital, annualGrowthRate)
	for i := 0; i < numberOfYears; i++ {
		portfolioGrowth := capital * annualGrowthRate / 100
		capital += portfolioGrowth
		capital += annualContribution
		fmt.Printf("Year %v. Capital: %v%.3f млн. Growth: %.0f Contribution: %.0f. Ratio: %.1f\n", i+1, currency, capital/1000000, portfolioGrowth, annualContribution, portfolioGrowth/annualContribution)
		annualContribution *= 1.1
	}
	fmt.Println("---------------------------------------------")

	return capital
}

// Compares long-term portoflio returns depending on the provided annual growth rates (e.g. 10%, 20%, 30%, etc.)
func CompareDifferentGrowthRates(rates []float64, annualContribution float64, numberOfYears int) []PortfolioReturn {
	returns := []PortfolioReturn{}
	for i := 0; i < len(rates); i++ {
		portfolioReturn := Compound(1_000_000, numberOfYears, rates[i], annualContribution)
		returns = append(returns, PortfolioReturn{
			portfolioReturn,
			numberOfYears,
			0, //TODO: - Fix this
		})
	}

	return returns
}

func NumberOfYearsToReachTargetSum(capital float64, targetSum float64, annualGrowthRate float64, annualContribution float64) float64 {
	numberOfYears := 0.0
	fmt.Printf("Starting capital: %.1f\n", capital)
	for capital < targetSum {
		annualGrowth := capital * annualGrowthRate / 100
		capital += annualContribution
		capital += annualGrowth
		numberOfYears += 1
		fmt.Printf("Year %v. Capital: %.0f. Growth: %.0f Contribution: %.0f. Ratio: %.2f\n", numberOfYears, capital, annualGrowth, annualContribution, annualGrowth/annualContribution)
	}

	fmt.Println("It took", numberOfYears, "years to reach the target sum of", targetSum)
	fmt.Println("---------------------------------------------")
	return numberOfYears
}
