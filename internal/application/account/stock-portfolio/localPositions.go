package portfolio

import (
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/application/shared"
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/portfolio/lot"
	"github.com/google/uuid"
)

// Returns positions that cannot be fetched from some external API and must thus be hardcoded here
func getHardCodedStockPositions() []lot.Lot {
	//TODO: Move this to the position_lot table:

	rosselHozId, _ := uuid.Parse(shared.ROSSELHOZ_ACCOUNT_ID)
	nlbId, _ := uuid.Parse(shared.NLB_ACCOUNT_ID_ID)

	jesvId := "BBG000BS7XH7"
	dunavId := "BBG000BMX476"
	mtlcId := "BBG000HP5RC7"
	nisId := "BBG0015L55D4"
	impolId := "BBG000HGH3F4"
	etalonId := "TCS50A10C1L6"

	serbianStocks := []lot.Lot{

		{
			Id:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Quantity:     125,
			PricePerUnit: 8457.9,
			Currency:     "RSD",
			AccountId:    nlbId,
			Figi:         jesvId,
			// Security: security.Stock{
			// 	Id:          jesvId,
			// 	Ticker:      "JESV",
			// 	Figi:        "BBG000BS7XH7",
			// 	CompanyName: "Jedinstvo Sevojno",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Id:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Quantity:     1232,
			PricePerUnit: 1495.58,
			Currency:     "RSD",
			AccountId:    nlbId,
			Figi:         dunavId,
			// Security: security.Stock{
			// 	Id:          dunavId,
			// 	Ticker:      "DNOS",
			// 	Figi:        "BBG000BMX476",
			// 	CompanyName: "Dunav Osiguranje",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Id:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Quantity:     459,
			PricePerUnit: 2018.78,
			Currency:     "RSD",
			AccountId:    nlbId,
			Figi:         mtlcId,
			// Security: security.Stock{
			// 	Id:          mtlcId,
			// 	Ticker:      "MTLC",
			// 	Figi:        "BBG000HP5RC7",
			// 	CompanyName: "Metalac",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Id:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Quantity:     380,
			PricePerUnit: 838.91,
			Currency:     "RSD",
			AccountId:    nlbId,
			Figi:         nisId,
			// Security: security.Stock{
			// 	Id:          nisId,
			// 	Ticker:      "NIIS",
			// 	Figi:        "BBG0015L55D4",
			// 	CompanyName: "NIS",
			// 	MIC:         "XBEL",
			// },
		},
		{
			Id:           uuid.New(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
			Quantity:     38,
			PricePerUnit: 7960.4,
			Currency:     "RSD",
			AccountId:    nlbId,
			Figi:         impolId,
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
		Id:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Quantity:     3459,
		PricePerUnit: 84,
		Currency:     "RUB",
		AccountId:    rosselHozId,
		Figi:         etalonId,
		// Security: security.Stock{
		// 	//TODO: Add security ID
		// 	Ticker:      "ETLN",
		// 	Figi:        "BBG00RM6M4V5", //TODO: Update it once the ISIN changes from US29760G1031 to RU000A10C1L6
		// 	CompanyName: "Эталон",
		// 	MIC:         "MISX",
		// },
	}

	allStocks := append(serbianStocks, rosselhozStocks)
	allStocks = append(allStocks, vtbLots()...)

	return allStocks
}

func vtbLots() []lot.Lot {
	lqdt := lot.Lot{
		AccountId:    uuid.MustParse(shared.VTB_ACCOUNT_ID),
		Id:           uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Quantity:     1191133,
		PricePerUnit: 1.9339,
		Currency:     "RUB",
		Figi:         "TCS60A1014L8",
	}

	return []lot.Lot{lqdt}
}
