package xlsx

import (
	"bytes"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/pulpfree/gdps-fs-dwnld/model"

	log "github.com/sirupsen/logrus"
)

// XLSX struct
type XLSX struct {
	file *excelize.File
}

// Defaults
const (
	abc             = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	floatFrmt       = "#,#0"
	timeShortForm   = "20060102"
	timeMonthForm   = "200601"
	dateDayFormat   = "Jan _2"
	dateMonthFormat = "January 2006"
)

// NewFile function
func NewFile() (x *XLSX, err error) {

	x = new(XLSX)
	x.file = excelize.NewFile()
	if err != nil {
		log.Errorf("xlsx err %s: ", err)
	}
	return x, err
}

// FuelSales method
func (x *XLSX) FuelSales(fs *model.FuelSales) (err error) {

	var cell string
	var style int

	xlsx := x.file
	sheetNm := "Sheet1"

	xlsx.SetSheetName(sheetNm, "Fuel Sales")

	// Merge cells to accommodate width of all fuel types
	endCell := toChar(len(fs.FuelTypes)+2) + "1"
	xlsx.MergeCell(sheetNm, "A1", endCell)

	style, _ = xlsx.NewStyle(`{"font":{"bold":true,"size":12}}`)

	title := fmt.Sprintf("%s Fuel Sales Detail - %s", fs.Station.Name, fs.Date.Format("January 2006"))
	xlsx.SetCellValue(sheetNm, "A1", title)
	xlsx.SetCellStyle(sheetNm, "A1", "A1", style)

	// Create second row with fuel type headings
	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellStyle(sheetNm, "A2", "A2", style)

	col := 2
	row := 2
	style, _ = xlsx.NewStyle(`{"font":{"bold":true}}`)
	for _, ft := range fs.FuelTypes {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, ft)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}

	// Fill in data
	col = 1
	row = 3
	style, _ = xlsx.NewStyle(`{"number_format": 3}`)

	for _, r := range fs.Report.StationSales {

		t, _ := time.Parse(timeShortForm, strconv.Itoa(int(r.Date)))
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, t.Format(dateDayFormat))
		col++

		for _, ft := range fs.FuelTypes {
			cell = toChar(col) + strconv.Itoa(row)
			xlsx.SetCellValue(sheetNm, cell, r.Sales[ft])
			xlsx.SetCellStyle(sheetNm, cell, cell, style)
			col++
		}
		col = 1
		row++
	}

	// Fueltype summary
	style, _ = xlsx.NewStyle(`{"number_format": 3, "font":{"bold":true}}`)
	cell = toChar(col) + strconv.Itoa(row)
	xlsx.SetCellValue(sheetNm, cell, "")
	col++

	for _, ft := range fs.FuelTypes {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, fs.Report.SalesSummary[ft])
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}
	row += 2
	col = 1
	cell = toChar(col) + strconv.Itoa(row)
	cellNext := toChar(col+1) + strconv.Itoa(row)
	xlsx.MergeCell(sheetNm, cell, cellNext)
	xlsx.SetCellValue(sheetNm, cell, "Total Sales")
	xlsx.SetCellStyle(sheetNm, cell, cell, style)

	col += 2
	cell = toChar(col) + strconv.Itoa(row)
	cellNext = toChar(col+1) + strconv.Itoa(row)
	xlsx.MergeCell(sheetNm, cell, cellNext)
	xlsx.SetCellValue(sheetNm, cell, fs.Report.SalesTotal)
	xlsx.SetCellStyle(sheetNm, cell, cell, style)

	return err
}

