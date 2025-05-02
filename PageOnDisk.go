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
	content := page.Content

	// add footer if requested
	if useFooter {

		content += "<div style='page-break-after: always;'></div>\n"
		content += "---\n"
		content += "<div style='page-break-after: always;'></div>\n"

		// only add tags if exist
		if len(page.Tags) > 0 {
			for _, tag := range page.Tags {
				content += "<a href='" + tag.LinkName + ".html'>" + tag.DisplayName + "</a>, "
			}
			content += "\n"
		}

		// create source link
		if page.Source != "" {
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
