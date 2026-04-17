package bondportfolio

import (
	"fmt"
	"strconv"
	"time"

	"github.com/compoundinvest/stockfundamentals/internal/domain/entities/bonds"
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

	f.SetCellValue(CALENDAR_SHEET_TITLE, "A1", "Дата")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "B1", "Событие")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "C1", "Сумма (₽)")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "D1", "Налог (₽)")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "E1", "Сумма ($)")
	f.SetCellValue(CALENDAR_SHEET_TITLE, "F1", "Налог ($)")

	lots, err := GetAccountPositions(uuid.MustParse("129274f9-ee80-4e74-aa1c-fea578bac6e6"))
	if err != nil {
		return err
	}

	lots, err = PopulateLotsWithBonds(lots)
	if err != nil {
		return err
	}

	lots = PopulateLotsWithCoupons(lots)

	// timeline, err := generateTimeLineForLots(lots, false)


	startWithRowIndex := 2 //The first row is reserved for headers
	for _, lot := range lots {
		startWithRowIndex = EnterLotInformationIntoSpreadsheet(f, lot, startWithRowIndex)
	}

	if err := f.SaveAs("Portfolio_Calendar.xlsx"); err != nil {
		logger.Log(err.Error(), logger.ERROR)
	}
	return nil
}

func EnterLotInformationIntoSpreadsheet(f *excelize.File, lot bonds.BondLot, currentRow int) int {
	totalRubPayout := 0.0
	totalUsdPayout := 0.0
	for _, coupon := range lot.Bond.Coupons {
		if coupon.CouponDate.Before(time.Now()) {
			continue
		}

		year, month, day := coupon.CouponDate.Date()
		formattedDate := strconv.Itoa(day) + "." + strconv.Itoa(int(month)) + "." + strconv.Itoa(year)
		f.SetCellValue(CALENDAR_SHEET_TITLE, "A"+strconv.Itoa(currentRow), formattedDate)
		f.SetCellValue(CALENDAR_SHEET_TITLE, "B"+strconv.Itoa(currentRow), "Выплата купона")

		if lot.Bond.IsRubleBond() {
			payout := lot.Quantity * coupon.PerBondAmount
			f.SetCellValue(CALENDAR_SHEET_TITLE, "C"+strconv.Itoa(currentRow), fmt.Sprint(payout))
			f.SetCellValue(CALENDAR_SHEET_TITLE, "D"+strconv.Itoa(currentRow), lot.Quantity*coupon.PerBondAmount*0.13)
			totalRubPayout += payout
		} else if lot.Bond.IsBondWithDifferentNominalCurrencyAndCurrency() {
			payout := lot.Quantity * coupon.PerBondAmount
			f.SetCellValue(CALENDAR_SHEET_TITLE, "E"+strconv.Itoa(currentRow), fmt.Sprint(payout))
			f.SetCellValue(CALENDAR_SHEET_TITLE, "F"+strconv.Itoa(currentRow), lot.Quantity*coupon.PerBondAmount*0.13)
			totalUsdPayout += payout
		} else {
			logger.Log("Unexpected scenario for lot "+lot.Figi, logger.ERROR)
		}

		currentRow++
	}

	return 0//FIXME
}
