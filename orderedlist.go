package main

import (
	"github.com/xuri/excelize/v2"
)

const PREFIX_LIST_ORDERED = "<ol>"
const SUFFIX_LIST_ORDERED = "</ol>\n"
const PREFIX_LIST_ORDERED_ENTRY = "<li>"
const SUFFIX_LIST_ORDERED_ENTRY = "</li>\n"

type OrderedList struct {
	RowMin  int    `json:"rowmin"`
	RowMax  int    `json:"rowmax"`
	Content string `json:"content"`
	BaseSubsection
	TextSubsection
}

func (l OrderedList) Write(input string, sheet string, allPages []string, file *excelize.File) string {

	// failed to get max rows
	rows, err := file.GetRows(sheet)
	if err != nil {
		panic(err)
	}

	// if <=0 set to min
	min := l.RowMin
	if min <= 0 {
		min = 1
	}

	// if <=0 set to max rows
	max := l.RowMax
	if max <= 0 {
		max = len(rows)
	}

	// begin list
	input = l.ModifyTextStart(input)
	input += PREFIX_LIST_ORDERED
	for i := min; i <= max; i++ {
		// iterate over each entry adding it to the list
		input += PREFIX_LIST_ORDERED_ENTRY + parseCompoundCollumnString(l.Content, sheet, i, allPages, file) + SUFFIX_LIST_ORDERED_ENTRY
	}
	// end list
	input += SUFFIX_LIST_ORDERED
	input = l.ModifyTextEnds(input)

	// return
	return input
}
