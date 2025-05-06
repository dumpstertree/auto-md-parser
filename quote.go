package main

import (
	"github.com/xuri/excelize/v2"
)

// const PREFIX_QUOTE = ">"
// const SUFFIX_QUOTE = "\n"
const PREFIX_QUOTE = ""
const SUFFIX_QUOTE = ""

type Quote struct {
	Break   bool   `json:"break"`
	RowMin  int    `json:"rowmin"`
	RowMax  int    `json:"rowmax"`
	Content string `json:"content"`
	BaseSubsection
	TextSubsection
}

func (q Quote) Write(page Page, allPages []Page, sheet string, file *excelize.File) []Page {

	// failed to get max rows
	rows, err := file.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	// if <=0 set to min
	min := q.RowMin
	if min <= 0 {
		min = 1
	}

	// if <=0 set to max rows
	max := q.RowMax
	if max <= 0 {
		max = len(rows)
	}

	// modify base
	page.Content = q.ModifyTextStart(page.Content)

	for i := min; i <= max; i++ {
		// iterate over each entry adding it to the quote
		page.Content += PREFIX_QUOTE + parseCompoundCollumnString(q.Content, sheet, i, file) + SUFFIX_QUOTE

		// if break add a new line
		if q.Break {
			page.Content += "\n"
		}
	}
	// modify base
	page.Content = q.ModifyTextEnds(page.Content)

	// return
	return allPages
}
