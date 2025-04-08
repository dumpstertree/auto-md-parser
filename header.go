package main

import (
	"github.com/xuri/excelize/v2"
)

const PREFIX_HEADER = "## "
const SUFFIX_HEADER = "\n"

type Header struct {
	Content string `json:"content"`
	BaseSubsection
	TextSubsection
}

func (h Header) Write(input string, sheet string, allPages []string, file *excelize.File) string {
	// modify base
	input = h.ModifyTextStart(input)
	input += PREFIX_HEADER + h.Content + SUFFIX_HEADER
	input = h.ModifyTextEnds(input)
	return input
}
