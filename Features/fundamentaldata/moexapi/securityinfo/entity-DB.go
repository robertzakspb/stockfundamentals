package securityinfo

import (
	"fmt"

	"github.com/compoundinvest/stockfundamentals/Utilities/sliceutils"
	"github.com/google/uuid"
)

func saveMoexSecuritiesToDB(securities []MoexSecurity) error {

	for _, security := range securities {
		if string(security.Ticker[len(security.Ticker)-1]) == "P" {
			fmt.Println(security.Ticker)
		}
	}

	sqlStatement := "UPSERT INTO " + COMPANY_TABLE_NAME + COMPANY_TABLE_COLUMN_LIST + " VALUES "

	sqlStatement += generateSQLforCompaniesWithPreferredShares(securities)
	sqlStatement += generateSQLforCompaniesWithoutPreferredShares(securities)

	fmt.Print(sqlStatement)

	return nil
}

func generateSQLforCompaniesWithPreferredShares(securities []MoexSecurity) string {
	sqlValuesStatement := ""

	parsedTickers := []string{}

	for _, security := range securities {
		if companyHasPreferredShares(security.Ticker) == false {
			continue
		}

		if sliceutils.Contains(parsedTickers, security.Ticker) {
			continue
		}

		// Generating a UUID for the DB table's primary key
		// uuid, err := uuid.NewRandom()
		// if err != nil {
		// return ""
		// }

		// parsedTickers = append(parsedTickers, security.Ticker[])
	}

	return sqlValuesStatement
}

func generateSQLforCompaniesWithoutPreferredShares(securities []MoexSecurity) string {
	sqlValuesStatement := ""

	for i, security := range securities {
		if companyHasPreferredShares(security.Ticker) {
			continue //Skipping such companies, as they are handled by a different function
		}
		// Generating a UUID for the DB table's primary key
		uuid, err := uuid.NewRandom()
		if err != nil {
			return ""
		}

		sqlValuesStatement += "(\"" + uuid.String() + "\", \"" + security.ShortName + "\", \"" + string(security.SecurityType) + "\", " + "\"RU\", " + "true, " + fmt.Sprint(int(security.SharesIssued)) + ", \"" + security.Ticker + "\", " + "null, " + "null)"

		if i < len(securities)-1 {
			sqlValuesStatement += ", "
		}
	}

	return sqlValuesStatement
}

func companyHasPreferredShares(targetTicker string) bool {
	companiesWithPreferredShares := []string{"BANEP", "BISVP", "BSPBP", "CNTLP", "DZRDP", "GAZAP", "HIMCP", "IGSTP", "JNOSP", "KAZTP", "KCHEP", "KGKCP", "KRKNP", "KRKOP", "KROTP", "KRSBP", "KTSBP", "KZOSP", "LNZLP", "LSNGP", "MAGEP", "MFGSP", "MGTSP", "MISBP", "MTLRP", "NKNCP", "NNSBP", "OMZZP", "PMSBP", "RTKMP", "RTSBP", "SAGOP", "SAREP", "SBERP", "SNGSP", "STSBP", "TASBP", "TATNP", "TGKBP", "TORSP", "TRNFP", "VGSBP", "VJGZP", "VRSBP", "VSYDP", "WTCMP", "YKENP", "YRSBP"}

	hasPreferredShares := false
	for _, ticker := range companiesWithPreferredShares {
		if targetTicker == ticker || targetTicker == ticker+"P" {
			hasPreferredShares = true
		}
	}

	return hasPreferredShares
}
