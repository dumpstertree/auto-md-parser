package main

import (
	"fmt"
	"os"
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
		content += "<br>"
	}

	content += page.Content

	// add footer if requested
	if useFooter {

		useTags := len(page.Tags) > 0
		useSource := page.Source != ""

		if useTags || useSource {
			content += "<div style='page-break-after: always;'></div>\n"
			content += "<div style='page-break-after: always;'></div>\n"
			content += "<hr/>\n"
			content += "<div style='page-break-after: always;'></div>\n"
			content += "<div style='page-break-after: always;'></div>\n"
		}

		// only add tags if exist
		if useTags {
			for _, tag := range page.Tags {
				content += "<a href='" + tag.LinkName + ".html'>" + tag.DisplayName + "</a>, "
			}
			content += "\n"
		}
		if useTags || useSource {
			content += "<div style='page-break-after: always;'></div>\n"
		}

		// create source link
		if useSource {
			content += "<div style='text-align: right'>\n"
			content += "<a href='" + page.Source + "'>SOURCE</a>\n"
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
