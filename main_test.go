package main

import (
	"bytes"
	"strings"
	"testing"
)

type TestingReadWriter struct {
	*bytes.Buffer
	*strings.Reader
}

func (rw TestingReadWriter) Read(p []byte) (n int, err error) {
	return rw.Reader.Read(p)
}

// Write writes data to the underlying bytes.Buffer
func (rw TestingReadWriter) Write(p []byte) (n int, err error) {
	return rw.Buffer.Write(p)
}
func TestHandleGetSetCommands(t *testing.T) {
	inputstring := "*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n"
	inputStringReader := strings.NewReader(inputstring)
	var b bytes.Buffer

	s := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}

	handleConnection(s)

	if !strings.Contains(b.String(), "OK") {
		t.Errorf("expected response of ok got response of %s", b.String())
	}

	b.Reset()
	inputstring = "*2\r\n$3\r\nget\r\n$5\r\nhello\r\n"
	inputStringReader = strings.NewReader(inputstring)

	s1 := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}

	handleConnection(s1)

	if !strings.Contains(b.String(), "world") {
		t.Errorf("expected response of world got response of %s", b.String())
	}
}

func TestHandlePingCommand(t *testing.T) {
	inputstring := "*1\r\n$4\r\nPING\r\n"
	inputStringReader := strings.NewReader(inputstring)
	var b bytes.Buffer

	s := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}

	handleConnection(s)

	if !strings.Contains(b.String(), "PONG") {
		t.Errorf("expected response of PONG got response of %s", b.String())
	}
}

func TestHandleLolWutCommand(t *testing.T) {
	inputstring := "*1\r\n$6\r\nLOLWUT\r\n"
	inputStringReader := strings.NewReader(inputstring)
	var b bytes.Buffer

	s := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}

	handleConnection(s)

	if !strings.Contains(b.String(), ":):):)") {
		t.Errorf("expected response of GoKV 0.1 :):):) got response of %s", b.String())
	}
}

func TestHandleSetCommendError(t *testing.T) {
	inputstring := "*1\r\n$3\r\nset\r\n"
	inputStringReader := strings.NewReader(inputstring)
	var b bytes.Buffer

	s := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}

	handleConnection(s)

	if !strings.Contains(b.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", b.String())
	}
}

func TestHandleInvalidCommand(t *testing.T) {
	inputstring := "*1\r\n$3\r\n123\r\n"
	inputStringReader := strings.NewReader(inputstring)
	var b bytes.Buffer

	s := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}

	handleConnection(s)

	if !strings.Contains(b.String(), "command not found") {
		t.Errorf("expected response of command not found got response of %s", b.String())
	}
}

func TestHandleIndvalidCommand(t *testing.T) {
	inputstring := "*1\r\n$3\r\n123\r\n"
	inputStringReader := strings.NewReader(inputstring)
	var b bytes.Buffer

	s := TestingReadWriter{
		Buffer: &b,
		Reader: inputStringReader,
	}

	handleConnection(s)

	if !strings.Contains(b.String(), "command not found") {
		t.Errorf("expected response of command not found got response of %s", b.String())
	}
}
