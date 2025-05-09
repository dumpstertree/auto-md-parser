package main

import (
	"github.com/xuri/excelize/v2"
)

type PageParser struct {
	ParsedPages map[*excelize.File][]string
}

func (p *PageParser) Load(loader *DataLoader) {
	p.ParsedPages = buildPageLists(loader.LoadedSpreadsheets)
}

func buildPageLists(layoutToFile map[*OrderedLayout]*excelize.File) map[*excelize.File][]string {

	// create return value
	allPages := make(map[*excelize.File][]string)

	// add pages for each path
	for layout, file := range layoutToFile {

		// get sheets on excel file
		allSheets := file.GetSheetList()

		// create array to hold pages for this file
		pages := []string{}

		// add all include pages
		if len(layout.IncludeSheets) == 0 {
			pages = append(pages, allSheets...)
		} else {
			pages = append(pages, layout.IncludeSheets...)
		}

		// remove all excluded pages
		if len(layout.ExcludeSheets) != 0 {
			for _, y := range layout.ExcludeSheets {
				pages = remove(pages, y)
			}
		}
		// add pages for file to map
		allPages[file] = pages
	}
	// return value
	return allPages
}
