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

func (h Header) Write(page Page, allPages []Page, sheet string, file *excelize.File) []Page {
	// modify base
	page.Content = h.ModifyTextStart(page.Content)
	page.Content += PREFIX_HEADER + h.Content + SUFFIX_HEADER
	page.Content = h.ModifyTextEnds(page.Content)
	return allPages
}