// FuelDelivery method
func (x *XLSX) FuelDelivery(fd *model.FuelDelivery) (err error) {

	var cell string
	var style int
	fuelTypes := fd.Report.FuelTypes
	numColWidth := 10.00

	xlsx := x.file
	sheetNm := "Sheet2"

	_ = xlsx.NewSheet(sheetNm)

	xlsx.SetSheetName(sheetNm, "Fuel Delivery")

	// Merge cells to accommodate width of all fuel types
	endCell := toChar(len(fuelTypes)+2) + "1"
	xlsx.MergeCell(sheetNm, "A1", endCell)

	style, _ = xlsx.NewStyle(`{"font":{"bold":true,"size":12}}`)

	title := fmt.Sprintf("%s Fuel Deliveries - %s", fd.Station.Name, fd.Date.Format("January 2006"))
	xlsx.SetCellValue(sheetNm, "A1", title)
	xlsx.SetCellStyle(sheetNm, "A1", "A1", style)

	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellStyle(sheetNm, "A2", "A2", style)

	// Create second row with fuel type headings
	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellStyle(sheetNm, "A2", "A2", style)

	xlsx.SetColWidth(sheetNm, "B", toChar(len(fuelTypes)+1), numColWidth)

	col := 2
	row := 2
	style, _ = xlsx.NewStyle(`{"font":{"bold":true}}`)
	for _, ft := range fuelTypes {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, ft)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}

	// Fill in data
	col = 1
	row = 3
	style, _ = xlsx.NewStyle(`{"number_format": 3}`)

	for _, r := range fd.Report.Deliveries {

		t, _ := time.Parse(timeShortForm, strconv.Itoa(int(r.Date)))
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, t.Format(dateDayFormat))
		col++

		for _, ft := range fuelTypes {

			cell = toChar(col) + strconv.Itoa(row)
			if r.Data[ft] > 0 {
				xlsx.SetCellValue(sheetNm, cell, r.Data[ft])
			} else {
				xlsx.SetCellValue(sheetNm, cell, "")
			}
			xlsx.SetCellStyle(sheetNm, cell, cell, style)
			col++
		}

		col = 1
		row++
	}

	// Summary Row
	style, _ = xlsx.NewStyle(`{"number_format": 3, "font":{"bold":true}}`)
	cell = toChar(col) + strconv.Itoa(row)
	xlsx.SetCellValue(sheetNm, cell, "")
	col++

	for _, ft := range fuelTypes {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, fd.Report.DeliverySummary[ft])
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}

	return err
}

// OverShortMonth method
func (x *XLSX) OverShortMonth(os *model.OverShortMonth) (err error) {

	var cell string
	var style int
	fuelTypes := os.Report.FuelTypes
	numColWidth := 10.00

	xlsx := x.file
	sheetNm := "Sheet3"

	_ = xlsx.NewSheet(sheetNm)

	xlsx.SetSheetName(sheetNm, "Over-Short Month")

	// Merge cells to accommodate width of all fuel types
	endCell := toChar(len(fuelTypes)+2) + "1"
	xlsx.MergeCell(sheetNm, "A1", endCell)

	style, _ = xlsx.NewStyle(`{"font":{"bold":true,"size":12}}`)

	title := fmt.Sprintf("%s Over-Short Month - %s", os.Station.Name, os.Date.Format("January 2006"))
	xlsx.SetCellValue(sheetNm, "A1", title)
	xlsx.SetCellStyle(sheetNm, "A1", "A1", style)

	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellStyle(sheetNm, "A2", "A2", style)

	// Create second row with fuel type headings
	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellStyle(sheetNm, "A2", "A2", style)

	xlsx.SetColWidth(sheetNm, "B", toChar(len(fuelTypes)+1), numColWidth)

	col := 2
	row := 2
	style, _ = xlsx.NewStyle(`{"font":{"bold":true}}`)
	for _, ft := range fuelTypes {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, ft)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}

	// Fill in data
	col = 1
	row = 3
	stylePos, _ := xlsx.NewStyle(`{"number_format": 4}`)
	styleNeg, _ := xlsx.NewStyle(`{"number_format": 4, "font":{"color": "#ff0000"}}`)

	for _, r := range os.Report.OverShort {

		t, _ := time.Parse(timeShortForm, strconv.Itoa(int(r.Date)))
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, t.Format(dateDayFormat))
		col++

		for _, ft := range fuelTypes {
			val := r.Data[ft].OverShort
			if val < 0 {
				style = styleNeg
			} else {
				style = stylePos
			}
			cell = toChar(col) + strconv.Itoa(row)
			xlsx.SetCellValue(sheetNm, cell, val)
			xlsx.SetCellStyle(sheetNm, cell, cell, style)
			col++
		}

		col = 1
		row++
	}

	// Summary Row
	stylePos, _ = xlsx.NewStyle(`{"number_format": 4, "font": {"bold":true}}`)
	styleNeg, _ = xlsx.NewStyle(`{"number_format": 4, "font":{"bold":true, "color": "#ff0000"}}`)

	cell = toChar(col) + strconv.Itoa(row)
	xlsx.SetCellValue(sheetNm, cell, "")
	col++

	for _, ft := range fuelTypes {

		val := os.Report.OverShortSummary[ft]
		if val < 0 {
			style = styleNeg
		} else {
			style = stylePos
		}
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, val)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}

	return err
}

