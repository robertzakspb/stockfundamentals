package portfolio

import (
	"github.com/compoundinvest/stockfundamentals/features/fundamentaldata/security"
	"github.com/google/uuid"
)

var tinkoffIisId, _ = uuid.Parse("3315bd1c-12a4-444e-a294-84ef339e26e1")

// Returns positions that cannot be fetched from some external API and must thus be hardcoded here
func getHardCodedStockPositions() []Lot {
	//TODO: Move this elsewhere:

	rosselHozId, _ := uuid.Parse("5e3e1fdb-5c18-43a5-a7c6-f898aff2d17f")
	ibkrId, _ := uuid.Parse("2f23017e-566b-48d3-bbb0-3c3766f9e560")
	nlbId, _ := uuid.Parse("3b450479-a136-4ecd-9f34-8bfac6488101")

	jesvId, _ := uuid.Parse("dd194350-4c61-4643-8c74-1120ceca8fae")
	dunavId, _ := uuid.Parse("4f96e511-34db-4f00-9bf5-0975cef04c2b")
	mtlcId, _ := uuid.Parse("8f0161f0-083d-431a-9fff-89deb073ce0f")
	nisId, _ := uuid.Parse("9f45ae88-70bd-48ef-aadc-aa6f01377b76")
	impolId, _ := uuid.Parse("d281b3b3-ac2a-49d2-9265-d030fe4142a9")
	csiqId, _ := uuid.Parse("6acbe4cc-80e1-4ab1-9256-d16661c7001a")

	serbianStocks := []Lot{

		{
			Quantity:     25,
			PricePerUnit: 8424.28,
			Currency:     "RSD",
			AccountId:    nlbId,
			Security: security.Stock{
				Id:          jesvId,
				Ticker:      "JESV",
				Figi:        "BBG000BS7XH7",
				CompanyName: "Jedinstvo Sevojno",
				MIC:         "XBEL",
			},
		},
		{
			Quantity:     534,
			PricePerUnit: 1484.5,
			Currency:     "RSD",
			AccountId:    nlbId,
			Security: security.Stock{
				Id:          dunavId,
				Ticker:      "DNOS",
				Figi:        "BBG000BMX476",
				CompanyName: "Dunav Osiguranje",
				MIC:         "XBEL",
			},
		},
		{
			Quantity:     289,
			PricePerUnit: 2018.78,
			Currency:     "RSD",
			AccountId:    nlbId,
			Security: security.Stock{
				Id:          mtlcId,
				Ticker:      "MTLC",
				Figi:        "BBG000HP5RC7",
				CompanyName: "Metalac",
				MIC:         "XBEL",
			},
		},
		{
			Quantity:     380,
			PricePerUnit: 838.91,
			Currency:     "RSD",
			AccountId:    nlbId,
			Security: security.Stock{
				Id:          nisId,
				Ticker:      "NIIS",
				Figi:        "BBG0015L55D4",
				CompanyName: "NIS",
				MIC:         "XBEL",
			},
		},
		{
			Quantity:     18,
			PricePerUnit: 7939.21,
			Currency:     "RSD",
			AccountId:    nlbId,
			Security: security.Stock{
				Id:          impolId,
				Ticker:      "IMPL",
				Figi:        "BBG000HGH3F4",
				CompanyName: "Impol Seval",
				MIC:         "XBEL",
			},
		},
	}

	americanStocks := []Lot{
		{
			Quantity:     1373,
			PricePerUnit: 15.07,
			Currency:     "USD",
			AccountId:    ibkrId,
			Security: security.Stock{
				Id:          csiqId,
				Ticker:      "CSIQ",
				Figi:        "BBG000K1JFJ0",
				CompanyName: "Canadian Solar",
				MIC:         "XNAS",
			},
		},
	}

	rosselhozStocks := Lot{
		Quantity:     3459,
		PricePerUnit: 84,
		Currency:     "RUB",
		AccountId:    rosselHozId,
		Security: security.Stock{
			//TODO: Add security ID
			Ticker:      "ETLN",
			Figi:        "BBG00RM6M4V5", //TODO: Update it once the ISIN changes from US29760G1031 to RU000A10C1L6
			CompanyName: "Эталон",
			MIC:         "MISX",
		},
	}
	allStocks := append(serbianStocks, americanStocks...)
	allStocks = append(allStocks, rosselhozStocks)

	return allStocks
}
