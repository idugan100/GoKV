package main

import (
	"strings"
	"testing"

	"github.com/idugan100/GoKV/handlers"
)

func TestHandleDel(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$3\r\nSET\r\n$10\r\ndeletedkey\r\n$13\r\ndeletedvalue\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$3\r\nDEL\r\n$10\r\ndeletedkey\r\n")
	HandleConnection(conn2)

	if !strings.Contains(conn2.String(), ":1") {
		t.Errorf("expected response of :1, recieved %s", conn.String())
	}

}

func TestHandleGetSetCommands(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "OK") {
		t.Errorf("expected response of ok got response of %s", conn.String())
	}

	conn1 := GetConnectionMock("*2\r\n$3\r\nget\r\n$5\r\nhello\r\n")

	HandleConnection(conn1)

	if !strings.Contains(conn1.Buffer.String(), "world") {
		t.Errorf("expected response of world got response of %s", conn1.Buffer.String())
	}
}

func TestHandleInvalidSetCommand(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*1\r\n$3\r\nset\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestHandleInvalidGetCommand(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*1\r\n$3\r\nget\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestHandleGetCommandOnMissingKey(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*2\r\n$3\r\nget\r\n$3\r\nmissing\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "_\r\n") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestExistsCommandOnMissingKey(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*2\r\n$6\r\nEXISTS\r\n$7\r\nmissing\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), ":0") {
		t.Errorf("expected response of command not found got response of %s", conn.String())
	}
}

func TestExisitsCommandOnExistingKey(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nthekey\r\n$5\r\nfound\r\n")
	HandleConnection(conn)

	conn1 := GetConnectionMock("*2\r\n$6\r\nexists\r\n$6\r\nthekey\r\n")
	HandleConnection(conn1)

	if !strings.Contains(conn1.String(), ":1") {
		t.Errorf("expected to find one, found %s", conn1.String())
	}

}

func TestStrLenCommand(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nlength\r\n$4\r\nfour\r\n")
	HandleConnection(conn)

	conn1 := GetConnectionMock("*2\r\n$6\r\nstrlen\r\n$6\r\nlength\r\n")
	HandleConnection(conn1)

	if !strings.Contains(conn1.String(), ":4") {
		t.Errorf("expected to find 4, found: %s", conn1.String())
	}

}

func TestSetNXInvalidNumberOfArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number of arguements"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find %s, found: %s", expectedResult, conn.String())
	}

}

func TestSetNXValueExisits(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	HandleConnection(conn)

	conn1 := GetConnectionMock("*3\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n$6\r\nnewval\r\n")
	HandleConnection(conn1)

	expectedResult := ":0"
	if !strings.Contains(conn1.String(), expectedResult) {
		t.Errorf("expected to find %s, found: %s", expectedResult, conn1.String())
	}
}

func TestSetNXValueNotSet(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n$6\r\nnewval\r\n")
	HandleConnection(conn)

	expectedResult := ":1"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find %s, found: %s", expectedResult, conn.String())
	}
}

func TestIncrInvalidNumberOfArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*1\r\n$4\r\nINCR\r\b")
	HandleConnection(conn)

	expectedResult := "incorrect number of arguements"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestIncrOnExistingKey(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$1\r\n1\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n")
	HandleConnection(conn2)

	expectedResult := ":2"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn2.String())
	}
}

func TestIncrOnNewKey(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n")
	HandleConnection(conn)

	expectedResult := ":0"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestIncrOnInvalidDataType(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$3\r\none\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n")
	HandleConnection(conn2)

	expectedResult := "incorrect data type"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn2.String())
	}
}
