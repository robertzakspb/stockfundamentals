package bondservice

import "github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"

func GetBondFigis(bondList []bonds.Bond) []string {
	figis := []string{}

	for _, bond := range bondList {
		if bond.Figi == "" {
			continue
		}
		figis = append(figis, bond.Figi)
	}

	return figis
}
