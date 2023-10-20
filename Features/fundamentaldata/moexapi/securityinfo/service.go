package securityinfo

import (
	"encoding/json"
	"fmt"
	"net/http"

	typeconverter "github.com/compoundinvest/stockfundamentals/Utilities/converters"
)

func FetchSecuritiesInfoFromMoex() {
	endpointURL := "https://iss.moex.com/iss/engines/stock/markets/shares/boards/TQBR/securities.json?iss.meta=off&iss.only=securities"

	response, err := http.Get(endpointURL)
	if err != nil {
		fmt.Println("Unable to fetch moex securities' info: ", err)
	}
	defer response.Body.Close()

	securitiesDTO := MoexSecuritiesDTO{}
	err = json.NewDecoder(response.Body).Decode(&securitiesDTO)
	if err != nil {
		fmt.Println("Failed to parse the Moex Securities JSON: ", err)
	}

	parsedSecurities := []MoexSecurity{}
	for _, security := range securitiesDTO.Securities.Data {
		ticker := security[0].(string)
		shortName := security[2].(string)
		fullName := security[9].(string)
		englishName := security[20].(string)
		lotSize, _ := typeconverter.GetFloat(security[4])
		faceValue, _ := typeconverter.GetFloat(security[5])
		sharesIssued, _ := typeconverter.GetFloat(security[18])
		isin := security[19].(string)
		securityType := parseSecurityType(security[24])

		parsedSecurity := MoexSecurity{
			Ticker:       ticker,
			ShortName:    shortName,
			FullName:     fullName,
			EnglishName:  englishName,
			LotSize:      lotSize,
			FaceValue:    faceValue,
			SharesIssued: sharesIssued,
			ISIN:         isin,
			SecurityType: securityType,
		}
		parsedSecurities = append(parsedSecurities, parsedSecurity)
	}

	saveMoexSecuritiesToDB(parsedSecurities)
}

func parseSecurityType(moexSecurityType interface{}) SecurityType {
	parsedSecurityType, ok := moexSecurityType.(string)
	if !ok {
		return OrdinaryShare
	}

	switch parsedSecurityType {
	case "1":
		return OrdinaryShare
	case "2":
		return PreferredShare
	case "D":
		return DepositoryReceipt
	default:
		return OrdinaryShare
	}
}
