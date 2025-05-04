package tinkoffapi

//TODO: - Delete the file

import (
	// "fmt"
	// "time"

	// "github.com/compoundinvest/stockfundamentals/infrastructure/logger"
	// tinkoff "github.com/russianinvestments/invest-api-go-sdk/investgo"
	// investapi "github.com/russianinvestments/invest-api-go-sdk/proto"
)

func FetchDividendsForAllSecurities() {

}

// func FetchDividendsFor() {

// 	client, err := newTinkoffAPIClient()
// 	if err != nil {
// 		return
// 	}

// 	securityService := client.NewInstrumentsServiceClient()


// 	earliestDividendTime := time.Date(1971, 1, 1, 0, 0, 0, 0, time.UTC) 
// 	upcomingDividendTime := time.Now().AddDate(1, 0, 0)
// 	dividends, err := securityService.GetDividents("RU000A0DKVS5", earliestDividendTime, upcomingDividendTime) // BBG00475KKY8
// 	if err != nil {
// 		logger.Log(err.Error(), logger.ERROR)
// 	}
// 	for _, dividend := range dividends.GetDividends() {
// 		if dividend == nil {
// 			continue
// 		}

// 		fmt.Println("Novatek dividends:")

// 		fmt.Println("Record date: ", time.Unix(dividend.RecordDate.Seconds, 0), "; DPS:", dividend.DividendNet.Units, ".", dividend.GetDividendNet().Nano)
// 		// fmt.Println("DPS net: ", dividend.DividendNet)
// 		// fmt.Println("Dividend type: ", dividend.DividendType)
// 		// fmt.Println("Payment date: ", time.Unix(dividend.PaymentDate.Seconds, 0))
// 		// fmt.Println("Regularity ", dividend.Regularity)
// 		// fmt.Println("Yield value:", dividend.YieldValue)
// 		fmt.Println("------------------------------------")
// 	}
// }

// func saveDividendsToDB() {

// }
