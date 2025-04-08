package main

import (
	"github.com/xuri/excelize/v2"
)

const PREFIX_CODE = "`"
const SUFFIX_CODE = "`\n"

type Code struct {
	Break   bool   `json:"break"`
	RowMin  int    `json:"rowmin"`
	RowMax  int    `json:"rowmax"`
	Content string `json:"content"`
	BaseSubsection
}

func (c Code) Write(input string, sheet string, allPages []string, file *excelize.File) string {

	// failed to get max rows
	rows, err := file.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	// if <=0 set to min
	min := c.RowMin
	if min <= 0 {
		min = 1
	}

	// if <=0 set to max rows
	max := c.RowMax
	if max <= 0 {
		max = len(rows)
	}

	for i := min; i <= max; i++ {
		// iterate over each entry adding it to the quote
		input += PREFIX_CODE + parseCompoundCollumnString(c.Content, sheet, i, allPages, file) + SUFFIX_CODE

		// if break add a new line
		if c.Break {
			input += "\n"
		}
	}

	// return
	return input
}
