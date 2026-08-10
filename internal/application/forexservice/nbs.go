package forexservice

import (
	"strings"
	"time"

	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
)

func makeNbsUrl(date time.Time) string {
	var sb strings.Builder
	sb.WriteString("https://webappcenter.nbs.rs/ExchangeRateWebApp/ExchangeRate/IndexByDate?isSearchExecuted=true&Date=")

	sb.WriteString(timehelpers.DateInDDMMYYYFormat(date))

	sb.WriteString(".&ExchangeRateListTypeID=3")

	return sb.String()
}
