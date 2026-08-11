package forexservice

import (
	"errors"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	typeconverter "github.com/compoundinvest/stockfundamentals/internal/utilities/converters"
	timehelpers "github.com/compoundinvest/stockfundamentals/internal/utilities/time-helpers"
)

func fetchUsdToRsdRate(date time.Time) (ForexRate, error) {
	url := makeNbsUrl(date)
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	content, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	html := string(content)

	ratesTableHtml := "<table class=\"table\">"
	indexOfratesTableHtml := strings.Index(html, ratesTableHtml)
	if indexOfratesTableHtml == -1 {
		return ForexRate{}, errors.New("Failed to find the exchange rate table in the HTML")
	}
	html = html[indexOfratesTableHtml:]

	tbodyIndex := strings.Index(html, "<tbody>")
	if tbodyIndex == -1 {
		return ForexRate{}, errors.New("Failed to find the tbody tag in the HTML")
	}
	html = html[tbodyIndex:]

	for {
		trIndex := strings.Index(html, "<tr>")
		if trIndex == -1 {
			return ForexRate{}, errors.New("Failed to find the <tr> tag in the HTML")
		}
		html = html[trIndex:]

		tdIndex := strings.Index(html, "<td>")
		if tdIndex == -1 {
			return ForexRate{}, errors.New("Failed to find the <td> tag with the currency name in the HTML")
		}
		html = html[tdIndex:]

		if len(html) < 7 { //We are lookign for the <td>USD tag which contains at least 7 characters
			return ForexRate{}, errors.New("Failed to find the currency name in the HTML")
		}

		currency := html[4:7]
		if currency != "USD" {
			continue
		}

		html = html[4:]                       //Skiping the <td> with the currency name to look for the next <td> with currency code
		tdIndex = strings.Index(html, "<td>") //The currency code td
		if tdIndex == -1 {
			return ForexRate{}, errors.New("Failed to find the <td> tag containing the currency code in the HTML")
		}
		html = html[tdIndex+4:]

		tdIndex = strings.Index(html, "<td>") //The country name td
		if tdIndex == -1 {
			return ForexRate{}, errors.New("Failed to find the <td> tag containing the country name in the HTML")
		}
		html = html[tdIndex+4:]

		tdIndex = strings.Index(html, "<td>") //The lot size td
		if tdIndex == -1 {
			return ForexRate{}, errors.New("Failed to find the <td> tag containing the lot size in the HTML")
		}
		html = html[tdIndex+4:]

		tdIndex = strings.Index(html, "<td>") //The exchange rate td
		if tdIndex == -1 {
			return ForexRate{}, errors.New("Failed to find the <td> tag containing the exchange rate in the HTML")
		}
		html = html[tdIndex+4:]

		closingTagIndex := strings.Index(html, "</td>") //Finding the closing tag to extract the rate between the opening and closing tag

		rateString := html[:closingTagIndex]
		rateString = strings.Replace(rateString, ",", ".", 1)
		rateFloat, err := typeconverter.GetFloat(rateString)
		if err != nil {
			return ForexRate{}, err
		}

		logger.Log("Fetched the USD/RSD exchange rate for "+date.String()+": "+rateString, logger.INFORMATION)

		rate := ForexRate{
			Currency1: "USD",
			Currency2: "RSD",
			Date:      date,
			Rate:      rateFloat,
		}
		return rate, nil
	}
}

func makeNbsUrl(date time.Time) string {
	var sb strings.Builder
	sb.WriteString("https://webappcenter.nbs.rs/ExchangeRateWebApp/ExchangeRate/IndexByDate?isSearchExecuted=true&Date=")

	sb.WriteString(timehelpers.DateInDDMMYYYFormat(date))

	sb.WriteString(".&ExchangeRateListTypeID=3")

	return sb.String()
}
