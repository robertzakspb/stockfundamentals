package fundamentals

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func FetchFundamentalDataFor(ticker string) (fundamentals StockFundamentals, err error) {
	fundamentalData, err := fetchFundamentalDataFor(ticker)
	if err != nil {
		fmt.Println(err)
	}
	if len(fundamentalData.Results) == 0 {
		return StockFundamentals{}, fmt.Errorf("API response contains zero financial metrics for %v", ticker)
	}

	return fundamentalData.AsStockFundamentals(), nil
}

func fetchFundamentalDataFor(ticker string) (fundamentals PolygonFundamentalDataDTO, err error) {

	endpointURL := "https://api.polygon.io/vX/reference/financials?ticker=" + ticker + "&timeframe=annual&order=asc&limit=100&sort=period_of_report_date&apiKey=" //Missing the API key here

	response, err := http.Get(endpointURL)
	if err != nil {
		fmt.Println("Unable to fetch fundamental data for ", ticker, ". ", err)
		return PolygonFundamentalDataDTO{}, err
	}
	defer response.Body.Close()

	var fundamentalsDTO PolygonFundamentalDataDTO
	json.NewDecoder(response.Body).Decode(&fundamentalsDTO)

	fundamentalsDTO.Ticker = ticker

	return fundamentalsDTO, nil
}
