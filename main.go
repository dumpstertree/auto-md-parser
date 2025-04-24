package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"

	"github.com/google/uuid"

	"github.com/cavaliergopher/grab/v3"
	"github.com/xuri/excelize/v2"
	"golang.org/x/text/collate"
	"golang.org/x/text/language"
)

var outputPath = "/media/dumpstertree/67FF-9C76/Development/Docker/mdBook/content/src/"
var inputPath = "./layout/"

const tempPath = "./.temp/"
const summaryName = "SUMMARY"

var images = []string{}

var isDirty = false

// main
func Reload() {

	isDirty = false

	// generate
	paths := find(inputPath, ".layout")
	layouts := loadLayouts(paths)
	spreadsheets := loadSpreadheet(layouts)

	// get the page list
	allPages := buildPageLists(spreadsheets)

	// build the pages and get a summaryPages
	summaryPages := buildPageContent(spreadsheets, allPages)

	// tags
	summaryTags := buildPageTags(allPages, spreadsheets)

	// build the summary page
	buildSummaryPage(summaryPages, summaryTags)

	// cleanup
	cleanupUnlinked(allPages)

	// clear any data from the web
	cleanupTempData()

	fmt.Println("Parse Complete")

}
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

	for {

		newPaths := find(inputPath, ".layout")
		if !arraysEqual(paths, newPaths) || w == nil {
			paths = newPaths
			fmt.Println("new array")
			if w != nil {
				w.Close()
			}
			w = watch(find(inputPath, ".layout"))
			Reload()
		}

		if w != nil {
			watchLoop(w)
		}

		time.Sleep(1 * time.Second)
	}

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

var paths []string

var w *fsnotify.Watcher

// load
func loadLayouts(filePaths []string) []*OrderedLayout {

	layouts := []*OrderedLayout{}
	for _, x := range filePaths {

		// create layout object
		var layout *OrderedLayout

		// guard - read text
		content, err := os.ReadFile(x)
		if err != nil {
			fmt.Println("Faled to read file at : " + x)
			continue
		}

		// guard - convert to layout object
		err = json.Unmarshal(content, &layout)
		if err != nil {
			fmt.Println("Faled to parse JSON at : " + x)
			continue
		}

		layouts = append(layouts, layout)

		fmt.Println("Found Layout at: " + x)
	}

	return layouts
}
func loadSpreadheet(layouts []*OrderedLayout) map[*OrderedLayout]*excelize.File {

	// make temp dir if doesnt exist
	os.Mkdir(tempPath, 0777)

	excel := make(map[*OrderedLayout]*excelize.File)
	for _, x := range layouts {

		// guard - grab file
		resp, err := grab.Get(tempPath, x.URL)
		if err != nil {
			fmt.Println("Failed to get response for : " + x.Title)
			continue
		}

		// guard - couldnt open file
		file, err := excelize.OpenFile(resp.Filename)
		if err != nil {
			fmt.Println("Failed to open file for : " + x.Title)
			continue
		}

		excel[x] = file

		fmt.Println("Found Spreadsheet at: " + x.URL)
	}
	return excel
}

