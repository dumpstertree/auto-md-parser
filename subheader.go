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

func (h Subheader) Write(page Page, allPages []Page, sheet string, file *excelize.File) []Page {
	page.Content = h.ModifyTextStart(page.Content)
	page.Content += PREFIX_SUBHEADER + h.Content + SUFFIX_SUBHEADER
	page.Content = h.ModifyTextEnds(page.Content)
	page.Content += "\n"
	return allPages
}
