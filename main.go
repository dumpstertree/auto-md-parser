package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/xuri/excelize/v2"
)

var outputPath = "/media/dumpstertree/67FF-9C76/Development/Docker/mdBook/content/src/"
var inputPath = "./layout/"

var images = []string{}

var generateTags = true
var generateContent = true
var generateSummary = true

// main
func main() {

	fmt.Println("Starting Parse")
	// variables declaration
	var inFlag string
	var outFlag string

	// flags declaration using flag package
	flag.StringVar(&inFlag, "i", "", "")
	flag.StringVar(&outFlag, "o", "", "")
	flag.Parse()
	if inFlag == "" {
		fmt.Println("No Input Path Provided")
		return
	}
	if outFlag == "" {
		fmt.Println("No Output Path Provided")
		return
	}

	outputPath = outFlag
	inputPath = inFlag

	// create a watcher that watches for file changes
	watcher := new(DirectoryWatcher)

	// loop waiting for changes
	for {

		if watcher.IsDirty() {
			Reload()
		}

		time.Sleep(1 * time.Second)
	}
}
func Reload() {

	// clear any old data
	clearCachedImages()

	// data loader
	d := new(DataLoader)
	d.Load(inputPath)

	//
	p := new(PageParser)
	p.Load(d)

	// get all pages
	allPages := []Page{}

	if generateContent {
		// load page content
		allPages = append(allPages, makePageContent(d, p).Pages...)
	}

	// load page tags
	if generateTags {
		allPages = append(allPages, makePageTags(allPages).Pages...)
	}

	// iterate over each page so far
	for i, p := range allPages {
		allPages[i].Content = p.applyExternalLinks(allPages)
	}
	for i, p := range allPages {
		allPages[i].Content = p.applyInternalLinks(allPages)
	}

	// load page summary
	if generateSummary {
		allPages = append(allPages, makePageSummary(allPages).Pages...)
	}

	// write all
	for _, i := range allPages {
		WriteToDisk(outputPath, i, i.DisplayName != "SUMMARY")
	}

	// cleanup
	clearUnusedMD(allPages)

	// unload
	d.Clear()
}
func arraysEqual(arr1, arr2 []string) bool {
	if len(arr1) != len(arr2) {
		return false
	}
	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return false
		}
	}
	return true
}
func clearCachedImages() {

	for _, i := range find(outputPath, ".png") {
		os.Remove(i)
	}
	for _, i := range find(outputPath, ".jpg") {
		os.Remove(i)
	}
}
func clearUnusedMD(pages []Page) {
	for _, i := range find(outputPath, ".md") {
		split := strings.Split(i, "/")
		path := strings.ReplaceAll(split[len(split)-1], ".md", "")
		contained := false
		for _, p := range pages {

			if strings.EqualFold(path, p.LinkName) {
				contained = true
				break
			}
		}
		if !contained {
			fmt.Println("Removing Unused Markdown at : " + i + " with name :" + path)
			os.Remove(i)
		}
	}
}

// array functions
func remove(s []string, t string) []string {
	for i := 0; i < len(s); i++ {
		if s[i] == t {
			s[i] = s[len(s)-1]
		}
	}
	return s[:len(s)-1]
}

// path functions
func find(root, ext string) []string {
	var a []string
	filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	return a
}

// parse
func parseCompoundCollumnString(input string, sheet string, row int, file *excelize.File) string {
	var lineEnd = ""
	strArr := []rune(input)
	last := ' '
	for x, y := range strArr {
		if last == '@' {

			// do nothing

		} else if y == '@' {

			nxt := string(strArr[x+1])
			z, err := file.GetPictures(sheet, nxt+strconv.Itoa(row))
			if err == nil {
				for _, pic := range z {

					u := uuid.New()
					name := u.String() + pic.Extension
					if err := os.WriteFile(outputPath+name, pic.File, 0644); err != nil {
						fmt.Println(err)

					}
					lineEnd += "<img src=" + name + ">"
					images = append(images, name)

				}
			}

			val, err := file.GetCellValue(sheet, nxt+strconv.Itoa(row))
			if err != nil {
				panic(err)
			}

			lineEnd += val
		} else {
			lineEnd += string(y)
		}

		last = y
	}

	return lineEnd
}
func (p Page) applyExternalLinks(pages []Page) string {
	content := p.Content
	for _, w := range strings.Split(content, " ") {

		_, err := url.ParseRequestURI(w)
		if err != nil {
			continue
		}

		u, err := url.Parse(w)
		if err != nil || u.Scheme == "" || u.Host == "" {
			continue
		}

		content = strings.ReplaceAll(content, w, "<a href="+w+">"+w+"</a>")

	}

	return content
}
func (p Page) applyInternalLinks(pages []Page) string {

	content := p.Content
	for _, c := range pages {
		if c.DisplayName == p.DisplayName {
			continue
		}

		re := regexp.MustCompile(`([^\n]|[^ ])` + c.DisplayName + `([^\n]|[^ ]|[^.]|[^?]|[^!][^,][^;])`)
		content = re.ReplaceAllString(content, "<a href='"+c.LinkName+".html'>"+c.DisplayName+"</a>")

	}
	return content
}

// constructors
func makePageExplicit(path string, name string, linkName string, content string, source string, tags []string) *Page {

	t := []PageTag{}
	for _, i := range tags {
		t = append(t, *makeTag(i))
	}

	name = strings.Replace(name, " ", "_", -1)

	return &Page{
		DisplayName: name,
		LinkName:    linkName,
		Path:        path,
		Content:     content,
		Source:      source,
		Tags:        t,
	}
}
func makePage(path string, name string, content string, source string, tags []string) *Page {

	t := []PageTag{}
	for _, i := range tags {
		t = append(t, *makeTag(i))
	}

	linkName := strings.Replace(name, " ", "_", -1)
	linkName = strings.ToLower(linkName)

	return &Page{
		DisplayName: name,
		LinkName:    linkName,
		Path:        path,
		Content:     content,
		Source:      source,
		Tags:        t,
	}
}
func makeTag(name string) *PageTag {

	return &PageTag{
		LinkName:    "tag-" + strings.ToLower(strings.ReplaceAll(name, " ", "_")),
		DisplayName: "#" + strings.ToLower(strings.ReplaceAll(name, " ", "_")),
	}
}

// data
type ISubsection interface {
	Write(page Page, allPages []Page, sheet string, file *excelize.File) []Page
}
type OrderedLayout struct {
	Title            string
	URL              string
	Path             string
	IncludeSheets    []string
	ExcludeSheets    []string
	Tags             []string
	LayoutSubsection []ISubsection     `json:"-"`
	Sections         []json.RawMessage `json:"sections"`
}
type BaseSubsection struct {
	Type string `json:"type"`
}
type TextSubsection struct {
	Color  string `json:"color"`
	Bold   bool   `json:"bold"`
	Italic bool   `json:"italic"`
}

// text
func (s TextSubsection) ModifyTextStart(text string) string {

	if s.Bold {
		text += "<b>"
	}
	if s.Italic {
		text += "<i>"
	}
	if s.Color != "" {
		text += "<font color=" + s.Color + ">"

	}
	return text
}
func (s TextSubsection) ModifyTextEnds(text string) string {

	if s.Bold {
		text += "</b>"
	}
	if s.Italic {
		text += "</i>"
	}
	if s.Color != "" {
		text += "</font>"

	}
	return text
}

type Page struct {
	DisplayName string
	LinkName    string
	Path        string
	Content     string
	Source      string
	Tags        []PageTag
}

type PageTag struct {
	LinkName    string
	DisplayName string
}
