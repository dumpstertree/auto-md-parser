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
	pathsA := []string{}
	pathsB := []string{}
	pathForName := make(map[string]string)
	for _, k := range pages {
		p := k.Path + k.DisplayName
		if strings.Contains(k.Path, "Tag") {
			pathsB = append(pathsB, p)

		} else {
			pathsA = append(pathsA, p)
		}
		pathForName[p] = k.LinkName
	}

	// sort
	c := collate.New(language.English, collate.IgnoreCase)
	c.SortStrings(pathsA)

	// entry line - always empty
	content := "#\n"

	// write all content
	content = Write(content, pathsA, pathForName)
	content = Write(content, pathsB, pathForName)

	// return new page
	return []Page{
		*makePageExplicit("", "SUMMARY", "SUMMARY", content, "", "", nil),
	}
}

func Write(content string, pathsA []string, pathForName map[string]string) string {
	// iterate over each path
	for pathIndex, path := range pathsA {

		// split this path
		split := strings.Split(path, "/")

		// iterate over each subpath
		for subPathIndex, subPath := range split {

			if pathIndex > 0 {

				//
				splitlast := strings.Split(pathsA[pathIndex-1], "/")

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
	return content
}
