package portfolio

// Returns positions that cannot be fetched from some external API and must thus be hardcoded here
func getHardCodedStockPositions() []Lot {
	serbianStocks := []Lot{
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
			Quantity:     25,
			OpeningPrice: 1296,
			Currency:     "RSD",
			Ticker:       "DNOS",
			CompanyName:  "Dunav Osiguranje",
			BrokerName:   "NLB",
			MIC:          "XBEL",
		},
		{
			Quantity:     13,
			OpeningPrice: 1949,
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

	sberbankStocks := []Lot{
		{
			Quantity:     598,
			OpeningPrice: 1620,
			Currency:     "RUB",
			Ticker:       "MBNK",
			CompanyName:  "МТС Банк",
			BrokerName:   "Sberbank",
			MIC:          "MISX",
		},
		{
			Quantity:     3010,
			OpeningPrice: 54.6,
			Currency:     "RUB",
			Ticker:       "SNGSP",
			CompanyName:  "Сургутнефтегаз-п",
			BrokerName:   "Sberbank",
			MIC:          "MISX",
		},
		{
			Quantity:     37,
			OpeningPrice: 4079,
			Currency:     "RUB",
			Ticker:       "HEAD",
			CompanyName:  "HeadHunter",
			BrokerName:   "Sberbank",
			MIC:          "MISX",
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
	allStocks := append(serbianStocks, sberbankStocks...)
	allStocks = append(allStocks, americanStocks...)
	allStocks = append(allStocks, rosselhozStocks)

	return allStocks
}
