// file: /src/infrastructure/lib/excel/excel.go
package excel

import (
	"bytes"
	"mime/multipart"

	"github.com/xuri/excelize/v2"
)

type ExcelHandler struct{}

func NewExcelHandler() *ExcelHandler {
	return &ExcelHandler{}
}

// ExcelData 表示Excel数据结构
type ExcelData struct {
	Headers []string
	Rows    [][]string
}

// CreateExcel 创建Excel文件
func (e *ExcelHandler) CreateExcel(sheetName string, data *ExcelData) (*bytes.Buffer, error) {
	f := excelize.NewFile()

	// 创建或获取工作表
	if sheetName != "" {
		f.SetSheetName("Sheet1", sheetName)
	} else {
		sheetName = "Sheet1"
	}

	// 写入表头
	for i, header := range data.Headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	// 写入数据行
	for rowIndex, row := range data.Rows {
		for colIndex, cellValue := range row {
			cell, _ := excelize.CoordinatesToCellName(colIndex+1, rowIndex+2)
			f.SetCellValue(sheetName, cell, cellValue)
		}
	}

	// 将文件写入缓冲区
	buffer := &bytes.Buffer{}
	if err := f.Write(buffer); err != nil {
		return nil, err
	}

	return buffer, nil
}

// ReadExcel 读取Excel文件
func (e *ExcelHandler) ReadExcel(file multipart.File, sheetName string) (*ExcelData, error) {
	// 读取Excel文件
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 获取工作表名称
	if sheetName == "" {
		sheetName = f.GetSheetName(0) // 获取第一个工作表
	}

	// 读取所有行
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return &ExcelData{}, nil
	}

	// 第一行作为表头
	headers := rows[0]

	// 其余行作为数据
	var dataRows [][]string
	if len(rows) > 1 {
		dataRows = rows[1:]
	}

	return &ExcelData{
		Headers: headers,
		Rows:    dataRows,
	}, nil
}

// CreateTemplate 创建API模板
func (e *ExcelHandler) CreateApiTemplate() (*bytes.Buffer, error) {
	templateData := &ExcelData{
		Headers: []string{"ID", "Path", "ApiGroup", "Method", "Description"},
		Rows:    [][]string{},
	}

	return e.CreateExcel("APIs", templateData)
}
