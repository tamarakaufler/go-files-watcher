package daemon

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// Watch watches for changes in files at regular intervals
func (d *Daemon) Watch(ctx context.Context, sigCh chan os.Signal) {
	cmdParts := strings.Split(d.Command, " ")
	doneCh := make(chan struct{})
	tick := time.NewTicker(d.frequency)

	go func() {
		for {
			select {
			case <-sigCh:
				fmt.Println("signal sent")
				os.Exit(0)
			case <-doneCh:
				d.mux.Lock()

				cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
				// cmd.Stdout = os.Stdout
				// cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					fmt.Printf("ERROR: %s\n", errors.Wrap(err, "error occurred processing during file watch"))
					d.mux.Unlock()
					continue
				}
				fmt.Println("command completed successfully")
				d.mux.Unlock()
			}
		}
	}()

	for {
		ctxR := context.Context(ctx)
		select {
		case <-tick.C:
			fmt.Println("time to repeat")
			files := d.collectFiles(ctxR)
			d.processFiles(ctxR, files, doneCh)
		}
	}
}

// collectFiles checks if any watched file has changed
func (d *Daemon) collectFiles(ctx context.Context) []os.FileInfo {
	var files []os.FileInfo

	filepath.Walk(d.BasePath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() ||
			strings.HasPrefix(path, ".git") ||
			(!info.IsDir() && filepath.Ext(path) != d.Extention) {
			return err
		}

		files = append(files, info)
		fmt.Printf("FILE info:  %s\n", info.Name())

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
