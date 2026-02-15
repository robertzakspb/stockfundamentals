package portfolio

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/google/uuid"
)



// Returns positions that cannot be fetched from some external API and must thus be hardcoded here
func getHardCodedStockPositions() []lot.Lot {
	//TODO: Move this to the position_lot table:

	rosselHozId, _ := uuid.Parse("5e3e1fdb-5c18-43a5-a7c6-f898aff2d17f")
	nlbId, _ := uuid.Parse("3b450479-a136-4ecd-9f34-8bfac6488101")

	jesvId := "BBG000BS7XH7"
	dunavId := "BBG000BMX476"
	mtlcId := "BBG000HP5RC7"
	nisId := "BBG0015L55D4"
	impolId := "BBG000HGH3F4"
	//TODO: Look up the figi in the DB
	etalonId := "08192dec-4141-4798-a342-0b62894285a2"

	serbianStocks := []lot.Lot{

		{
			Quantity:     125,
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
			Quantity:     1232,
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
			Quantity:     459,
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
			Quantity:     38,
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
