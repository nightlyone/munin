package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/nightlyone/munin/client"
)

func main() {
	done := make(chan os.Signal, 32)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	valChan := client.NewMuninClient("127.0.0.1:munin", 1 * time.Second, done)
	for values := range valChan {
		for key, value := range values {
			println(key, ":", value)
		}
	}
}
