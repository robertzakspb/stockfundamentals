package dateutils

import (
	"fmt"
	"time"

	"github.com/compoundinvest/stockfundamentals/Utilities/stringutils"
)

// Ancillary struct used to decode ISO-formatted datetimes returned in the responses of various API
type ISODate struct {
	time.Time
}

// Converts dates in the ISO format (e.g. "2009-09-26") into instances of time.Time
func ParseISODate(date []byte) (time.Time, error) {
	dateAsString := string(date)
	//Removing the quotation mark at the beginning of the date
	dateAsString = stringutils.TrimFirstCharacter(dateAsString)
	//Removing the quotation mark at the end of the date
	dateAsString = stringutils.TrimLastCharacter(dateAsString)

	parseDate, err := time.Parse("2006-01-02", dateAsString)
	if err != nil {
		fmt.Println(err)
		return time.Time{}, err
	}

	return parseDate, nil
}
