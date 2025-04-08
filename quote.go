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

func (q Quote) Write(input string, sheet string, allPages []string, file *excelize.File) string {

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
	input = q.ModifyTextStart(input)

	for i := min; i <= max; i++ {
		// iterate over each entry adding it to the quote
		input += PREFIX_QUOTE + parseCompoundCollumnString(q.Content, sheet, i, allPages, file) + SUFFIX_QUOTE

		// if break add a new line
		if q.Break {
			input += "\n"
		}
	}
	// modify base
	input = q.ModifyTextEnds(input)

	// return
	return input
}
