package munin

import (
	"io"
	"net/textproto"
	"os"
	"strings"
	"time"
)

type KeyValueMap map[string]string

type Client struct {
	conn    *textproto.Conn
	headers []string
	values  chan KeyValueMap
}

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
		panic("error in fetch dotlines for " + what + ", err is " + err.Error())
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

func (c *Client) Run(interval time.Duration, done <-chan os.Signal) <-chan KeyValueMap {
	go func() {
		ticker := time.NewTicker(interval)
		for {
			kv := make(KeyValueMap)
			select {
			case <-ticker.C:
				c.headers = list(c.conn)
				for _, prefix := range c.headers {
					for key, value := range fetch(c.conn, prefix) {
						kv[prefix+"."+key] = value
					}
				}
				c.values <- kv
			case <-done:
				ticker.Stop()
				close(c.values)
				c.conn.Close()
				return
			}
		}
	}()
	return c.values
}

func newMuninClient(connection io.ReadWriteCloser) (client *Client, err error) {
	conn := textproto.NewConn(connection)
	// skip the banner
	if _, err = conn.ReadLine(); err != nil {
		return nil, err
	}
	client = &Client{
		conn:   conn,
		values: make(chan KeyValueMap, 1),
	}
	return client, nil
}

func NewMuninClient(connection io.ReadWriteCloser, interval time.Duration, done <-chan os.Signal) <-chan KeyValueMap {
	client, err := newMuninClient(connection)
	if err != nil {
		panic(err)
	}
	return client.Run(interval, done)
}
