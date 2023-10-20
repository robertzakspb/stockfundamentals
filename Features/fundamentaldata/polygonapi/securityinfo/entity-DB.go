package securityinfo

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/compoundinvest/stockfundamentals/Utilities/stringutils"
	"github.com/google/uuid"
)

func saveSecurityInfoToDB(securities []PolygonSecurity, db *sql.DB) error {
	sqlStatement := generateSQLStatementForSecuritiesInfo(securities)

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println("Failed to execute the query. " + err.Error())
		return err
	}
	defer rows.Close()

	return nil
}

func generateSQLStatementForSecuritiesInfo(securities []PolygonSecurity) string {
	sqlStatement := "INSERT INTO company VALUES "

	for _, security := range securities {
		uuid, err := uuid.NewRandom()
		if err != nil {
			continue
		}
		//Replacing single quotation marks with pairs of single quotation marks for Postgres
		sanitizedName := strings.ReplaceAll(security.Name, "'", "''")
		sqlStatement += "(" + "'" + uuid.String() + "', " + "'" + sanitizedName + "'" + ", " + "'stock', " + "'US', " + "true, " + "0, " + "'" + security.Ticker + "'" + ", " + "0, " + "NULL" + ")" + ","
	}

	sqlStatement = stringutils.TrimLastCharacter(sqlStatement) //removing the last comma
	return sqlStatement
}
