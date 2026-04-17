package bondportfolio

import (
	"fmt"
	"strconv"

	"github.com/compoundinvest/stockfundamentals/internal/infrastructure/logger"
	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

const CALENDAR_SHEET_TITLE = "Календарь Портфеля"

func GenerateTimeLineExcel() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	err := f.SetSheetName("Sheet1", CALENDAR_SHEET_TITLE)
	if err != nil {
		logger.Log(err.Error(), logger.ERROR)
		return err
	}

	f.SetCellValue(CALENDAR_SHEET_TITLE, "A1", "Дата")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "B1", "Эмитент")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "C1", "Событие")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "D1", "Сумма (₽)")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "E1", "Налог (₽)")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "F1", "Сумма ($)")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "G1", "Налог ($)")

	lots, err := GetAccountPositions(uuid.MustParse("129274f9-ee80-4e74-aa1c-fea578bac6e6"))
	if err != nil {
		return err
	}

	lots, err = PopulateLotsWithBonds(lots)
	if err != nil {
		return err
	}

	lots = PopulateLotsWithCoupons(lots)

	timeline, err := generateTimeLineForLots(lots, false)

	currentRow := 2                                                               //The first row is reserved for headers
	currentRow = EnterTimelineInformationIntoSpreadsheet(f, timeline, currentRow) //FIXME: Change the method's name

	if err := f.SaveAs("Portfolio_Calendar.xlsx"); err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}
	return nil
}

func EnterTimelineInformationIntoSpreadsheet(f *excelize.File, timeline []TimeLineItem, currentRow int) int {
	totalRubPayout := 0.0
	totalRubTaxes := 0.0
	totalUsdPayout := 0.0
	totalUsdTaxes := 0.0

	for _, event := range timeline {
		year, month, day := event.Timestamp.Date()
		formattedDate := strconv.Itoa(day) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(year)
		f.SetCellValue(CALENDAR_SHEET_TITLE, "A"+strconv.Itoa(currentRow), formattedDate)
		f.SetCellValue(CALENDAR_SHEET_TITLE, "C"+strconv.Itoa(currentRow), event.EventName)
		f.SetCellValue(CALENDAR_SHEET_TITLE, "B"+strconv.Itoa(currentRow), event.BondName)

		switch event.Currency {
		case "RUB":
			payout := event.Amount
			if payout == 0.0 {
				continue
			}

			f.SetCellValue(CALENDAR_SHEET_TITLE, "D"+strconv.Itoa(currentRow), "₽"+fmt.Sprintf("%.1f", payout))

			if event.EventName == "Выплата купона" {
				totalRubTaxes += payout * 0.13
				f.SetCellValue(CALENDAR_SHEET_TITLE, "E"+strconv.Itoa(currentRow), "₽"+fmt.Sprintf("%.1f", payout*0.13))

			}

			totalRubPayout += payout
		case "USD":
			payout := event.Amount
			if payout == 0.0 {
				continue
			}

			f.SetCellValue(CALENDAR_SHEET_TITLE, "F"+strconv.Itoa(currentRow), "$"+fmt.Sprintf("%.1f", payout))

			if event.EventName == "Выплата купона" {
				f.SetCellValue(CALENDAR_SHEET_TITLE, "G"+strconv.Itoa(currentRow), "$"+fmt.Sprintf("%.1f", payout*0.13))
				totalUsdTaxes += payout * 0.13
			}

			totalUsdPayout += payout
		default:
			logger.Log("Unexpected currency in the timeline: "+event.Currency, logger.WARNING)
		}

		currentRow++
	}

	currentRow++

	style, _ := f.NewStyle(&excelize.Style{
		Fill: excelize.Fill{Type: "pattern", Color: []string{"E0EBF5"}, Pattern: 1},
	})
	f.SetCellStyle(CALENDAR_SHEET_TITLE, "C"+strconv.Itoa(currentRow), "G"+strconv.Itoa(currentRow), style)

	f.SetCellValue(CALENDAR_SHEET_TITLE, "C"+strconv.Itoa(currentRow), "Итого: ")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "D"+strconv.Itoa(currentRow), "₽"+fmt.Sprintf("%.1f", totalRubPayout))
	f.SetCellValue(CALENDAR_SHEET_TITLE, "E"+strconv.Itoa(currentRow), "₽"+fmt.Sprintf("%.1f", totalRubTaxes))
	f.SetCellValue(CALENDAR_SHEET_TITLE, "F"+strconv.Itoa(currentRow), "$"+fmt.Sprintf("%.1f", totalUsdPayout))
	f.SetCellValue(CALENDAR_SHEET_TITLE, "G"+strconv.Itoa(currentRow), "$"+fmt.Sprintf("%.1f", totalUsdTaxes))
	currentRow++

	return currentRow
}
