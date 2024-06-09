package main

import (
	"strings"
	"testing"

	"github.com/idugan100/GoKV/handlers"
)

func TestHandlePingCommand(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*1\r\n$4\r\nPING\r\n")

	HandleConnection(conn)

	expectedResult := "PONG"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result of '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestHandleLolWutCommand(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*1\r\n$6\r\nLOLWUT\r\n")

	HandleConnection(conn)

	expectedResult := "GoKV 0.1 :):):)"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result of '%s' got response of '%s'", expectedResult, conn.String())
	}
}

// HSET and HGET tests
func TestHandleHSetandHGet(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	HandleConnection(conn)

	expectedResult := ":3"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}

	conn1 := GetConnectionMock("*6\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nage\r\n$3\r\njob\r\n$3\r\nlol\r\n")

	HandleConnection(conn1)

	expectedName := "isaac"
	expectedAge := "20"
	expectedJob := "swe"

	if !strings.Contains(conn1.Buffer.String(), expectedName) || !strings.Contains(conn1.Buffer.String(), expectedAge) || !strings.Contains(conn1.Buffer.String(), expectedJob) || !strings.Contains(conn1.Buffer.String(), "_") {
		t.Errorf("expected response of 1) '%s' 2) '%s' 3) '%s' 4) (nil) got response of %s", expectedName, expectedAge, expectedJob, conn1.Buffer.String())
	}
}

func TestHandleHsetIncorrectNumberArgs(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n")

	HandleConnection(conn)

	expectedResult := "incorrect number"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestHandleHgetIncorrectNumberArgs(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*2\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n")

	HandleConnection(conn)

	expectedResult := "incorrect number"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

// HEXISTS tests
func TestHandleHexists(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	HandleConnection(conn)

	conn1 := GetConnectionMock("*3\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n$3\r\nage\r\n")
	HandleConnection(conn1)

	expectedResult := ":1"
	if !strings.Contains(conn1.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn1.String())
	}
}

func TestHandleHexistsMissingArgs(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*2\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestHandleHexistsNotFound(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$7\r\nhexists\r\n$6\r\nkey\r\n$5\r\nvalue\r\n")
	HandleConnection(conn)

	expectedResult := ":0"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s', got '%s'", expectedResult, conn.String())
	}
}

// HSTRLEN tests
func TestHandleHstrlen(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*4\r\n$4\r\nhset\r\n$6\r\nlength\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*3\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n$3\r\nkey\r\n")
	HandleConnection(conn2)

	expectedResult := ":3"
	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected results '%s' got value '%s'", expectedResult, conn2.String())
	}

}

func TestHandleHstrlenInvalidArgs(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*2\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

// HLEN tests
func TestHlenIncorrectArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*1\r\n$4\r\nHLEN\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestHlenHashNotFound(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$4\r\nHLEN\r\n$6\r\nmissing\r\n")
	HandleConnection(conn)

	expectedResult := ":0"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestHlen(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$4\r\ndata\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$4\r\nHLEN\r\n$4\r\ndata\r\n")
	HandleConnection(conn2)

	expectedResult := ":3"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

// HGETALL tests
func TestHGetAllIncorrectArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*1\r\n$7\r\nhgetall\r\n")

	HandleConnection(conn)

	expectedResult := "incorrect number of args"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected '%s', got '%s'", expectedResult, conn.String())
	}
}

func TestHGetAllKeyNotFound(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$7\r\nhgetall\r\n$7\r\nmissing\r\n")

	HandleConnection(conn)

	expectedResult := "*0"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected '%s', got '%s'", expectedResult, conn.String())
	}
}

func TestHGetAll(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$4\r\ndata\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")

	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$7\r\nhgetall\r\n$4\r\ndata\r\n")

	HandleConnection(conn2)

	expectedResults := [...]string{"isaac", "20", "swe"}

	for _, expectedResult := range expectedResults {
		if !strings.Contains(conn2.String(), expectedResult) {
			t.Errorf("expected '%s', got '%s'", expectedResult, conn2.String())
		}
	}

}

// HSETNX tests
func TestHSetNXIncorrectArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$6\r\nHSETNX\r\n$3\r\nkey\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number of args"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to get '%s' error, found: %s", expectedResult, conn.String())
	}
}

func TestHSetNXAlreadyExists(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*4\r\n$6\r\nHSETNX\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nbob\r\n")
	HandleConnection(conn2)

	expectedResult := ":0"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected '%s', got '%s'", expectedResult, conn2.String())
	}

}

func TestHSetNX(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*4\r\n$6\r\nHSETNX\r\n$6\r\nmyinfo\r\n$8\r\nfavcolor\r\n$4\r\nblue\r\n")
	HandleConnection(conn2)

	expectedResult := ":1"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected '%s', got '%s'", expectedResult, conn2.String())
	}

}

// HDEL tests
func TestHDelIncorrectNumberArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$4\r\nHDEL\r\ndata\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number of args"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected '%s' recieved '%s'", expectedResult, conn.String())
	}
}

func TestHDel(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*5\r\n$4\r\nHDEL\r\n$6\r\nmyinfo\r\n$8\r\nfavcolor\r\n$4\r\nname\r\n$3\r\nage\r\n")
	HandleConnection(conn2)

	expectedResult := ":2"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected '%s', got '%s'", expectedResult, conn2.String())
	}
}
