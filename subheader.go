package main

import (
	"github.com/xuri/excelize/v2"
)

const PREFIX_SUBHEADER = "## "
const SUFFIX_SUBHEADER = "\n"

type Subheader struct {
	Content string `json:"content"`
	BaseSubsection
	TextSubsection
}

func (h Subheader) Write(input string, sheet string, allPages []string, file *excelize.File) string {
	input = h.ModifyTextStart(input)
	input += PREFIX_SUBHEADER + h.Content + SUFFIX_SUBHEADER
	input = h.ModifyTextEnds(input)
	return input
}
