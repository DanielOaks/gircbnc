// written by Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package ircbnc

import (
	"bufio"
	"io"
	"net"
	"strings"
)

// Socket represents an IRC socket.
type Socket struct {
	Closed bool
	conn   net.Conn
	reader *bufio.Reader
	buffer string
}

// NewSocket returns a new Socket.
func NewSocket(conn net.Conn) Socket {
	return Socket{
		conn:   conn,
		reader: bufio.NewReader(conn),
	}
}

// Close stops a Socket from being able to send/receive any more data.
func (socket *Socket) Close() {
	if socket.Closed {
		return
	}
	socket.Closed = true
	socket.conn.Close()
}

// Read returns a single IRC line from a Socket.
func (socket *Socket) Read() (string, error) {
	if socket.Closed {
		return "", io.EOF
	}

	lineBytes, err := socket.reader.ReadBytes('\n')

	// convert bytes to string
	line := string(lineBytes[:])

	// read last message properly (such as ERROR/QUIT/etc), just fail next reads/writes
	if err == io.EOF {
		socket.Close()
	}

	if err == io.EOF && strings.TrimSpace(line) != "" {
		// don't do anything
	} else if err != nil {
		return "", err
	}

	return line, nil
}

// Write sends the given line out of Socket.
func (socket *Socket) Write(line string) error {
	if socket.Closed {
		return io.EOF
	}

	// write data
	_, err := socket.conn.Write([]byte(line))
	if err != nil {
		socket.Close()
		return err
	}
	return nil
}
