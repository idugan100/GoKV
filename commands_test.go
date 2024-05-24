package main

import (
	"strings"
	"testing"
)

func TestHandlePingCommand(t *testing.T) {
	conn := GetConnectionMock("*1\r\n$4\r\nPING\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "PONG") {
		t.Errorf("expected response of PONG got response of %s", conn.String())
	}
}

func TestHandleLolWutCommand(t *testing.T) {
	conn := GetConnectionMock("*1\r\n$6\r\nLOLWUT\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), ":):):)") {
		t.Errorf("expected response of GoKV 0.1 :):):) got response of %s", conn.String())
	}
}

func TestHandleHSetandHGet(t *testing.T) {
	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "3") {
		t.Errorf("expected response of 3 got response of %s", conn.String())
	}

	conn1 := GetConnectionMock("*6\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nage\r\n$3\r\njob\r\n$3\r\nlol\r\n")

	HandleConnection(conn1)

	if !strings.Contains(conn1.Buffer.String(), "isaac") || !strings.Contains(conn1.Buffer.String(), "20") || !strings.Contains(conn1.Buffer.String(), "swe") || !strings.Contains(conn1.Buffer.String(), "_") {
		t.Errorf("expected response of 1) \"isaac\" 2) \"20\" 3) \"swe\" 4) (nil) got response of %s", conn1.Buffer.String())
	}
}

func TestHandleHsetIncorrectNumberArgs(t *testing.T) {
	conn := GetConnectionMock("*3\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected error about incorrect number of arguements got response of %s", conn.String())
	}
}

func TestHandleHgetIncorrectNumberArgs(t *testing.T) {
	conn := GetConnectionMock("*2\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected error about incorrect number of arguements got response of %s", conn.String())
	}
}

func TestHandleHexists(t *testing.T) {
	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	HandleConnection(conn)

	conn1 := GetConnectionMock("*3\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n$3\r\nage\r\n")
	HandleConnection(conn1)

	if !strings.Contains(conn1.String(), ":1") {
		t.Errorf("expected result of 1, got %s", conn1.String())
	}
}

func TestHandleHexistsMissingArgs(t *testing.T) {
	conn := GetConnectionMock("*2\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n")
	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected result of incorrect number of args error, got %s", conn.String())
	}
}

func TestHandleHexistsNotFound(t *testing.T) {
	conn := GetConnectionMock("*3\r\n$7\r\nhexists\r\n$6\r\nkey\r\n$5\r\nvalue\r\n")
	HandleConnection(conn)

	if !strings.Contains(conn.String(), ":0") {
		t.Errorf("expected result of 0, got %s", conn.String())
	}
}

func TestHandleDel(t *testing.T) {
	conn := GetConnectionMock("*3\r\n$3\r\nSET\r\n$10\r\ndeletedkey\r\n$13\r\ndeletedvalue\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$3\r\nDEL\r\n$10\r\ndeletedkey\r\n")
	HandleConnection(conn2)

	if !strings.Contains(conn2.String(), ":1") {
		t.Errorf("expected response of :1, recieved %s", conn.String())
	}

}

func TestHandleHstrlen(t *testing.T) {
	conn := GetConnectionMock("*4\r\n$4\r\nhset\r\n$6\r\nlength\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*3\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n$3\r\nkey\r\n")
	HandleConnection(conn2)

	if !strings.Contains(conn2.String(), ":3") {
		t.Errorf("expected :3 got value %s", conn2.String())
	}

}

func TestHandleHstrlenInvalidArgs(t *testing.T) {
	conn := GetConnectionMock("*2\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n")
	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect number") {
		t.Errorf("expected incorrect number of args error got value %s", conn.String())
	}

}

func TestHandleGetSetCommands(t *testing.T) {
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
	conn := GetConnectionMock("*1\r\n$3\r\nset\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestHandleInvalidGetCommand(t *testing.T) {
	conn := GetConnectionMock("*1\r\n$3\r\nget\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "incorrect") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}

func TestHandleGetCommandOnMissingKey(t *testing.T) {
	conn := GetConnectionMock("*2\r\n$3\r\nget\r\n$3\r\nmissing\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), "_\r\n") {
		t.Errorf("expected response of incorrect number of args got response of %s", conn.String())
	}
}
func TestExistsCommandOnMissingKey(t *testing.T) {
	conn := GetConnectionMock("*2\r\n$6\r\nEXISTS\r\n$7\r\nmissing\r\n")

	HandleConnection(conn)

	if !strings.Contains(conn.String(), ":0") {
		t.Errorf("expected response of command not found got response of %s", conn.String())
	}
}

func TestExisitsCommandOnExistingKey(t *testing.T) {
	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nthekey\r\n$5\r\nfound\r\n")
	HandleConnection(conn)

	conn1 := GetConnectionMock("*2\r\n$6\r\nexists\r\n$6\r\nthekey\r\n")
	HandleConnection(conn1)

	if !strings.Contains(conn1.String(), ":1") {
		t.Errorf("expected to find one, found %s", conn1.String())
	}

}

func TestStrLenCommand(t *testing.T) {
	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nlength\r\n$4\r\nfour\r\n")
	HandleConnection(conn)

	conn1 := GetConnectionMock("*2\r\n$6\r\nstrlen\r\n$6\r\nlength\r\n")
	HandleConnection(conn1)

	if !strings.Contains(conn1.String(), ":4") {
		t.Errorf("expected to find 4, found: %s", conn1.String())
	}

}