// cleanup
func cleanupUnlinked(allPages map[*excelize.File][]string) {
	for _, p1 := range find(outputPath, ".png") {
		// remove the ext
		p1Split := strings.Split(p1, "/")
		file := p1Split[len(p1Split)-1]

		// setup a flag
		match := false

		for _, p3 := range images {

			// if they are equal mark flag as match
			isMatch := file == p3
			if isMatch {
				match = true
				break
			}
		}

		// no match was found for file so remove
		if !match {
			os.Remove(p1)
		}
	}
	for _, p1 := range find(outputPath, ".jpg") {
		// remove the ext
		p1Split := strings.Split(p1, "/")
		file := p1Split[len(p1Split)-1]

		// setup a flag
		match := false

		for _, p3 := range images {

			// if they are equal mark flag as match
			isMatch := file == p3
			if isMatch {
				match = true
				break
			}
		}

		// no match was found for file so remove
		if !match {
			os.Remove(p1)
		}
	}
	// remove all pages currently there
	for _, p1 := range find(outputPath, ".md") {

		// remove the ext
		p1Split := strings.Split(p1, "/")
		file := p1Split[len(p1Split)-1]
		file = strings.ReplaceAll(file, ".md", "")

		// if this is the summary skip
		if file == summaryName {
			continue
		}
		if file == "Tags" {
			continue
		}
		if strings.ContainsRune(file, '_') {
			continue
		}

		if strings.Contains(file, "tag") {
			continue
		}

		// setup a flag
		match := false
		for _, p2 := range allPages {

			for _, p3 := range p2 {

				// if they are equal mark flag as match
				isMatch := file == p3
				if isMatch {
					match = true
					break
				}
			}
		}
		// no match was found for file so remove
		if !match {
			os.Remove(p1)
		}
	}
}
func cleanupTempData() {
	os.RemoveAll(tempPath)
}

