package main

import (
	"strings"

	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

type PageSummary struct {
	Pages []Page
}

func makePageSummary(pages []Page) *PageSummary {
	x := new(PageSummary)
	x.Pages = buildPageSumary(pages)
	//x.Pages = buildPageSumary(d.LoadedSpreadsheets, p.ParsedPages)
	return x
}

// func buildPageContent(layoutToFile map[*OrderedLayout]*excelize.File, allPages map[*excelize.File][]string) map[string][]string {
func buildPageSumary(pages []Page) []Page {

	// get all paths
	paths := []string{}
	pathForName := make(map[string]string)
	for _, k := range pages {
		p := k.Path + k.DisplayName
		paths = append(paths, p)
		pathForName[p] = k.LinkName
	}

	// sort
	c := collate.New(language.English, collate.IgnoreCase)
	c.SortStrings(paths)

	// entry line - always empty
	content := "#\n"

	// iterate over each path
	for pathIndex, path := range paths {

		// split this path
		split := strings.Split(path, "/")

		// iterate over each subpath
		for subPathIndex, subPath := range split {

			if pathIndex > 0 {

				//
				splitlast := strings.Split(paths[pathIndex-1], "/")

				if len(splitlast) > subPathIndex && splitlast[subPathIndex] == subPath {
					continue
				}
			}
			//
			if subPathIndex == 0 {
				// if this is the first level make it a header
				content += "# " + subPath + "\n"

			} else if subPathIndex == len(split)-1 {

				// add space
				for f := 1; f < subPathIndex; f++ {
					content += "	"
				}

				// add link
				content += "- [" + subPath + "](" + pathForName[path] + ".md)\n"

			} else {

				// add space
				for f := 1; f < subPathIndex; f++ {
					content += "	"
				}

				// add draft
				content += "- [" + subPath + "]()\n"
			}
		}
	}
	return []Page{
		*makePageExplicit("", "SUMMARY", "SUMMARY", content, "", nil),
	}
}
