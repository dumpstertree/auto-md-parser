package main

import (
	"github.com/xuri/excelize/v2"
)

const PREFIX_PARAGRAPH = ""
const SUFFIX_PARAGRAPH = "\n"

type Paragraph struct {
	RowMin  int    `json:"rowmin"`
	RowMax  int    `json:"rowmax"`
	Content string `json:"content"`
	BaseSubsection
	TextSubsection
}

func (re Paragraph) Write(page Page, allPages []Page, sheet string, file *excelize.File) []Page {

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

	//
	page.Content = re.ModifyTextStart(page.Content)
	for i := min; i <= max; i++ {
		page.Content += PREFIX_PARAGRAPH + parseCompoundCollumnString(re.Content, sheet, i, file) + SUFFIX_PARAGRAPH
	}
	page.Content = re.ModifyTextEnds(page.Content)
	page.Content += "\n"

	// return
	return allPages
}
