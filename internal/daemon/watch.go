package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type dirContent struct {
	files []string
	dirs  []string
}

// Watch watches for changes in files at regular intervals
func (d Daemon) Watch() {
	for {
		fmt.Println("I am watching you ...")
		d.processFiles()
		time.Sleep(time.Duration(5 * time.Second))
	}
}

// processFiles
func (d Daemon) processFiles() {
	var files []string

	filepath.Walk(d.BasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() ||
			//info.Mode() == "" ||
			strings.HasPrefix(path, ".git") ||
			!info.IsDir() && filepath.Ext(path) != d.Extention {
			return err
		}
		files = append(files, strings.Join([]string{path, info.Name() + ".go"}, string(filepath.Separator)))

		fmt.Printf("path: %s, name: %+v, mode: %+v err: %+v\n\n", path, info.Name(), info.Mode(), err)
		return nil
	})

	fmt.Printf("ALL files: %s\n\n", strings.Join(files, "\n"))

	//return err
}
