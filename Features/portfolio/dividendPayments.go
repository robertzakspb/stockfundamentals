package portfolio

type DividendPayment struct {
	Ticker string
	DPS float64
	Quantity int
}

func (payment DividendPayment) GrossPayout() float64 {
	return payment.DPS * float64(payment.Quantity)
}

func (portfolio Portfolio) UpcomingDividends() []DividendPayment  {
	payments := []DividendPayment{}
	for _, position := range portfolio.UniquePositions() {
		dividend := DividendPayment{
			position.Ticker,
			projectedDPSFor(position.Ticker), 
			int(position.Quantity),
		}
		payments = append(payments, dividend)
	}
	
	return payments
}

func projectedDPSFor(ticker string) float64 {
	projectedDPS := map[string]float64 {
		"NVTK": 46,
		"SNGSP": 10,
		"SIBN": 70,
		"ETLN": 12,
		"CHMF": 100,
		"NLMK": 15,
		"MAGN": 4,
		"LKOH": 1000,
		"TATNP": 100,
		"ROSN": 60,
		"NMTP": 1,
		"JESV": 647,
		"DNOS": 121,
		"MTLC": 100,
		"NIIS": 30,
	}
	
	return projectedDPS[ticker]
}