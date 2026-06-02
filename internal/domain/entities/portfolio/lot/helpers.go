package lot

func FindLotsByFigi(lots []Lot, figi string) []Lot {
	filteredLots := []Lot{}

	for i := range lots {
		if lots[i].Figi == figi {
			filteredLots = append(filteredLots, lots[i])
		}
	}

	return filteredLots
}
