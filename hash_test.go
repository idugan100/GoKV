package main

import (
	"strings"
	"testing"
)

func TestHandleHSetandHGet(t *testing.T) {
	conn := getConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "3") {
		t.Errorf("expected response of 3 got response of %s", conn.String())
	}

	conn1 := getConnectionMock("*6\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nage\r\n$3\r\njob\r\n$3\r\nlol\r\n")

	handleConnection(conn1)

	if !strings.Contains(conn1.Buffer.String(), "isaac") || !strings.Contains(conn1.Buffer.String(), "20") || !strings.Contains(conn1.Buffer.String(), "swe") || !strings.Contains(conn1.Buffer.String(), "_") {
		t.Errorf("expected response of 1) \"isaac\" 2) \"20\" 3) \"swe\" 4) (nil) got response of %s", conn1.Buffer.String())
	}
}

func TestHandleHsetIncorrectNumberArgs(t *testing.T) {
	conn := getConnectionMock("*3\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected error about incorrect number of arguements got response of %s", conn.String())
	}
}

func TestHandleHgetIncorrectNumberArgs(t *testing.T) {
	conn := getConnectionMock("*2\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected error about incorrect number of arguements got response of %s", conn.String())
	}
}

func TestHandleHexists(t *testing.T) {
	conn := getConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	handleConnection(conn)

	conn1 := getConnectionMock("*3\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n$3\r\nage\r\n")
	handleConnection(conn1)

	if !strings.Contains(conn1.String(), ":1") {
		t.Errorf("expected result of 1, got %s", conn1.String())
	}
}

func TestHandleHexistsMissingArgs(t *testing.T) {
	conn := getConnectionMock("*2\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n")
	handleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected result of incorrect number of args error, got %s", conn.String())
	}
}

func TestHandleHexistsNotFound(t *testing.T) {
	conn := getConnectionMock("*3\r\n$7\r\nhexists\r\n$6\r\nkey\r\n$5\r\nvalue\r\n")
	handleConnection(conn)

	if !strings.Contains(conn.String(), ":0") {
		t.Errorf("expected result of 0, got %s", conn.String())
	}
}

func TestHandleDel(t *testing.T) {
	conn := getConnectionMock("*3\r\n$3\r\nSET\r\n$10\r\ndeletedkey\r\n$13\r\ndeletedvalue\r\n")
	handleConnection(conn)

	conn2 := getConnectionMock("*2\r\n$3\r\nDEL\r\n$10\r\ndeletedkey\r\n")
	handleConnection(conn2)

	if !strings.Contains(conn2.String(), ":1") {
		t.Errorf("expected response of :1, recieved %s", conn.String())
	}

}

func TestHandleHstrlen(t *testing.T) {
	conn := getConnectionMock("*4\r\n$4\r\nhset\r\n$6\r\nlength\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	handleConnection(conn)

	conn2 := getConnectionMock("*3\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n$3\r\nkey\r\n")
	handleConnection(conn2)

	if !strings.Contains(conn2.String(), ":3") {
		t.Errorf("expected :3 got value %s", conn2.String())
	}

}

func TestHandleHstrlenInvalidArgs(t *testing.T) {
	conn := getConnectionMock("*2\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n")
	handleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected incorrect number of args error got value %s", conn.String())
	}

}
