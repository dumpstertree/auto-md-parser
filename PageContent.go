package main

import (
	"github.com/xuri/excelize/v2"
)

type PageContent struct {
	Content map[string][]string
	Pages   []Page
}

func makePageContent(d *DataLoader, p *PageParser) *PageContent {
	x := new(PageContent)
	x.Pages = buildPageContent(d.LoadedSpreadsheets, p.ParsedPages)
	return x
}

// func buildPageContent(layoutToFile map[*OrderedLayout]*excelize.File, allPages map[*excelize.File][]string) map[string][]string {
func buildPageContent(layoutToFile map[*OrderedLayout]*excelize.File, allPages map[*excelize.File][]string) []Page {

	// generate all pages
	pages := []Page{}

	// add pages for each path
	for layout, file := range layoutToFile {

		// get sheets for this file
		reqSheets := allPages[file]

		//iterate on the required pages
		for _, sheet := range reqSheets {

			// make the current page
			p := *makePage(
				layout.Path+"/",
				sheet,
				"",
				layout.URL,
				layout.Tags,
			)

			// add all subsections
			for _, subsection := range layout.LayoutSubsection {
				pages = subsection.Write(p, pages, sheet, file)
			}

			// add page
			pages = append(pages, p)
		}

	}

	// return value
	return pages
}
