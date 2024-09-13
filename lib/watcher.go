package lib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
	gitignore "github.com/sabhiram/go-gitignore"
)

func watchForChanges(dir string, reloadChan chan bool) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Printf("Error creating file watcher: %v\n", err)
		return
	}
	defer watcher.Close()

	// Load .gitignore if it exists
	ignorer, err := gitignore.CompileIgnoreFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		// If .gitignore doesn't exist or can't be read, create an empty ignorer
		ignorer = gitignore.CompileIgnoreLines([]string{}...)
	}

	// Function to check if a path should be ignored
	shouldIgnore := func(path string) bool {
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			return false
		}
		return ignorer.MatchesPath(relPath)
	}

	// Function to add a directory and its subdirectories to the watcher
	addDir := func(path string) error {
		return filepath.Walk(path, func(walkPath string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() && !shouldIgnore(walkPath) {
				return watcher.Add(walkPath)
			}
			return nil
		})
	}

	// Add the initial directory
	if err := addDir(dir); err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if !shouldIgnore(event.Name) {
				switch {
				case event.Op&fsnotify.Write == fsnotify.Write,
					event.Op&fsnotify.Create == fsnotify.Create,
					event.Op&fsnotify.Remove == fsnotify.Remove,
					event.Op&fsnotify.Rename == fsnotify.Rename:

					reloadChan <- true

					// If a new directory is created, add it to the watcher
					if event.Op&fsnotify.Create == fsnotify.Create {
						if info, err := os.Stat(event.Name); err == nil && info.IsDir() {
							if err := addDir(event.Name); err != nil {
								fmt.Printf("Error adding new directory to watcher: %v\n", err)
							}
						}
					}
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Error watching files: %v\n", err)
		}
	}
}
