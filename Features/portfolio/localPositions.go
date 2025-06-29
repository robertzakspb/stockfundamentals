package portfolio

// Returns positions that cannot be fetched from some external API and must thus be hardcoded here
func getHardCodedStockPositions() []Lot {
	serbianStocks := []Lot{
		//TODO: - Fill out the ISINs
		//TODO: Update after the Jun/July? updates
		{
			Quantity:     5,
			OpeningPrice: 7772,
			Currency:     "RSD",
			Ticker:       "JESV",
			CompanyName:  "Jedinstvo Sevojno",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     431,
			OpeningPrice: 1407.30,
			Currency:     "RSD",
			Ticker:       "DNOS",
			CompanyName:  "Dunav Osiguranje",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     277,
			OpeningPrice: 1961.57,
			Currency:     "RSD",
			Ticker:       "MTLC",
			CompanyName:  "Metalac",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     380,
			OpeningPrice: 830.2,
			Currency:     "RSD",
			Ticker:       "NIIS",
			CompanyName:  "NIS",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
	}

	americanStocks := []Lot{
		{
			Quantity:     1373,
			OpeningPrice: 15.07,
			Currency:     "USD",
			Ticker:       "CSIQ",
			CompanyName:  "Canadian Solar",
			BrokerName:   "IBKR",
			MIC:          "XNAS",
		},
	}

	rosselhozStocks := Lot{
		Quantity:     3459,
		OpeningPrice: 84,
		Currency:     "RUB",
		Ticker:       "ETLN",
		CompanyName:  "Эталон",
		BrokerName:   "Rosselhoz",
		MIC:          "MISX",
	}
	allStocks := append(serbianStocks, americanStocks...)
	allStocks = append(allStocks, rosselhozStocks)

	return allStocks
}
