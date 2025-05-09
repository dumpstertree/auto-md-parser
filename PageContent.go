package main

import (
	"strings"

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

			x := ""
			x += layout.Path
			if !strings.HasSuffix(x, "/") {
				x += "/"
			}

			x = x + strings.Replace(sheet, "|", "/", -1)
			y := strings.Split(x, "/")

			path := ""
			for i, r := range y {
				last := i == len(y)-1
				if last {
					break
				}
				path += r + "/"
			}

			// make the current page
			curPage := *makePage(
				path,
				y[len(y)-1],
				"",
				layout.URL,
				layout.Tags,
			)

			// add all subsections
			for _, subsection := range layout.LayoutSubsection {
				pages = subsection.Write(&curPage, pages, sheet, file)
			}

			// add page
			pages = append(pages, curPage)
		}

	}

	// return value
	return pages
}
