package output

import (
	excel "github.com/360EntSecGroup-Skylar/excelize"
	"github.com/Newlooc/fundtools/pkg/apis"
	"github.com/Newlooc/fundtools/pkg/parser"
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Excel struct {
	Filename           string
	ExcelFile          *excel.File
	Sheet              string
	currentColume      string
	currentRow         string
	headStartColume    string
	headStartRow       string
	contentStartColume string
	contentStartRow    string
}

func NewExcel(filename, sheet string) *Excel {
	return &Excel{
		Filename:           filename,
		ExcelFile:          excel.NewFile(),
		Sheet:              sheet,
		headStartColume:    "A",
		headStartRow:       "1",
		contentStartColume: "B",
		contentStartRow:    "2",
	}
}

func (e *Excel) Write(data map[apis.URLConfig]*parser.DTMock, x []time.Time, y []time.Time) error {
	index := e.ExcelFile.NewSheet(e.Sheet)
	e.ExcelFile.SetActiveSheet(index)

	log.Info("Start writing xhead")
	e.writeXHead(x)

	log.Info("Start writing yhead")
	e.writeYHead(y)

	log.Info("Start writing content")
	e.PosContentStart()
	for _, startPoint := range y {
		for _, endPoint := range x {
			dtData := data[apis.URLConfig{
				Start: startPoint,
				End:   endPoint,
			}]
			if dtData == nil {
				e.Right()
				continue
			}
			log.Debugf("Set Cell Value: %s, %s, %+v", e.Sheet, e.GetAxis(), dtData)

			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.ProfitPercent)
			e.Right()
		}
		e.PosContentNextRow()
	}

	log.Info("Start writing rows")
	e.PosRawNextRow(20)
	for _, startPoint := range y {
		for _, endPoint := range x {
			dtData := data[apis.URLConfig{
				Start: startPoint,
				End:   endPoint,
			}]
			if dtData == nil {
				continue
			}
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), startPoint.Format(apis.DateFormat))
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), endPoint.Format(apis.DateFormat))
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.Code)
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.CNName)
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.InvestType)
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.ProfitPercent)
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.TotalRound)
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.BaseInvestAmount)
			e.Right()
			e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), dtData.FinalAmount)
			e.Right()
			e.PosRawNextRow(1)
		}
	}

	log.Info("Start saving file")
	if err := e.ExcelFile.SaveAs(e.Filename); err != nil {
		log.WithError(err).Error("Failed to save EXCEL file")
	}
	log.Infof("File saved at %s", e.Filename)

	return nil
}

func (e *Excel) writeXHead(x []time.Time) {
	e.PosHeadStart()
	e.Right()
	for _, timePoint := range x {
		e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), timePoint.Format(apis.DateFormat))
		e.Right()
	}
}

func (e *Excel) writeYHead(y []time.Time) {
	e.PosHeadStart()
	e.Down()
	for _, timePoint := range y {
		e.ExcelFile.SetCellValue(e.Sheet, e.GetAxis(), timePoint.Format(apis.DateFormat))
		e.Down()
	}
}

func (e *Excel) Right() {
	byteArr := []byte(e.currentColume)
	arrayRevers(byteArr)
	carry := 1
	i := 0

	for carry != 0 {
		if i == len(byteArr) {
			byteArr = append(byteArr, 'A')
			carry--
		} else {
			if byteArr[i] == 'Z' {
				byteArr[i] = 'A'
				i++
			} else {
				byteArr[i]++
				carry--
			}
		}
	}
	arrayRevers(byteArr)
	e.currentColume = string(byteArr)
}

func arrayRevers(arr []byte) {
	h := 0
	t := len(arr) - 1
	for h < t {
		temp := arr[h]
		arr[h] = arr[t]
		arr[t] = temp
		h++
		t--
	}
}

func (e *Excel) Down() {
	rowNum, _ := strconv.Atoi(e.currentRow)
	rowNum++
	e.currentRow = strconv.Itoa(rowNum)
}

func (e *Excel) PosContentStart() {
	e.currentColume = e.contentStartColume
	e.currentRow = e.contentStartRow
}

func (e *Excel) PosContentNextRow() {
	e.currentColume = e.contentStartColume
	e.Down()
}

func (e *Excel) PosHeadStart() {
	e.currentColume = e.headStartColume
	e.currentRow = e.headStartRow
}

func (e *Excel) GetAxis() string {
	return e.currentColume + e.currentRow
}

func (e *Excel) PosRawNextRow(rowCount int) {
	e.currentColume = "A"
	for i := 0; i < rowCount; i++ {
		e.Down()
	}
}
