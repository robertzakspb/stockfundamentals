package lot

import (
	"testing"

	"github.com/compoundinvest/stockfundamentals/internal/test"
)

func Test_FindLotsByFigi(t *testing.T) {
	targetFigi := "figi1"
	lots := []Lot{
		{
			Figi: targetFigi,
		},
		{
			Figi: "figi5",
		},
		{
			Figi: targetFigi,
		},
		{
			Figi: "figi2",
		},
		{
			Figi: "figi3",
		},
		{
			Figi: targetFigi,
		},
	}

	filteredLots := FindLotsByFigi(lots, targetFigi)

	test.AssertEqual(t, 3, len(filteredLots))
	test.AssertEqual(t, targetFigi, filteredLots[0].Figi)
	test.AssertEqual(t, targetFigi, filteredLots[1].Figi)
	test.AssertEqual(t, targetFigi, filteredLots[2].Figi)
}
