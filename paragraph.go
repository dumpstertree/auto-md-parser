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

func (re Paragraph) Write(input string, sheet string, allPages []string, file *excelize.File) string {

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
	input = re.ModifyTextStart(input)
	for i := min; i <= max; i++ {
		input += PREFIX_PARAGRAPH + parseCompoundCollumnString(re.Content, sheet, i, allPages, file) + SUFFIX_PARAGRAPH
	}
	input = re.ModifyTextEnds(input)

	// return
	return input
}
