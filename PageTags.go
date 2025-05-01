package main

type PageTags struct {
	Content  []string
	Content2 string
	Pages    []Page
}

func makePageTags(allPages []Page) *PageTags {
	x := new(PageTags)
	x.Pages = buildPageTags(allPages)
	return x
}

func buildPageTags(allPages []Page) []Page {

	pages := []Page{}

	// match tags to the pages they access
	allTags := make(map[PageTag][]Page)
	for _, p := range allPages {
		for _, t := range p.Tags {
			allTags[t] = append(allTags[t], p)
		}
	}

	// generate pages for each tag
	for tag, taggedPages := range allTags {

		content := ""
		for i, page := range taggedPages {
			content += "<a href='" + page.Name + ".html'>" + page.Name + "</a>\n"
			if i != len(taggedPages)-1 {
				content += "\\"
			}
		}

		pages = append(pages, *makePage("Tags/All/", tag.LinkName, content, "", nil))
	}

	// generate page for all tags
	allTagsContent := ""
	x := 0
	for tag, _ := range allTags {

		allTagsContent += "<a href='" + tag.LinkName + ".html'>" + tag.DisplayName + "</a>\n"
		if x != len(allTags)-1 {
			allTagsContent += "\\"
		}
		x++
	}
	pages = append(pages, *makePage("Tags/", "All", allTagsContent, "", nil))

	// return
	return pages
}
