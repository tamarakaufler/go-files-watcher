package main

import (
	"fmt"

	"github.com/tamarakaufler/go-files-watcher/internal/daemon"
)

func main() {
	d := daemon.New(daemon.WithCommand("go build -o go-files-watcher cmd/go-files-watcher"))
	fmt.Printf("%+v\n", d)

	d.Watch()
}
