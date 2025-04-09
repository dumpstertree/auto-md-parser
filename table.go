package main

import (
	"github.com/tzfqh/gmdtable"
	"github.com/xuri/excelize/v2"
)

const PREFIX_TABLE = ""
const SUFFIX_TABLE = "\n"

type Table struct {
	RowMin  int            `json:"rowmin"`
	RowMax  int            `json:"rowmax"`
	Content []TableCollumn `json:"content"`
	BaseSubsection
}

type TableCollumn struct {
	Header  string `json:"header"`
	Content string `json:"content"`
}

func (re Table) Write(input string, sheet string, allPages []string, file *excelize.File) string {

	// failed to get max rows
	rows, err := file.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	// if <=0 set to min
	min := re.RowMin
	if min <= 0 {
		min = 1
	}

	// if <=0 set to max rows
	max := re.RowMax
	if max <= 0 {
		max = len(rows)
	}

	// create header and data
	headers := getHeaders(re.Content)
	data := getData(re.Content, sheet, allPages, min, max, file)

	// get table from external package
	table, err := gmdtable.Convert(headers, data)
	if err != nil {
		panic(err)
	}

	// return
	return input + PREFIX_TABLE + table + SUFFIX_TABLE
}

func getHeaders(content []TableCollumn) []string {
	headers := []string{}
	for _, x := range content {
		headers = append(headers, x.Header)
	}
	return headers
}
func getData(content []TableCollumn, sheet string, allPages []string, min int, max int, file *excelize.File) []map[string]interface{} {

	// add each for data in collumn
	data := []map[string]interface{}{}
	for i := min; i <= max; i++ {

		// create a map for all the headers : parsed value for row
		m := make(map[string]interface{})
		for _, x := range content {

			// assign after parsing string
			m[x.Header] = parseCompoundCollumnString(x.Content, sheet, i, allPages, file)
		}
		// add map to list
		data = append(data, m)
	}
	return data
}
