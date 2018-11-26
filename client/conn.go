package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "github.com/scotow/goxy/common"
)

func Dial(remoteAddr *net.TCPAddr) (*Conn, error) {
	httpAddr := fmt.Sprintf("http://%s/", remoteAddr.String())

	resp, err := http.Get(httpAddr)
	if err != nil {
		return nil, err
	}

	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	id, err := NewIdFromToken(string(token))
	if err != nil {
		return nil, err
	}

	conn := Conn{id, remoteAddr}

	return &conn, nil
}

type Conn struct {
	id         *Id
	remoteAddr *net.TCPAddr
}

func (c *Conn) Read(b []byte) (n int, err error) {
	httpAddr := fmt.Sprintf("http://%s/read/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", strings.NewReader(strconv.Itoa(len(b))))
	if err != nil {
		fmt.Println("HTTP read POST request", err.Error())
		return
	}
	defer resp.Body.Close()

	for {
		read, er := resp.Body.Read(b[n:])
		n += read
		err = er

		if err == io.EOF {
			err = nil
			break
		}

		if err != nil {
			break
		}
	}

	//fmt.Fprintf(c.logger, "Read: buffer size: %d. Read: %d.\n", len(b), n)

	return
}

func (c *Conn) Write(b []byte) (n int, err error) {
	httpAddr := fmt.Sprintf("http://%s/write/%s", c.remoteAddr.String(), c.id)

	resp, err := http.Post(httpAddr, "*/*", bytes.NewReader(b))
	if err != nil {
		n = 0
		return
	}
	defer resp.Body.Close()

	n = len(b)

	//fmt.Fprintf(c.logger, "Write: buffer size: %d. Written: %d.\n", len(b), n)
	return
}

func (c *Conn) Close() error {
	return nil
}

func (c *Conn) LocalAddr() net.Addr {
	panic("implement me")
}

func (c *Conn) RemoteAddr() net.Addr {
	panic("implement me")
}

func (c *Conn) SetDeadline(t time.Time) error {
	panic("implement me")
}

func (c *Conn) SetReadDeadline(t time.Time) error {
	panic("implement me")
}

func (c *Conn) SetWriteDeadline(t time.Time) error {
	panic("implement me")
}
