package dividend

import (
	"fmt"
	"time"

	"github.com/compoundinvest/stockfundamentals/Utilities/dateutils"
)

type MoexDividendDTO []struct {
	Dividends []struct {
		Secid             string   `json:"secid,omitempty"`
		Isin              string   `json:"isin"`
		Registryclosedate MoexDate `json:"registryclosedate"`
		Value             float64  `json:"value"`
		Currencyid        string   `json:"currencyid"`
	} `json:"dividends,omitempty"`
}

func (dividend MoexDividendDTO) asDividends() []Dividend {

	var parsedDividends []Dividend
	for _, dividendItem := range dividend[1].Dividends {
		newDividend := Dividend{
			Isin:       dividendItem.Isin,
			Ticker:     dividendItem.Secid,
			AmountPaid: dividendItem.Value,
			Currency:   dividendItem.Currencyid,
			Date:       dividendItem.Registryclosedate.Time,
		}
		parsedDividends = append(parsedDividends, newDividend)
	}

	return parsedDividends
}

// Ancillary struct used to decode datetimes returned in the responses of MOEX API
type MoexDate struct {
	time.Time //MOEX API supplies dates in the ISO format: "2009-09-26"
}

func (t *MoexDate) UnmarshalJSON(bytes []byte) error {
	date, err := dateutils.ParseISODate(bytes)
	if err != nil {
		fmt.Println(err)
	}

	t.Time = date
	return nil
}
