package main

import (
	"encoding/json"

	"github.com/xuri/excelize/v2"
)

const PREFIX_LIST_SUBPAGE = "<ol>"
const SUFFIX_LIST_SUBPAGE = "</ol>\n"
const PREFIX_LIST_SUBPAGE_ENTRY = "<li>"
const SUFFIX_LIST_SUBPAGE_ENTRY = "</li>\n"

type Subpage struct {
	RowMin           int               `json:"rowmin"`
	RowMax           int               `json:"rowmax"`
	Title            string            `json:"title"`
	LayoutSubsection []ISubSubsection  `json:"-"`
	Sections         []json.RawMessage `json:"sections"`
	BaseSubsection
	TextSubsection
}

func (l Subpage) Write(page *Page, allPages []Page, sheet string, file *excelize.File) []Page {

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

	page.Content += "<ul>\n"
	for i := min; i <= max; i++ {

		// parse for title
		title := parseCompoundCollumnString(l.Title, sheet, i, file)

		// add link from this page to next
		page.Content += "<li>" + title + "</li>\n"

		// make subpage
		subpage := *makePage(
			page.Path+page.DisplayName+"/",
			title,
			"",
			page.Source,
			page.OverrideSource,
			nil,
		)

		// write subpages
		for _, x := range l.LayoutSubsection {
			x.WriteSubsection(&subpage, allPages, sheet, file, i)
		}

		// add new page
		allPages = append(allPages, subpage)
	}
	page.Content += "</ul>\n"

	// add space
	page.Content += "\n"

	// return
	return allPages
}
