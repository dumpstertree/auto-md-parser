package main

import "sort"

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

		sortedTags := taggedPages
		sort.Slice(sortedTags, func(i, j int) bool {
			return sortedTags[i].DisplayName < sortedTags[j].DisplayName
		})

		content := ""
		for i, page := range sortedTags {
			content += page.DisplayName
			if i != len(sortedTags)-1 {
				content += " \\"
			}
			content += "\n"
		}

		pages = append(pages, *makePageExplicit("Tags/Index/", tag.DisplayName, tag.LinkName, content, "", "", nil))
	}

	// generate page for all tags

	allTagsContent := ""
	x := 0

	keys := make([]PageTag, 0, len(allTags))
	for pt := range allTags {
		keys = append(keys, pt)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].DisplayName < keys[j].DisplayName
	})

	for i, key := range keys {

		//allTagsContent += "<a href='" + tag.LinkName + ".html'>" + tag.DisplayName + "</a>"
		allTagsContent += " " + key.DisplayName
		if i != len(allTags)-1 {
			allTagsContent += " \\"
		}
		allTagsContent += "\n"
		x++
	}
	pages = append(pages, *makePage("Tags/", "Index", allTagsContent, "", "", nil))

	// return
	return pages
}
