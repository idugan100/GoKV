package main

import (
	"strings"
	"testing"
)

func TestHandlePingCommand(t *testing.T) {
	conn := getConnectionMock("*1\r\n$4\r\nPING\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "PONG") {
		t.Errorf("expected response of PONG got response of %s", conn.String())
	}
}

func TestHandleLolWutCommand(t *testing.T) {
	conn := getConnectionMock("*1\r\n$6\r\nLOLWUT\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), ":):):)") {
		t.Errorf("expected response of GoKV 0.1 :):):) got response of %s", conn.String())
	}
}
