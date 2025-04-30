package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cavaliergopher/grab/v3"
	"github.com/xuri/excelize/v2"
)

type DataLoader struct {
	LoadedLayouts      []*OrderedLayout
	LoadedSpreadsheets map[*OrderedLayout]*excelize.File
}

func (d *DataLoader) Load(inputPath string) {

	// generate
	paths := find(inputPath, FILE_EXT)

	// load data
	d.LoadedLayouts = loadLayouts(paths)
	d.LoadedSpreadsheets = loadSpreadheet(d.LoadedLayouts)
}

func (d DataLoader) Clear() {

	os.RemoveAll(tempPath)
}

const tempPath = "./temp/"
const FILE_EXT = ".layout"

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
