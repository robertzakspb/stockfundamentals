package bondportfolio

import (
	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
)

func matchLotsWithBonds(lots []bonds.BondLot, bonds []bonds.Bond) []bonds.BondLot {
	if len(bonds) == 0 {
		return lots
	}

	for i, lot := range lots {
		for _, b := range bonds {
			if lot.Figi == b.Figi || lot.Isin == b.Isin {
				lots[i].Bond = b
			}
		}
	}

	return lots
}

func GetLotBonds(lots []bonds.BondLot) []bonds.Bond {
	bondList := []bonds.Bond{}

	for _, lot := range lots {
		if lot.Bond.Figi == "" && lot.Bond.Isin == "" && lot.Bond.NominalValue == 0 {
			continue
		}
		bondList = append(bondList, lot.Bond)
	}

	return bondList
}

func GetLotFigis(lots []bonds.BondLot) []string {
	figis := []string{}
	for _, lot := range lots {
		if lot.Figi != "" {
			figis = append(figis, lot.Figi)
		}
	}
	return figis
}

func GetLotIsins(lots []bonds.BondLot) []string {
	isins := []string{}
	for _, lot := range lots {
		if lot.Isin != "" {
			isins = append(isins, lot.Isin)
		}
	}
	return isins
}