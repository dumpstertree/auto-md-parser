package main

import (
	"fmt"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type DirectoryWatcher struct {
}

var w *fsnotify.Watcher
var paths []string

func (t *DirectoryWatcher) IsDirty() bool {
	newPaths := find(inputPath, ".layout")
	if !arraysEqual(paths, newPaths) || w == nil {
		paths = newPaths
		fmt.Println("new array")
		if w != nil {
			w.Close()
		}
		w = watch(find(inputPath, ".layout"))
		return true
	}

	if w != nil {
		fmt.Println("Watching directory : " + inputPath)
		//for {
		select {

		case e, ok := <-w.Events:

			// not ok
			if !ok {
				return false
			}

			// invalid type
			if !strings.Contains(e.Name, ".layout") {
				return false
			}

			fmt.Println("Change found in file : " + e.Name)
			//
			return true
		}
	}

	return false
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
