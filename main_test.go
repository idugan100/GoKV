package main

import (
	"bytes"
	"strings"
	"testing"
)

// this is a mock for the network connection that is read from and written to
type TestingReadWriter struct {
	*bytes.Buffer
	*strings.Reader
}

func (rw TestingReadWriter) Read(p []byte) (n int, err error) {
	return rw.Reader.Read(p)
}

func (rw TestingReadWriter) Write(p []byte) (n int, err error) {
	return rw.Buffer.Write(p)
}

func (rw TestingReadWriter) Close() error {
	return nil
}

func getTestingReadWriter(inputString string) TestingReadWriter {
	inputStringReader := strings.NewReader(inputString)
	var b bytes.Buffer

	s := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}
	return s
}
func TestHandleGetSetCommands(t *testing.T) {
	conn := getTestingReadWriter("*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "OK") {
		t.Errorf("expected response of ok got response of %s", conn.String())
	}

	conn1 := getTestingReadWriter("*2\r\n$3\r\nget\r\n$5\r\nhello\r\n")

	handleConnection(conn1)

	if !strings.Contains(conn1.Buffer.String(), "world") {
		t.Errorf("expected response of world got response of %s", conn1.Buffer.String())
	}
}

func TestHandlePingCommand(t *testing.T) {
	conn := getTestingReadWriter("*1\r\n$4\r\nPING\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "PONG") {
		t.Errorf("expected response of PONG got response of %s", conn.String())
	}
}

func TestHandleLolWutCommand(t *testing.T) {
	conn := getTestingReadWriter("*1\r\n$6\r\nLOLWUT\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), ":):):)") {
		t.Errorf("expected response of GoKV 0.1 :):):) got response of %s", conn.String())
	}
}

func TestHandleSetCommendError(t *testing.T) {
	conn := getTestingReadWriter("*1\r\n$3\r\nset\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestHandleInvalidCommand(t *testing.T) {
	conn := getTestingReadWriter("*1\r\n$3\r\n123\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "command not found") {
		t.Errorf("expected response of command not found got response of %s", conn.String())
	}
}

func TestHandleIndvalidCommand(t *testing.T) {
	conn := getTestingReadWriter("*1\r\n$3\r\n123\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "command not found") {
		t.Errorf("expected response of command not found got response of %s", conn.String())
	}
}
