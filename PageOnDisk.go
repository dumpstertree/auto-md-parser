package main

import (
	"fmt"
	"os"
	"sort"
)

type PageOnDisk struct {
	Path    string
	Content string
}

func WriteToDisk(path string, page Page, useFooter bool) *PageOnDisk {
	//
	content := ""

	if useFooter {
		content += "# " + "<p style='font-size: 15px;'>" + page.Path + "</p>" + "\n"
		content += "# " + "<p style='font-size: 40px;'>" + page.DisplayName + "</p>" + "\n"
		content += "\n"
	}

	content += page.Content

	// add footer if requested
	if useFooter {

		source := page.Source
		if page.OverrideSource != "" {
			source = page.OverrideSource
		}

		useTags := len(page.Tags) > 0
		useSource := source != ""

		if useTags || useSource {
			content += "<div style='page-break-after: always;'></div>\n"
			content += "<div style='page-break-after: always;'></div>\n"
			content += "\n"
			content += "<hr/>\n"
			content += "\n"
			content += "<div style='page-break-after: always;'></div>\n"
			content += "<div style='page-break-after: always;'></div>\n"
			content += "\n"
		}

		// only add tags if exist
		if useTags {
			sortedTags := page.Tags
			sort.Slice(sortedTags, func(i, j int) bool {
				return sortedTags[i].DisplayName < sortedTags[j].DisplayName
			})

			for i, tag := range sortedTags {
				content += "<a href='" + tag.LinkName + ".html'>" + tag.DisplayName + "</a>"
				if i < len(page.Tags)-1 {
					content += ", "
				}
			}
			content += "\n"
			content += "<div style='page-break-after: always;'></div>\n"
			content += "\n"
		}

		// create source link
		if useSource {

			content += "<div style='text-align: right'>\n"
			content += "<a href='" + source + "'>SOURCE</a>\n"
			content += "</div>\n"
		}
	}

	err := os.WriteFile(path+page.LinkName+".md", []byte(content), 0644)
	if err != nil {
		fmt.Println("Failed to write file at : " + path)
		return nil
	}
	return &PageOnDisk{Path: path, Content: content}
}
func (t *PageOnDisk) Delete() {

}
