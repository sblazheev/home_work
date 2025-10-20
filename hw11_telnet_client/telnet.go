package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ErrEOFSend = errors.New("EOF")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type MyTelnetClient struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	connect net.Conn
	UserEOF bool
}

type ProxyReader struct {
	in    io.Reader
	isEOF bool
}

func (r *ProxyReader) Read(p []byte) (int, error) {
	rb, er := r.in.Read(p)
	if er != nil && er == io.EOF {
		r.isEOF = true
	}
	return rb, er
}

func (c *MyTelnetClient) Connect() error {
	connect, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.connect = connect
	return nil
}

func (c *MyTelnetClient) GetAddress() string {
	return c.address
}

func (c *MyTelnetClient) Close() (err error) {
	return c.connect.Close()
}

func (c *MyTelnetClient) Receive() error {
	_, err := io.Copy(c.out, c.connect)
	return err
}

func (c *MyTelnetClient) Send() error {
	proxyReader := ProxyReader{
		in: c.in,
	}
	_, err := io.Copy(c.connect, &proxyReader)
	if proxyReader.isEOF {
		c.UserEOF = true
	}
	return err
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &MyTelnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func RunTelnetClient(tc MyTelnetClient) error {
	if err := tc.Connect(); err != nil {
		return fmt.Errorf("error connecting: %w", err)
	}
	defer tc.Close()

	os.Stderr.Write([]byte("...Connected to " + tc.GetAddress() + "\n"))

	osSignal := make(chan os.Signal, 1)
	signal.Notify(osSignal, syscall.SIGINT)
	errorConnect := make(chan error, 1)

	go func() {
		if err := tc.Send(); err != nil {
			errorConnect <- err
			return
		}
		if tc.UserEOF {
			errorConnect <- ErrEOFSend
		}
	}()

	go func() {
		if err := tc.Receive(); err != nil {
			errorConnect <- err
			return
		}
	}()

	for {
		select {
		case <-osSignal:
			tc.Close()
			return nil
		case er := <-errorConnect:
			if errors.Is(er, ErrEOFSend) {
				os.Stderr.Write([]byte("...EOF\n"))
			} else {
				os.Stderr.Write([]byte("...Connection was closed by peer\n"))
			}
			return nil
		}
	}
}
