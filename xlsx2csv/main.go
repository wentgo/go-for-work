package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path"
	"encoding/csv"

	"github.com/tealeg/xlsx"
)

var xlsxPath = flag.String("f", "", "Path to an XLSX file")
var csvPath = flag.String("o", "", "Output CSV filename")
var sheetIndex = flag.Int("i", 0, "Index of sheet to convert, zero based")
var delimiter = flag.String("d", ",", "Delimiter to use between fields")

func xlsx2csv(xlsxFilename string, csvFilename string, sheetIndex int) (err error) {
	xlsFile, err := xlsx.OpenFile(xlsxFilename)
	if err != nil {
		return
	}

	sheetLen := len(xlsFile.Sheets)
	switch {
	case sheetLen == 0:
		return errors.New("This XLSX file contains no sheets.")
	case sheetIndex >= sheetLen:
		return fmt.Errorf("No sheet %d available, please select a sheet between 0 and %d\n", sheetIndex, sheetLen-1)
	}

	csvFile, err := os.Create(csvFilename)
	if err != nil {
		return
	}
	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)
	writer.Comma = []rune(*delimiter)[0]
	writer.UseCRLF = true
	if *delimiter == "\\t" || *delimiter == "TAB" {
		writer.Comma = '\t'
	}
	defer writer.Flush()

	sheet := xlsFile.Sheets[sheetIndex]
	for _, row := range sheet.Rows {
		var vals []string
		if row != nil {
			for _, cell := range row.Cells {
				str, err := cell.FormattedValue()
				if err != nil {
					vals = append(vals, err.Error())
				}
				vals = append(vals, str)
			}
			err = writer.Write(vals)
		}
	}

	return nil
}

func main() {
	flag.Parse()

	if len(os.Args) < 3 {
		flag.PrintDefaults()
		return
	}

	ext := path.Ext(*xlsxPath)
	if *csvPath == "" {
		*csvPath = (*xlsxPath)[0:len(*xlsxPath)-len(ext)] + ".csv"
	}

	if err := xlsx2csv(*xlsxPath, *csvPath, *sheetIndex); err != nil {
		fmt.Println(err)
	}
}
