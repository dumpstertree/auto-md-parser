package main

import (
	"github.com/xuri/excelize/v2"
)

const PREFIX_LIST_UNORDERED = "<ul>"
const SUFFIX_LIST_UNORDERED = "</ul>\n"
const PREFIX_LIST_UNORDERED_ENTRY = "<li>"
const SUFFIX_LIST_UNORDERED_ENTRY = "</li>\n"

type UnorderedList struct {
	RowMin  int    `json:"rowmin"`
	RowMax  int    `json:"rowmax"`
	Content string `json:"content"`
	BaseSubsection
	TextSubsection
}

func (l UnorderedList) Write(input string, sheet string, allPages []string, file *excelize.File) string {

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
	input += PREFIX_LIST_UNORDERED
	for i := min; i <= max; i++ {
		// iterate over each entry adding it to the list
		input += PREFIX_LIST_UNORDERED_ENTRY + parseCompoundCollumnString(l.Content, sheet, i, allPages, file) + SUFFIX_LIST_UNORDERED_ENTRY
	}
	// end list
	input += SUFFIX_LIST_UNORDERED
	input = l.ModifyTextEnds(input)

	// return
	return input
}
