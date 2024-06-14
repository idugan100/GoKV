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

	expectedResult := ":1"
	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn2.String())
	}

}

func TestHandleGetSetCommands(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n")

	HandleConnection(conn)

	expectedResult := "OK"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}

	conn1 := GetConnectionMock("*2\r\n$3\r\nget\r\n$5\r\nhello\r\n")

	HandleConnection(conn1)
	expectedResult = "world"
	if !strings.Contains(conn1.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn1.String())
	}
}

func TestHandleInvalidSetCommand(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*1\r\n$3\r\nset\r\n")

	HandleConnection(conn)

	expectedResult := "incorrect number"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestHandleInvalidGetCommand(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*1\r\n$3\r\nget\r\n")

	HandleConnection(conn)

	expectedResult := "incorrect number"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestHandleGetCommandOnMissingKey(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*2\r\n$3\r\nget\r\n$3\r\nmissing\r\n")

	HandleConnection(conn)

	expectedResult := "_\r\n"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestExistsCommandOnMissingKey(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*2\r\n$6\r\nEXISTS\r\n$7\r\nmissing\r\n")

	HandleConnection(conn)

	expectedResult := ":0"
	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn.String())
	}
}

func TestExisitsCommandOnExistingKey(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nthekey\r\n$5\r\nfound\r\n")
	HandleConnection(conn)

	conn1 := GetConnectionMock("*2\r\n$6\r\nexists\r\n$6\r\nthekey\r\n")
	HandleConnection(conn1)

	expectedResult := ":1"
	if !strings.Contains(conn1.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn1.String())
	}
}

func TestStrLenCommand(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$6\r\nlength\r\n$4\r\nfour\r\n")
	HandleConnection(conn)

	conn1 := GetConnectionMock("*2\r\n$6\r\nstrlen\r\n$6\r\nlength\r\n")
	HandleConnection(conn1)

	expectedResult := ":4"
	if !strings.Contains(conn1.String(), expectedResult) {
		t.Errorf("expected result '%s' got response of '%s'", expectedResult, conn1.String())
	}
}

func TestSetNXInvalidNumberOfArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number of arguments"
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

	expectedResult := "incorrect number of arguments"

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

	expectedResult := ":1"

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

func TestDecrOnExistingKey(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$1\r\n1\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n")
	HandleConnection(conn2)

	expectedResult := ":0"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn2.String())
	}
}

func TestDecrOnNewKey(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n")
	HandleConnection(conn)

	expectedResult := ":-1"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestDecrOnInvalidDataType(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$3\r\none\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n")
	HandleConnection(conn2)

	expectedResult := "incorrect data type"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn2.String())
	}
}

func TestDecrInvalidNumberOfArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*1\r\n$4\r\nDECR\r\b")
	HandleConnection(conn)

	expectedResult := "incorrect number of arguments"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestRenameInvalidNumberOfArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number of arguments"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestRenameKeyNotFound(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n$6\r\nnewkey\r\n")
	HandleConnection(conn)

	expectedResult := "key to be renamed not found"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestRenameKeyFound(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*3\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n$6\r\nnewkey\r\n")
	HandleConnection(conn2)

	expectedResult := "OK"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn2.String())
	}

	conn3 := GetConnectionMock("*2\r\n$3\r\nGET\r\n$6\r\nnewkey\r\n")
	HandleConnection(conn3)
	expectedResult2 := "val"

	if !strings.Contains(conn3.String(), expectedResult2) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult2, conn3.String())
	}
}

func TestDecrbyInvalidNumberOfArgs(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*2\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect number of arguments"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestDecrbyInvalidType(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$6\r\nDECRBY\r\n$3\r\nkey\r\n$3\r\nval\r\n")
	HandleConnection(conn)

	expectedResult := "incorrect data type"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestDecrbyKeyNotFound(t *testing.T) {
	defer handlers.ClearData()

	conn := GetConnectionMock("*3\r\n$6\r\nDECRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n")
	HandleConnection(conn)

	expectedResult := ":-2"

	if !strings.Contains(conn.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn.String())
	}
}

func TestDecrby(t *testing.T) {
	defer handlers.ClearData()
	conn := GetConnectionMock("*3\r\n$3\r\nSET\r\n$3\r\nnum\r\n$1\r\n5\r\n")
	HandleConnection(conn)

	conn2 := GetConnectionMock("*3\r\n$6\r\nDECRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n")
	HandleConnection(conn2)

	expectedResult := ":3"

	if !strings.Contains(conn2.String(), expectedResult) {
		t.Errorf("expected to find '%s', found '%s'", expectedResult, conn2.String())
	}
}
