package dividend

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/security"
	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_CumulativeReturn(t *testing.T) {
	forecast1 := DividendForecast{
		Yield: 0.1,
	}
	forecast2 := DividendForecast{
		Yield: 0.2,
	}
	securityDivForecast := SecurityDivForecast{
		Forecasts: []DividendForecast{forecast1, forecast2},
	}
	expectedYield := 0.32

	actualYield := securityDivForecast.CumulativeReturn()

	test.AssertEqualFloat(t, expectedYield, actualYield, 0.001)
}

func Test_GroupForecastsBySecurity(t *testing.T) {
	forecast1 := DividendForecast{
		Stock: security.Stock{
			Figi: "ABC",
		},
		Yield: 0.1,
	}
	forecast2 := DividendForecast{
		Stock: security.Stock{
			Figi: "ABC",
		},
		Yield: 0.2,
	}
	forecast3 := DividendForecast{
		Stock: security.Stock{
			Figi: "DDD",
		},
		Yield: 0.3,
	}
	forecast4 := DividendForecast{
		Stock: security.Stock{
			Figi: "DDD",
		},
		Yield: 0.4,
	}
	forecast5 := DividendForecast{
		Stock: security.Stock{
			Figi: "XYZ",
		},
		Yield: 3,
	}

	groupedForecasts := GroupForecastsBySecurity([]DividendForecast{forecast1, forecast2, forecast3, forecast4, forecast5})

	test.AssertEqual(t, "XYZ", groupedForecasts[0].Figi)
	test.AssertEqual(t, "DDD", groupedForecasts[1].Figi)
	test.AssertEqual(t, "ABC", groupedForecasts[2].Figi)
}
