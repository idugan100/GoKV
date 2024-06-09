package main

import (
	"bytes"
	"net"
	"strings"
	"testing"
)

// this is a mock for the network connection that is read from and written to
type ConnectionMock struct {
	*bytes.Buffer
	*strings.Reader
}

func (c ConnectionMock) Read(p []byte) (n int, err error) {
	return c.Reader.Read(p)
}

func (c ConnectionMock) Write(p []byte) (n int, err error) {
	return c.Buffer.Write(p)
}

func (c ConnectionMock) Close() error {
	return nil
}

func GetConnectionMock(inputString string) ConnectionMock {
	inputStringReader := strings.NewReader(inputString)
	var b bytes.Buffer

	c := ConnectionMock{
		Buffer: &b,
		Reader: inputStringReader,
	}
	return c
}
func TestStartServer(t *testing.T) {
	go startServer()
	conn, err := net.Dial("tcp", "localhost:6379")

	if err != nil {
		t.Errorf("error when connecting to GoKV server %s", err.Error())
	}
	defer conn.Close()

	message := "*1\r\n$4\r\nPING\r\n"
	_, err = conn.Write([]byte(message))
	if err != nil {
		t.Errorf("error when sending message over tcp connection %s", err.Error())
	}

	response := make([]byte, 1024)
	_, err = conn.Read(response)

	if err != nil {
		t.Errorf("error when sending message over tcp connection %s", err.Error())
	}

	expectedResult := "PONG"
	if !strings.Contains(string(response), expectedResult) {
		t.Errorf("expected result of '%s', recieved '%s'", expectedResult, string(response))
	}

}

func TestHandleInvalidCommand(t *testing.T) {
	conn := GetConnectionMock("*1\r\n$3\r\n123\r\n")

	HandleConnection(conn)

	expectedResult := "command not found"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result of '%s' got response of '%s'", expectedResult, conn.String())
	}
}
