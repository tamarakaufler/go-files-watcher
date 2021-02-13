package daemon

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

// Watch watches for changes in files at regular intervals
func (d *Daemon) Watch(ctx context.Context, sigCh chan os.Signal) {
	fmt.Print("\nStarting the watcher daemon ⌚ 👀 ... \n\n")
	cmdParts := strings.Split(d.Command, " ")

	// use when a change is detected to avoid processing further files
	doneCh := make(chan struct{})
	// use when a change is detected, after successfully running the command,
	// to cancel already created goroutines
	cancelCh := make(chan struct{})
	tick := time.NewTicker(d.frequency)

	go func() {
		for {
			select {
			case <-sigCh:
				fmt.Println("You interrupted me 👹!")
				os.Exit(0)
			case <-doneCh:
				d.mux.Lock()

				cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
				// these can be commented out if not needed
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Printf("ERROR: %s\n", errors.Wrap(err, "error occurred processing during file watch"))
					cancelCh <- struct{}{}
					d.mux.Unlock()
					continue
				}
				fmt.Print("command completed successfully\n\n")
				d.mux.Unlock()
			}
		}
	}()

	for {
		ctxR, cancel := context.WithCancel(ctx)
		select {
		case <-tick.C:
			//fmt.Println("time to repeat")

			// implementation 1
			// d.walkThroughFiles(ctxR, doneCh)

			// implementation 2
			// files := d.CollectFiles(ctxR)
			// d.processFiles(ctxR, files, doneCh)

			// implementation 3
			files := d.CollectFiles(ctxR)
			d.processFilesInParallel(ctxR, files, doneCh)
		case <-cancelCh:
			cancel()
		}
	}
}

// CollectFiles checks if any watched file has changed.
// The Walk function continues the walk while theere is no error and stops
// when the filepath.WalkFunc exits with error.
func (d *Daemon) walkThroughFiles(ctx context.Context, doneCh chan struct{}) {
	filepath.Walk(d.BasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() ||
			strings.HasPrefix(path, ".git") ||
			(!info.IsDir() && filepath.Ext(path) != d.Extention) {
			return nil
		}

		fmt.Printf("FILE info:  %s\n", info.Name())

		lastChecked := time.Now().Add(-d.frequency)
		if info.ModTime().After(lastChecked) {
			fmt.Printf("\tFile %s has changed\n", info.Name())
			// trigger running of the command
			doneCh <- struct{}{}
			// return any known error to stop walking through the dir content
			return io.EOF
		}
		return nil
	})
}

// CollectFiles checks if any watched file has changed
func (d *Daemon) CollectFiles(ctx context.Context) []os.FileInfo {
	var files []os.FileInfo

	filepath.Walk(d.BasePath, func(path string, info os.FileInfo, err error) error {
		//fmt.Println(os.Getwd())
		//fmt.Printf("\n\npath: %s,\ninfo: %+v, error: %+v\n\n", path, info, err)

		if info.IsDir() ||
			strings.HasPrefix(path, ".git") ||
			(!info.IsDir() && filepath.Ext(path) != d.Extention) {
			return err // this will be nil if there is no problem with the file
		}

		if len(d.Excluded) != 0 {
			isExcl, err := d.IsExcluded(ctx, path, info)
			if err != nil {
				panic(errors.Wrap(err, "cannot proccess exclusion of files"))
			}
			if isExcl {
				return nil
			}
		}

		files = append(files, info)
		//fmt.Printf("FILE info:  %s - %s\n", path, info.Name())

		return nil
	})

	return files
}

func (d *Daemon) processFiles(ctx context.Context, files []os.FileInfo, doneCh chan struct{}) {

	fmt.Println("GOT to processing ...")

	for _, f := range files {
		time.Sleep(100 * time.Millisecond)

		lastChecked := time.Now().Add(-d.frequency)
		if f.ModTime().After(lastChecked) {
			fmt.Printf("File %s has changed\n", f.Name())
			doneCh <- struct{}{}
			break
		}
	}
}

func (d *Daemon) processFilesInParallel(ctx context.Context, files []os.FileInfo, doneCh chan struct{}) {
	wg := &sync.WaitGroup{}

	stopCh := make(chan struct{})
	continueCh := make(chan struct{})

	// Files are checked in parallel. When a change is found, a message is sent
	// to the doneCh channel to interrupt the looping through the rest of the files. When no chenge is found, a message is sent to the continueCh channel to continue looping.
	// Note: I triend to use select default to continue the looping but that
	// did not work.
	for _, f := range files {
		wg.Add(1)
		go func(wg *sync.WaitGroup, f os.FileInfo, doneCh chan struct{}, stopCh chan struct{}) {
			defer wg.Done()
			time.Sleep(100 * time.Millisecond)

			lastChecked := time.Now().Add(-d.frequency)
			if f.ModTime().After(lastChecked) {
				fmt.Printf("File %s has changed\n", f.Name())
				stopCh <- struct{}{}
				return
			}
			continueCh <- struct{}{}
		}(wg, f, doneCh, stopCh)

		select {
		case <-stopCh:
			doneCh <- struct{}{}
			break
		case <-continueCh:
		}
	}

	wg.Wait()
}

// IsExcluded filters files based on custom exclusion configuration
func (d *Daemon) IsExcluded(ctx context.Context, path string, info os.FileInfo) (bool, error) {
	toExclude := false

	for _, ex := range d.Excluded {
		// deal with regex
		if strings.Contains(ex, "*") {
			r, err := regexp.Compile(ex)
			if err != nil {
				return false, errors.Wrap(err, "cannot exclude files")
			}
			if r.Match([]byte(path)) {
				return true, nil
			}
			// deal with exact matches
		} else if info.Name() == ex || path == ex {
			return true, nil
		}

	}
	return toExclude, nil
}
