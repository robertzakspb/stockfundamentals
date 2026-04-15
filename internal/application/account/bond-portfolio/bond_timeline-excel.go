package bondportfolio

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func GenerateTimeLineExcel() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	//index, err := f.NewSheet("Sheet2")
}
