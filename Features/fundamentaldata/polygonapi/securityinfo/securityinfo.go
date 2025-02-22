package securityinfo

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	polygon "github.com/compoundinvest/stockfundamentals/Features/fundamentaldata/polygonapi"
)

func FetchSecuritiesInfo(db *sql.DB) error {
	polygon.DelayRequestIfAPILimitReached()

	endpointURL := "https://api.polygon.io/v3/reference/tickers?market=stocks&active=true&order=asc&limit=1000&sort=ticker&apiKey=" //TODO: the API key here
	response, err := http.Get(endpointURL)
	if err != nil {
		fmt.Println("Unable to fetch securities info")
	}
	defer response.Body.Close()

	var apiDTO PolygonSecurityDTO
	json.NewDecoder(response.Body).Decode(&apiDTO)
	//Parsing the first 1000 securities from Polygon API
	parsedSecurities := apiDTO.convertToPolygonSecurity()

	//Iterating through the remaining securities via Polygon API
	for apiDTO.NextURL != "" {
		polygon.DelayRequestIfAPILimitReached()

		response, err := http.Get(apiDTO.NextURL + "&apiKey=") //Missing the API key here
		if err != nil {
			fmt.Println("Unable to fetch securities info")
			break
		}

		apiDTO = PolygonSecurityDTO{}
		json.NewDecoder(response.Body).Decode(&apiDTO)
		if apiDTO.Error != "" {
			fmt.Println("Encountered an error:", apiDTO.Error)
		}
		parsedSecurities = append(parsedSecurities, apiDTO.convertToPolygonSecurity()...)
	}

	if len(parsedSecurities) > 0 {
		saveSecurityInfoToDB(parsedSecurities, db)
		return nil
	} else {
		return fmt.Errorf("trying to save zero Polygon API securities to the DB")
	}
}
