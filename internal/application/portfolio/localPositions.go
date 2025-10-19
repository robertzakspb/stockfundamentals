package portfolio

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/google/uuid"
)



// Returns positions that cannot be fetched from some external API and must thus be hardcoded here
func getHardCodedStockPositions() []lot.Lot {
	//TODO: Move this elsewhere:

	rosselHozId, _ := uuid.Parse("5e3e1fdb-5c18-43a5-a7c6-f898aff2d17f")
	nlbId, _ := uuid.Parse("3b450479-a136-4ecd-9f34-8bfac6488101")

	jesvId, _ := uuid.Parse("dd194350-4c61-4643-8c74-1120ceca8fae")
	dunavId, _ := uuid.Parse("4f96e511-34db-4f00-9bf5-0975cef04c2b")
	mtlcId, _ := uuid.Parse("8f0161f0-083d-431a-9fff-89deb073ce0f")
	nisId, _ := uuid.Parse("9f45ae88-70bd-48ef-aadc-aa6f01377b76")
	impolId, _ := uuid.Parse("d281b3b3-ac2a-49d2-9265-d030fe4142a9")
	etalonId, _ := uuid.Parse("08192dec-4141-4798-a342-0b62894285a2")

	serbianStocks := []lot.Lot{

		{
			Quantity:     43,
			PricePerUnit: 8457.9,
			Currency:     "RSD",
			AccountId:    nlbId,
			SecurityId:   jesvId,
			// Security: security.Stock{
			// 	Id:          jesvId,
			// 	Ticker:      "JESV",
			// 	Figi:        "BBG000BS7XH7",
			// 	CompanyName: "Jedinstvo Sevojno",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Quantity:     567,
			PricePerUnit: 1495.58,
			Currency:     "RSD",
			AccountId:    nlbId,
			SecurityId:   dunavId,
			// Security: security.Stock{
			// 	Id:          dunavId,
			// 	Ticker:      "DNOS",
			// 	Figi:        "BBG000BMX476",
			// 	CompanyName: "Dunav Osiguranje",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Quantity:     289,
			PricePerUnit: 2018.78,
			Currency:     "RSD",
			AccountId:    nlbId,
			SecurityId:   mtlcId,
			// Security: security.Stock{
			// 	Id:          mtlcId,
			// 	Ticker:      "MTLC",
			// 	Figi:        "BBG000HP5RC7",
			// 	CompanyName: "Metalac",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Quantity:     380,
			PricePerUnit: 838.91,
			Currency:     "RSD",
			AccountId:    nlbId,
			SecurityId:   nisId,
			// Security: security.Stock{
			// 	Id:          nisId,
			// 	Ticker:      "NIIS",
			// 	Figi:        "BBG0015L55D4",
			// 	CompanyName: "NIS",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Quantity:     28,
			PricePerUnit: 7960.4,
			Currency:     "RSD",
			AccountId:    nlbId,
			SecurityId:   impolId,
			// Security: security.Stock{
			// 	Id:          impolId,
			// 	Ticker:      "IMPL",
			// 	Figi:        "BBG000HGH3F4",
			// 	CompanyName: "Impol Seval",
			// 	MIC:         "XBEL",
			// },
		},
	}

	rosselhozStocks := lot.Lot{
		Quantity:     3459,
		PricePerUnit: 84,
		Currency:     "RUB",
		AccountId:    rosselHozId,
		SecurityId: etalonId,
		// Security: security.Stock{
		// 	//TODO: Add security ID
		// 	Ticker:      "ETLN",
		// 	Figi:        "BBG00RM6M4V5", //TODO: Update it once the ISIN changes from US29760G1031 to RU000A10C1L6
		// 	CompanyName: "Эталон",
		// 	MIC:         "MISX",
		// },
	}
	
	allStocks := append(serbianStocks, rosselhozStocks)

	return allStocks
}
