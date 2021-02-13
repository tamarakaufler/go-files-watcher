package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tamarakaufler/go-files-watcher/internal/daemon"
)

func main() {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	d := daemon.New(
		daemon.WithCommand("tree"),
		//daemon.WithCommand("go build -o go-files-watcher cmd/go-files-watcher/main.go"))
		daemon.WithExcluded([]string{"internal/daemon/fixtures/*"}),
		daemon.WithFrequency(5),
	)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.Frequency)*time.Second)
	defer cancel()
	d.Watch(ctx, sigCh)
}
