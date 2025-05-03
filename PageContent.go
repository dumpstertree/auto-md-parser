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

	// flatten the pages for use when writing
	allPagesFlat := []string{}
	for _, y := range allPages {
		allPagesFlat = append(allPagesFlat, y...)
	}

	// add pages for each path
	for layout, file := range layoutToFile {

		// get sheets for this file
		reqSheets := allPages[file]

		//iterate on the required pages
		for _, sheet := range reqSheets {

			// page content
			var content = ""

			// add all subsections
			for _, subsection := range layout.LayoutSubsection {
				content = subsection.Write(content, sheet, allPagesFlat, file)
			}

			// add page
			pages = append(pages,
				*makePage(
					layout.Path+"/",
					sheet,
					content,
					layout.URL,
					layout.Tags,
				),
			)
		}

	}

	// return value
	return pages
}
