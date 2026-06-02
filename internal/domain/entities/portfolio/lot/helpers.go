package lot

func FindLotIndicesByFigi(lots []Lot, figi string) []int {
	filteredLotIndices := []int{}

	for i := range lots {
		if lots[i].Figi == figi {
			filteredLotIndices = append(filteredLotIndices, i)
		}
	}

	return filteredLotIndices
}