// OverShortAnnual method
func (x *XLSX) OverShortAnnual(os *model.OverShortAnnual) (err error) {

	var cell string
	var style int
	fuelTypes := os.Report.FuelTypes
	numColWidth := 10.00
	months := setMonths(os.Report.Year, len(os.Report.Months))

	xlsx := x.file
	sheetNm := "Sheet4"

	_ = xlsx.NewSheet(sheetNm)

	xlsx.SetSheetName(sheetNm, "Over-Short Annual")

	// Merge cells to accommodate width of all fuel types
	endCell := toChar(len(fuelTypes)+2) + "1"
	xlsx.MergeCell(sheetNm, "A1", endCell)

	style, _ = xlsx.NewStyle(`{"font":{"bold":true,"size":12}}`)

	title := fmt.Sprintf("%s Over-Short Annual - %s", os.Station.Name, os.Date.Format("January 2006"))
	xlsx.SetCellValue(sheetNm, "A1", title)
	xlsx.SetCellStyle(sheetNm, "A1", "A1", style)

	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellStyle(sheetNm, "A2", "A2", style)

	// Create second row with fuel type headings
	xlsx.SetCellValue(sheetNm, "A2", "Date")
	xlsx.SetCellStyle(sheetNm, "A2", "A2", style)

	xlsx.SetColWidth(sheetNm, "A", toChar(len(fuelTypes)+1), numColWidth)

	col := 2
	row := 2
	style, _ = xlsx.NewStyle(`{"font":{"bold":true}}`)

	for _, ft := range fuelTypes {
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, ft)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}

	// Fill in data
	col = 1
	row = 3
	stylePos, _ := xlsx.NewStyle(`{"number_format": 4}`)
	styleNeg, _ := xlsx.NewStyle(`{"number_format": 4, "font":{"color": "#ff0000"}}`)

	for _, m := range months {

		t, _ := time.Parse(timeMonthForm, m)
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, t.Format("January"))
		col++

		for _, ft := range fuelTypes {

			val := os.Report.Months[m][ft]
			if val < 0 {
				style = styleNeg
			} else {
				style = stylePos
			}
			cell = toChar(col) + strconv.Itoa(row)
			xlsx.SetCellValue(sheetNm, cell, val)
			xlsx.SetCellStyle(sheetNm, cell, cell, style)
			col++
		}

		col = 1
		row++
	}

	// Summary Row
	stylePos, _ = xlsx.NewStyle(`{"number_format": 4, "font": {"bold":true}}`)
	styleNeg, _ = xlsx.NewStyle(`{"number_format": 4, "font":{"bold":true, "color": "#ff0000"}}`)

	cell = toChar(col) + strconv.Itoa(row)
	xlsx.SetCellValue(sheetNm, cell, "")
	col++

	for _, ft := range fuelTypes {
		val := os.Report.Summary[ft]
		if val < 0 {
			style = styleNeg
		} else {
			style = stylePos
		}
		cell = toChar(col) + strconv.Itoa(row)
		xlsx.SetCellValue(sheetNm, cell, val)
		xlsx.SetCellStyle(sheetNm, cell, cell, style)
		col++
	}

	return err
}

// OutputFile method
func (x *XLSX) OutputFile() (buf bytes.Buffer, err error) {
	err = x.file.Write(&buf)
	if err != nil {
		log.Errorf("xlsx err: %s", err)
	}
	return buf, err
}

// OutputToDisk method
func (x *XLSX) OutputToDisk(path string) (fp string, err error) {
	err = x.file.SaveAs(path)
	return path, err
}

// ======================== Helper Methods ================================= //

// see: https://stackoverflow.com/questions/36803999/golang-alphabetic-representation-of-a-number
// for a way to map int to letters
func toChar(i int) string {
	return abc[i-1 : i]
}

// Found these function at: https://stackoverflow.com/questions/18390266/how-can-we-truncate-float64-type-to-a-particular-precision-in-golang
// Looks like a good way to deal with precision
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func setMonths(year, numMonths int) (months []string) {
	dte := time.Date(year, time.January, 1, 12, 0, 0, 0, time.UTC)
	months = append(months, dte.Format("200601"))
	for n := 1; n < numMonths; n++ {
		nextMn := dte.AddDate(0, n, 0)
		months = append(months, nextMn.Format("200601"))
	}
	return months
}
