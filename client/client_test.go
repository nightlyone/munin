package client

import (
	"bufio"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// Thanks rsc (Russ Cox, Google Inc.) for this nice fake code
type fakePipes struct {
	io.ReadCloser
	io.WriteCloser
}

func (p *fakePipes) Close() error {
	p.ReadCloser.Close()
	p.WriteCloser.Close()
	return nil
}

func fakeConnect(t *testing.T) (io.ReadWriteCloser, error) {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	go fakeServer(t, &fakePipes{r1, w2})
	return &fakePipes{r2, w1}, nil
}

func fakeServer(t *testing.T, rw io.ReadWriteCloser) {
	b := bufio.NewReader(rw)
	_, err := rw.Write([]byte(fakeReply[""]))
	if err != nil {
		t.Errorf("fakeServer: Banner write causes %s", err)
	}
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				t.Errorf("fakeServer: Reading causes %s", err)
			}
			break
		}
		reply := fakeReply[strings.TrimSpace(line)]
		if reply == "" {
			break
		}
		_, err = rw.Write([]byte(reply))
		if err != nil {
			t.Errorf("fakeServer: Writing causes %s", err)
			break
		}
	}
	rw.Close()
}

// End of rsc (Russ Cox, Google Inc.) code

var fakeReply = map[string]string{
	"":     "# munin node at localhost\n",
	"list": "cpu\n",
	"quit": "",
	"fetch cpu": "user.value 234600\n" +
		"nice.value 1931\n" +
		"system.value 80354\n" +
		"idle.value 11153645\n" +
		"iowait.value 98142\n" +
		"irq.value 1\n" +
		"softirq.value 706\n" +
		"steal.value 0\n" +
		".\n",
}

var expectedReply = map[string]string{
	"cpu.user":    "234600",
	"cpu.nice":    "1931",
	"cpu.system":  "80354",
	"cpu.idle":    "11153645",
	"cpu.iowait":  "98142",
	"cpu.irq":     "1",
	"cpu.softirq": "706",
	"cpu.steal":   "0",
}

func TestConnect(t *testing.T) {
	t.Logf("Seting up fake server\n")
	conn, _ := fakeConnect(t)
	t.Logf("Seting up fake server done\n")
	interval := time.Millisecond * 200
	done := make(chan os.Signal, 32)
	go func(die chan<- os.Signal) {
		time.Sleep(interval * 2)
		die <- os.Interrupt
	}(done)
	valChan := NewMuninClient(conn, interval, done)
	reply := make(map[string]string)
	for values := range valChan {
		for key, value := range values {
			reply[key] = value
		}
	}

	for key, value := range expectedReply {
		if _, exists := reply[key]; !exists {
			t.Errorf("missing key %s in reply", key)
		}
		if reply[key] != value {
			t.Errorf("bad value for key %s in reply: got = %s, want = %s", key, reply[key], value)
		}
	}
	for key, value := range reply {
		if _, exists := expectedReply[key]; !exists {
			t.Errorf("extra key %s in reply, value %s", key, value)
		}
	}
}
