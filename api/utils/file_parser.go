package utils

import (
	"encoding/csv"
	"mime/multipart"
	"strings"

	"github.com/xuri/excelize/v2"
)

// CSV
func ParseCSV(file multipart.File) ([][]string, error) {
	reader := csv.NewReader(file)
	return reader.ReadAll()
}

// EXCEL
func ParseExcel(file multipart.File) ([][]string, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return nil, err
	}

	sheet := f.GetSheetName(0)
	return f.GetRows(sheet)
}

// AUTO DETECT
func ParseFile(file multipart.File, filename string) ([][]string, error) {
	if strings.HasSuffix(filename, ".csv") {
		return ParseCSV(file)
	}
	return ParseExcel(file)
}