// build
func buildPageTags(fileToPath map[*excelize.File][]string, layoutToFile map[*OrderedLayout]*excelize.File) []string {
	allTags := make(map[string][]string)

	for z, y := range layoutToFile {
		for _, x := range z.Tags {

			for _, path := range fileToPath[y] {
				allTags["_"+x] = append(allTags["_"+x], path)
			}
		}
	}

	for x, y := range allTags {
		content := ""
		for _, path := range y {
			content += "<a href='" + path + ".html'>" + path + "</a>"
			content += "\\"
			content += "\n"
		}
		// create a file
		err := os.WriteFile(outputPath+x+".md", []byte(content), 0644)
		fmt.Println("add: " + outputPath + x + ".md")

		if err != nil {
			panic(err)
		}
	}

	content := ""
	for y, _ := range allTags {

		//for _, path := range x {
		content += "<a href='" + y + ".html'>" + y + "</a>"
		content += "\\"
		content += "\n"

		//}

	}
	err := os.WriteFile(outputPath+"Tags"+".md", []byte(content), 0644)
	if err != nil {
		panic(err)
	}

	//c := collate.New(language.English, collate.IgnoreCase)
	//c.SortStrings(allTags)

	keys := make([]string, 0, len(allTags))
	for k := range allTags {
		keys = append(keys, k)
	}

	return keys
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
func buildPageContent(layoutToFile map[*OrderedLayout]*excelize.File, allPages map[*excelize.File][]string) map[string][]string {

	// flatten the pages for use when writing
	allPagesFlat := []string{}
	for _, y := range allPages {
		allPagesFlat = append(allPagesFlat, y...)
	}

	// initialize return value
	summary := make(map[string][]string)

	// add pages for each path
	for layout, file := range layoutToFile {

		// get sheets for this file
		reqSheets := allPages[file]

		//iterate on the required pages
		for _, sheet := range reqSheets {

			// page content
			var content = ""

			// add mandatory header
			content += "## " + sheet
			content += "\n"

			// add all subsections
			for _, subsection := range layout.LayoutSubsection {
				content = subsection.Write(content, sheet, allPagesFlat, file)
			}

			content += "<div style='page-break-after: always;'></div>\n"
			content += "---\n"
			content += "<div style='page-break-after: always;'></div>\n"

			for _, tag := range layout.Tags {
				content += "<a href='" + tag + ".html'>" + tag + "</a>, "
			}
			content += "\n"

			// create source link

			content += "<div style='text-align: right'>\n"
			content += "<a href='" + layout.URL + "'>SOURCE</a>\n"
			content += "</div>\n"

			// guard- output to .md
			err := os.WriteFile(outputPath+sheet+".md", []byte(content), 0644)
			if err != nil {
				panic(err)
			}

			summary[layout.Path] = append(summary[layout.Path], sheet)
		}

	}

	// return value
	return summary
}
func buildSummaryPage(summary map[string][]string, summaryTags []string) {

	// get all paths
	keys := make([]string, 0, len(summary))
	for k := range summary {
		keys = append(keys, k)
	}
	// sort
	c := collate.New(language.English, collate.IgnoreCase)
	c.SortStrings(keys)

	out := "#\n"
	for m, y := range keys {
		split := strings.Split(y, "/")
		for z, x := range split {

			if m > 0 {
				splitlast := strings.Split(keys[m-1], "/")

				if len(splitlast) > z && splitlast[z] == x {
					continue
				}
			}

			if z == 0 {
				// if this is the first level make it a header
				out += "# " + x
				out += "\n"

			} else {
				// if this is more than first level add a draft
				for f := 1; f < z; f++ {
					out += "	"
				}

				out += "- [" + x + "]()"
				out += "\n"
			}
		}
		for _, g := range summary[y] {
			for f := 1; f < len(split); f++ {
				out += "	"
			}

			out += "- [" + g + "](" + strings.ReplaceAll(g, " ", "") + ".md)"
			out += "\n"
		}
	}

	out += "# Tags\n"
	out += "- [Tags](Tags.md)\n"
	for _, x := range summaryTags {
		out += "	- [" + x + "](" + strings.ReplaceAll(x, " ", "") + ".md)"
		out += "\n"
	}

	err := os.WriteFile(outputPath+"SUMMARY.md", []byte(out), 0644)
	if err != nil {
		panic(err)
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
func parseCompoundCollumnString(input string, sheet string, row int, allPages []string, file *excelize.File) string {
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

	for _, w := range strings.Split(lineEnd, " ") {
		for _, x := range allPages {
			has := strings.EqualFold(w, x)
			if has {
				fmt.Println("match")
			}
			//fmt.Println("check page " + w + " : " + x)
			if has {
				lineEnd = strings.ReplaceAll(lineEnd, w, "<a href='"+x+".html'>"+strings.ReplaceAll(w, " ", "")+"</a>")

			}
		}
	}
	for _, w := range strings.Split(lineEnd, " ") {
		_, err := url.ParseRequestURI(w)
		if err != nil {
			continue
		}

		u, err := url.Parse(w)
		if err != nil || u.Scheme == "" || u.Host == "" {
			continue
		}

		lineEnd = strings.ReplaceAll(lineEnd, w, "<a href="+w+">"+w+"</a>")

	}

	return lineEnd
}

// data
type ISubsection interface {
	Write(input string, sheet string, allPages []string, file *excelize.File) string
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

// This is the most basic example: it prints events to the terminal as we
// receive them.
func watch(paths []string) *fsnotify.Watcher {
	if len(paths) < 1 {
		fmt.Println("must specify at least one path to watch")
		return nil
	}

	// Create a new watcher.
	w, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("creating a new watcher: %s", err)
		return nil
	}
	//defer w.Close()

	// Start listening for events.
	//go watchLoop(w)

	// Add all paths from the commandline.
	for _, p := range paths {
		err = w.Add(p)
		if err != nil {
			fmt.Println("%q: %s", p, err)
		}
	}

	//fmt.Println("ready; press ^C to exit")
	//<-make(chan struct{}) // Block forever

	return w
}

func watchLoop(w *fsnotify.Watcher) {
	fmt.Println("watch loop ")
	i := 0
	//for {
	select {
	// Read from Errors.
	case err, ok := <-w.Errors:
		if !ok { // Channel was closed (i.e. Watcher.Close() was called).
			return
		}
		fmt.Println("ERROR: %s", err)
	// Read from Events.
	case e, ok := <-w.Events:
		if !ok { // Channel was closed (i.e. Watcher.Close() was called).
			return
		}

		// Just print the event nicely aligned, and keep track how many
		// events we've seen.
		i++
		fmt.Println("%3d %s", i, e)
		isDirty = true
		fmt.Println("set dirty")

		//
		Reload()
	}
	//}
}
