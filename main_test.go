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

func getConnectionMock(inputString string) ConnectionMock {
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

	if !strings.Contains(string(response), "PONG") {
		t.Errorf("expected response of PONG, recieved %s", string(response))
	}

}

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

func TestHandleInvalidCommand(t *testing.T) {
	conn := getConnectionMock("*1\r\n$3\r\n123\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "command not found") {
		t.Errorf("expected response of command not found got response of %s", conn.String())
	}
}

func TestHandleIndvalidCommand(t *testing.T) {
	conn := getConnectionMock("*1\r\n$3\r\n123\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "command not found") {
		t.Errorf("expected response of command not found got response of %s", conn.String())
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

func TestUnknownCommand(t *testing.T) {
	conn := getConnectionMock("*1\r\n$3\r\nxyz\r\n")
	handleConnection(conn)
	if !strings.Contains(conn.String(), "command not found") {
		t.Errorf("expected command not found. Got: %s", conn.String())
	}
}

func TestHandleHSetandHGet(t *testing.T) {
	conn := getConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	handleConnection(conn)

	if !strings.Contains(conn.String(), "3") {
		t.Errorf("expected response of 3 got response of %s", conn.String())
	}
	// hget myinfo name age job lol
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
