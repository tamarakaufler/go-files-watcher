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
		daemon.WithFrequency(5),
		//daemon.WithCommand("go build -o go-files-watcher cmd/go-files-watcher/main.go"))
		daemon.WithCommand("tree"),
		daemon.WithExcluded([]string{"internal/daemon/fixtures/basepath/subdir1/*", "test.go"}))

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(d.Frequency)*time.Second)
	defer cancel()
	d.Watch(ctx, sigCh)
}
