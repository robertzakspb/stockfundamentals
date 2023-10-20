package moexdividends

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchDividendsFor(ticker string) (Dividends, error) {
	return fetchDividendsFor(ticker)
}

func fetchDividendsFor(ticker string) (Dividends, error) {
	if ticker == "" {
		return Dividends{}, fmt.Errorf("empty ticker provided when fetching moex dividends")
	}

	endpointURL := generateDividendURL(ticker)

	response, err := http.Get(endpointURL)
	if err != nil {
		fmt.Println("Unable to fetch moex dividends for ", ticker, ". ", err)
		return Dividends{}, err
	}
	defer response.Body.Close()

	dividendDTO := MoexDividendDTO{}
	err = json.NewDecoder(response.Body).Decode(&dividendDTO)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(dividendDTO); i++ {
		if len(dividendDTO[i].Dividends) == 0 {
			dividendDTO[i] = dividendDTO[len(dividendDTO)-1]
			dividendDTO = dividendDTO[:len(dividendDTO)-1]
		}
	}

	JSONisEmpty := len(dividendDTO) == 0
	dividendArryIsEmpty := len(dividendDTO[0].Dividends) == 0
	if JSONisEmpty || dividendArryIsEmpty {
		fmt.Println("Missing dividend information for", ticker)
		return Dividends{}, nil
	}

	dividends := dividendDTO.asDividends()

	return dividends, nil
}

func generateDividendURL(ticker string) string {
	const moexDividendBaseURL = "https://iss.moex.com/iss/securities/"
	const moexDividendURLEnd = "/dividends.json?iss.meta=off&iss.json=extended"
	return moexDividendBaseURL + ticker + moexDividendURLEnd
}
