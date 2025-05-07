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

func (l OrderedList) Write(page *Page, allPages []Page, sheet string, file *excelize.File) []Page {

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
	page.Content = l.ModifyTextStart(page.Content)
	page.Content += PREFIX_LIST_ORDERED
	for i := min; i <= max; i++ {
		// iterate over each entry adding it to the list
		page.Content += PREFIX_LIST_ORDERED_ENTRY + parseCompoundCollumnString(l.Content, sheet, i, file) + SUFFIX_LIST_ORDERED_ENTRY
	}
	// end list
	page.Content += SUFFIX_LIST_ORDERED
	page.Content = l.ModifyTextEnds(page.Content)

	// return
	return allPages
}
