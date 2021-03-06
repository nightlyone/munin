package main

import (
	"flag"
	"github.com/nightlyone/munin"
	"net"
	"os"
	"os/signal"
	"time"
)

var server = flag.String("server", "localhost", "host we query munin from")
var interval = flag.Duration("interval", 1*time.Minute, "interval between queries. Valid time units are ns, us (or µs), ms, s, m, h.")

func main() {
	flag.Parse()
	done := make(chan os.Signal, 32)
	signal.Notify(done, os.Interrupt, os.Kill)

	conn, err := net.Dial("tcp", net.JoinHostPort(*server, "munin"))
	if err != nil {
		println("error connecting to " + *server + ", error" + err.Error())
		os.Exit(2)
	}
	valChan := munin.NewMuninClient(conn, *interval, done)
	for values := range valChan {
		for key, value := range values {
			println(key, ":", value)
		}
	}
}
