package securityinfo

import (
	"fmt"
	"strings"
)

func SaveMoexSecuritiesToDB(securities []MoexSecurity) error {

	sqlStatement := generateSQLtoInsertCompanies(securities)

	fmt.Print(sqlStatement)

	return nil
}

func generateSQLtoInsertCompanies(securities []MoexSecurity) string {
	sqlValuesStatement := "UPSERT INTO stock (company_name, is_public, isin, security_type, country_iso2, share_ticker, share_count) VALUES \n"

	for i, security := range securities {
		sqlValuesStatement += "("
		sqlValuesStatement += "\"" + strings.Replace(security.ShortName, "\"", "", -1) + "\", "
		sqlValuesStatement += "true, "
		sqlValuesStatement += "\"" + security.ISIN + "\", "
		sqlValuesStatement += "\"" + string(security.SecurityType) + "\", "
		sqlValuesStatement += "\"RU\", "
		sqlValuesStatement += "\"" + security.Ticker + "\", "
		sqlValuesStatement += fmt.Sprint(int(security.SharesIssued))
		sqlValuesStatement += ")"

		//YDB doesn't accept queries ending with a comma
		if i < len(securities)-1 {
			sqlValuesStatement += ", "
		}

		sqlValuesStatement += "\n"
	}

	return sqlValuesStatement
}
