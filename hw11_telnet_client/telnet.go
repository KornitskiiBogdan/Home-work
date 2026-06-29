package main

import (
	"errors"
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{address: address, timeout: timeout, in: in, out: out}
}

type telnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    net.Conn
}

func (c *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn
	return nil
}

func (c *telnetClient) Send() error {
	buf := make([]byte, 1024)
	n, err := c.in.Read(buf)
	if n > 0 {
		_, writeErr := c.conn.Write(buf[:n])
		if writeErr != nil {
			return writeErr
		}
	}
	if errors.Is(err, io.EOF) {
		return io.EOF
	}
	return err
}

func (c *telnetClient) Receive() error {
	buf := make([]byte, 1024)
	n, err := c.conn.Read(buf)
	if n > 0 {
		_, writeErr := c.out.Write(buf[:n])
		if writeErr != nil {
			return writeErr
		}
	}
	if errors.Is(err, io.EOF) {
		return io.EOF
	}
	return err
}

func (c *telnetClient) Close() error {
	return c.conn.Close()
}
