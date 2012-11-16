package client

import (
	"net/textproto"
	"os"
	"strings"
	"time"
)

type KeyValueMap map[string]string

func list(conn *textproto.Conn) []string {
	id, err := conn.Cmd("list")
	values := make([]string, 0, 16)
	if err != nil {
		panic("error in list connection err is " + err.Error())
		return values
	}
	conn.StartResponse(id)
	for line, err := conn.ReadLine(); err == nil && len(line) > 0; line, err = conn.ReadLine() {
		if line[0] == '#' {
			continue
		}
		values = append(values, strings.Fields(line)...)
		break
	}
	conn.EndResponse(id)
	if err != nil {
		panic("error in list readline err is " + err.Error())
	}
	return values
}

func pair(input, sep string) (string, string) {
	both := strings.SplitN(input, sep, 2)
	return both[0], both[1]
}

func fetch(conn *textproto.Conn, what string) (values KeyValueMap) {
	id, err := conn.Cmd("fetch %s", what)
	values = make(KeyValueMap)
	if err != nil {
		panic("error in fetch connection for " + what + ", err is " + err.Error())
		return values
	}

	conn.StartResponse(id)
	dotlines, err := conn.ReadDotLines()
	conn.EndResponse(id)
	if err != nil {
		panic("error in fetch dotlines for" + what + ", err is " + err.Error())
		return values
	}
	for _, line := range dotlines {
		if line[0] == '.' {
			break
		}
		if line[0] == '#' {
			continue
		}
		key, value := pair(line, " ")
		key, _ = pair(key, ".")
		values[key] = value
	}
	return values
}

func NewMuninClient(hostport string, interval time.Duration, done <-chan os.Signal) <-chan KeyValueMap {
	conn, err := textproto.Dial("tcp4", hostport)
	if err != nil {
		panic("error connecting to munin " + err.Error())
		return nil
	}
	data := make(chan KeyValueMap, 1)
	go func() {
		ticker := time.NewTicker(interval)
		for {
			kv := make(KeyValueMap)
			select {
			case <-ticker.C:
				headers := list(conn)
				for _, prefix := range headers {
					for key, value := range fetch(conn, prefix) {
						kv[prefix+"."+key] = value
					}
				}
				data <- kv
			case <-done:
				ticker.Stop()
				close(data)
				conn.Close()
				return
			}
		}
	}()
	return data
}
