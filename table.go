package main

import (
	"fmt"

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

func (re Table) Write(page *Page, allPages []Page, sheet string, file *excelize.File) []Page {

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
	data := getData(re.Content, sheet, min, max, file)

	// // get table from external package
	// table, err := gmdtable.Convert(headers, data)
	// if err != nil {
	// 	panic(err)
	// }

	table := getTable(headers, data)

	// return
	page.Content += PREFIX_TABLE + table + SUFFIX_TABLE
	return allPages
}

func getHeaders(content []TableCollumn) []string {
	headers := []string{}
	for _, x := range content {
		headers = append(headers, x.Header)
	}
	return headers
}
func getData(content []TableCollumn, sheet string, min int, max int, file *excelize.File) []map[string]string {

	// add each for data in collumn
	data := []map[string]string{}
	for i := min; i <= max; i++ {

		// create a map for all the headers : parsed value for row
		m := make(map[string]string)
		for _, x := range content {

			// assign after parsing string
			m[x.Header] = parseCompoundCollumnString(x.Content, sheet, i, file)
		}
		// add map to list
		data = append(data, m)
	}
	return data
}
func getTable(header []string, data []map[string]string) string {
	fmt.Println("start make table")
	content := ""

	content += "<div class='table-container'><table><thead>\n"

	// headers
	content += "<tr>\n"
	for _, h := range header {
		content += "<th>" + h + "</th>\n"
	}
	content += "</tr>\n"

	// body
	content += "</thead><tbody>"
	for _, d := range data {
		content += "<tr>\n"
		for _, h := range header {
			content += "<td>" + d[h] + "</td>\n"
		}
		content += "</tr>\n"
	}

	// end
	content += "</tbody></table></div>\n"
	content += "\n"
	fmt.Println("finish make table")
	return content
}
