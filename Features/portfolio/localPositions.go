package portfolio

// Returns positions that cannot be fetched from some external API and must thus be hardcoded here
func getHardCodedStockPositions() []Lot {
	serbianStocks := []Lot{
		{
			Quantity:     25,
			OpeningPrice: 8424.28,
			Currency:     "RSD",
			Ticker:       "JESV",
			Figi:         "BBG000BS7XH7",
			CompanyName:  "Jedinstvo Sevojno",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     534,
			OpeningPrice: 1484.5,
			Currency:     "RSD",
			Ticker:       "DNOS",
			Figi:         "BBG000BMX476",
			CompanyName:  "Dunav Osiguranje",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     289,
			OpeningPrice: 2018.78,
			Currency:     "RSD",
			Ticker:       "MTLC",
			Figi:         "BBG000HP5RC7",
			CompanyName:  "Metalac",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     380,
			OpeningPrice: 838.91,
			Currency:     "RSD",
			Ticker:       "NIIS",
			Figi:         "BBG0015L55D4",
			CompanyName:  "NIS",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     18,
			OpeningPrice: 7939.21,
			Currency:     "RSD",
			Ticker:       "IMPL",
			Figi:         "BBG000HGH3F4",
			CompanyName:  "Impol Seval",
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
			Figi:         "BBG000K1JFJ0",
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
		Figi:         "BBG00RM6M4V5", //TODO: Update it once the ISIN changes from US29760G1031 to RU000A10C1L6
		CompanyName:  "Эталон",
		BrokerName:   "Rosselhoz",
		MIC:          "MISX",
	}
	allStocks := append(serbianStocks, americanStocks...)
	allStocks = append(allStocks, rosselhozStocks)

	return allStocks
}
