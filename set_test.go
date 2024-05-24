package main

import (
	"strings"
	"testing"
)

func TestHandleGetSetCommands(t *testing.T) {
	conn := getConnectionMock("*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "OK") {
		t.Errorf("expected response of ok got response of %s", conn.String())
	}

	conn1 := getConnectionMock("*2\r\n$3\r\nget\r\n$5\r\nhello\r\n")

	handleConnection(conn1)

	if !strings.Contains(conn1.Buffer.String(), "world") {
		t.Errorf("expected response of world got response of %s", conn1.Buffer.String())
	}
}

func TestHandleInvalidSetCommand(t *testing.T) {
	conn := getConnectionMock("*1\r\n$3\r\nset\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestHandleInvalidGetCommand(t *testing.T) {
	conn := getConnectionMock("*1\r\n$3\r\nget\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestHandleGetCommandOnMissingKey(t *testing.T) {
	conn := getConnectionMock("*2\r\n$3\r\nget\r\n$3\r\nmissing\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "_\r\n") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}
func TestExistsCommandOnMissingKey(t *testing.T) {
	conn := getConnectionMock("*2\r\n$6\r\nEXISTS\r\n$7\r\nmissing\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), ":0") {
		t.Errorf("expected response of command not found got response of %s", conn.String())
	}
}

func TestExisitsCommandOnExistingKey(t *testing.T) {
	conn := getConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nthekey\r\n$5\r\nfound\r\n")
	handleConnection(conn)

	conn1 := getConnectionMock("*2\r\n$6\r\nexists\r\n$6\r\nthekey\r\n")
	handleConnection(conn1)

	if !strings.Contains(conn1.String(), ":1") {
		t.Errorf("expected to find one, found %s", conn1.String())
	}

}

func TestStrLenCommand(t *testing.T) {
	conn := getConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nlength\r\n$4\r\nfour\r\n")
	handleConnection(conn)

	conn1 := getConnectionMock("*2\r\n$6\r\nstrlen\r\n$6\r\nlength\r\n")
	handleConnection(conn1)

	if !strings.Contains(conn1.String(), ":4") {
		t.Errorf("expected to find 4, found: %s", conn1.String())
	}

}
