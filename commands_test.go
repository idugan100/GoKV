package main

import (
	"strings"
	"testing"

	"github.com/idugan100/GoKV/handlers"
)

type TableTest struct {
	Commands       []string
	ExpectedOutput string
}

var CommandTableTests = []TableTest{
	{[]string{"*1\r\n$4\r\nPING\r\n"}, "PONG"},
	{[]string{"*1\r\n$6\r\nLOLWUT\r\n"}, "GoKV 0.1 :):):)"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n"}, ":3"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*6\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nage\r\n$3\r\njob\r\n$3\r\nlol\r\n"}, "isaac"},
	{[]string{"*3\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n"}, handlers.InvalidArgsNumberError{Command: "HSET"}.Error()},
	{[]string{"*2\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n"}, handlers.InvalidArgsNumberError{Command: "HGET"}.Error()},
	{[]string{"*2\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n"}, handlers.InvalidArgsNumberError{Command: "HEXISTS"}.Error()},
	{[]string{"*3\r\n$7\r\nhexists\r\n$6\r\nkey\r\n$5\r\nvalue\r\n"}, ":0"},
	{[]string{"*4\r\n$4\r\nhset\r\n$6\r\nlength\r\n$3\r\nkey\r\n$3\r\nval\r\n", "*3\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n$3\r\nkey\r\n"}, ":3"},
	{[]string{"*2\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n"}, handlers.InvalidArgsNumberError{Command: "HSTRLEN"}.Error()},
	{[]string{"*1\r\n$4\r\nHLEN\r\n"}, handlers.InvalidArgsNumberError{Command: "HLEN"}.Error()},
	{[]string{"*2\r\n$4\r\nHLEN\r\n$6\r\nmissing\r\n"}, ":0"},
	{[]string{"*8\r\n$4\r\nhset\r\n$4\r\ndata\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*2\r\n$4\r\nHLEN\r\n$4\r\ndata\r\n"}, ":3"},
	{[]string{"*1\r\n$7\r\nhgetall\r\n"}, handlers.InvalidArgsNumberError{Command: "HGETALL"}.Error()},
	{[]string{"*2\r\n$7\r\nhgetall\r\n$7\r\nmissing\r\n"}, "*0"},
	{[]string{"*8\r\n$4\r\nhset\r\n$4\r\ndata\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*2\r\n$7\r\nhgetall\r\n$4\r\ndata\r\n"}, "swe"},
	{[]string{"*2\r\n$6\r\nHSETNX\r\n$3\r\nkey\r\n"}, handlers.InvalidArgsNumberError{Command: "HSETNX"}.Error()},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*4\r\n$6\r\nHSETNX\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nbob\r\n"}, ":0"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*4\r\n$6\r\nHSETNX\r\n$6\r\nmyinfo\r\n$8\r\nfavcolor\r\n$4\r\nblue\r\n"}, ":1"},
	{[]string{"*2\r\n$4\r\nHDEL\r\ndata\r\n"}, handlers.InvalidArgsNumberError{Command: "HDEL"}.Error()},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*5\r\n$4\r\nHDEL\r\n$6\r\nmyinfo\r\n$8\r\nfavcolor\r\n$4\r\nname\r\n$3\r\nage\r\n"}, ":2"},
	{[]string{"*1\r\n$3\r\n123\r\n"}, "command not found"},
	{[]string{"*1\r\n$3\r\nset\r\n"}, handlers.InvalidArgsNumberError{Command: "SET"}.Error()},
	{[]string{"*1\r\n$3\r\nget\r\n"}, handlers.InvalidArgsNumberError{Command: "GET"}.Error()},
	{[]string{"*2\r\n$3\r\nget\r\n$3\r\nmissing\r\n"}, "_\r\n"},
	{[]string{"*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n"}, "OK"},
	{[]string{"*3\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n$6\r\nnewval\r\n"}, ":1"},
	{[]string{"*1\r\n$4\r\nINCR\r\b"}, handlers.InvalidArgsNumberError{Command: "INCR"}.Error()},
	{[]string{"*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n"}, ":1"},
	{[]string{"*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n"}, ":-1"},
	{[]string{"*1\r\n$4\r\nDECR\r\b"}, handlers.InvalidArgsNumberError{Command: "DECR"}.Error()},
	{[]string{"*1\r\n$6\r\nDECRBY\r\b"}, handlers.InvalidArgsNumberError{Command: "DECRBY"}.Error()},
	{[]string{"*3\r\n$6\r\nDECRBY\r\n$3\r\nkey\r\n$3\r\nval\r\n"}, handlers.InvalidDataTypeError{Command: "DECRBY"}.Error()},
	{[]string{"*3\r\n$6\r\nDECRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n"}, ":-2"},
	{[]string{"*3\r\n$3\r\nSET\r\n$3\r\nnum\r\n$1\r\n5\r\n", "*3\r\n$6\r\nDECRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n"}, ":3"},
	{[]string{"*1\r\n$6\r\nINCRBY\r\b"}, handlers.InvalidArgsNumberError{Command: "INCRBY"}.Error()},
	{[]string{"*3\r\n$6\r\nINCRBY\r\n$3\r\nkey\r\n$3\r\nval\r\n"}, handlers.InvalidDataTypeError{Command: "INCRBY"}.Error()},
	{[]string{"*3\r\n$6\r\nINCRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n"}, ":2"},
	{[]string{"*3\r\n$3\r\nSET\r\n$3\r\nnum\r\n$1\r\n5\r\n", "*3\r\n$6\r\nINCRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n"}, ":7"},
	{[]string{"*2\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n"}, handlers.InvalidArgsNumberError{Command: "RENAME"}.Error()},
	{[]string{"*3\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n$6\r\nnewkey\r\n"}, "key to be renamed not found"},
	{[]string{"*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n", "*2\r\n$3\r\nget\r\n$5\r\nhello\r\n"}, "world"},
	{[]string{"*3\r\n$3\r\nSET\r\n$10\r\ndeletedkey\r\n$13\r\ndeletedvalue\r\n", "*2\r\n$3\r\nDEL\r\n$10\r\ndeletedkey\r\n"}, ":1"},
	{[]string{"*2\r\n$6\r\nEXISTS\r\n$7\r\nmissing\r\n"}, ":0"},
	{[]string{"*3\r\n$3\r\nset\r\n$6\r\nthekey\r\n$5\r\nfound\r\n", "*2\r\n$6\r\nexists\r\n$6\r\nthekey\r\n"}, ":1"},
	{[]string{"*3\r\n$3\r\nset\r\n$6\r\nlength\r\n$4\r\nfour\r\n", "*2\r\n$6\r\nstrlen\r\n$6\r\nlength\r\n"}, ":4"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$3\r\nval\r\n", "*3\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n$6\r\nnewval\r\n"}, ":0"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$1\r\n1\r\n", "*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n"}, ":2"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$3\r\none\r\n", "*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n"}, handlers.InvalidDataTypeError{Command: "INCR"}.Error()},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$1\r\n1\r\n", "*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n"}, ":0"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$3\r\none\r\n", "*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n"}, handlers.InvalidDataTypeError{Command: "DECR"}.Error()},
	{[]string{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$3\r\nval\r\n", "*3\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n$6\r\nnewkey\r\n", "*2\r\n$3\r\nGET\r\n$6\r\nnewkey\r\n"}, "val"},
	{[]string{"*2\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n"}, handlers.InvalidArgsNumberError{Command: "SETNX"}.Error()},
	{[]string{"*2\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n"}, handlers.InvalidArgsNumberError{Command: "LPUSH"}.Error()},
	{[]string{"*3\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"}, ":1"},
	{[]string{"*3\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$5\r\nvalue\r\n", "*4\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$6\r\nvalue2\r\n$6\r\nvalue3\r\n"}, ":3"},
	{[]string{"*4\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$2\r\nb\r\n", "*2\r\n$4\r\nLPOP\r\n$3\r\nkey\r\n"}, "b"},
	{[]string{"*1\r\n$4\r\nLPOP\r\n"}, handlers.InvalidArgsNumberError{Command: "LPOP"}.Error()},
	{[]string{"*2\r\n$4\r\nLPOP\r\n$3\r\nkey\r\n"}, "_\r\n"},
	{[]string{"*4\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n", "*3\r\n$4\r\nLPOP\r\n$3\r\nkey\r\n$1\r\n2\r\n"}, "*2\r\n$1\r\nb\r\n$1\r\na\r\n"},
	{[]string{"*4\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n", "*3\r\n$4\r\nLPOP\r\n$3\r\nkey\r\n$1\r\n3\r\n"}, "*2\r\n$1\r\nb\r\n$1\r\na\r\n"},
	{[]string{"*4\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n", "*3\r\n$4\r\nLPOP\r\n$3\r\nkey\r\n$1\r\na\r\n"}, handlers.InvalidDataTypeError{Command: "LPOP"}.Error()},
	{[]string{"*4\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n", "*2\r\n$4\r\nLLEN\r\n$3\r\nkey\r\n"}, ":2"},
	{[]string{"*2\r\n$4\r\nLLEN\r\n$3\r\nkey\r\n"}, ":0"},
	{[]string{"*3\r\n$4\r\nLLEN\r\n$3\r\nkey\r\n$5\r\nvalue\r\n"}, handlers.InvalidArgsNumberError{Command: "LLEN"}.Error()},
	{[]string{"*4\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n", "*3\r\n$4\r\nLPOP\r\n$3\r\nkey\r\n$1\r\n2\r\n"}, "*2\r\n$1\r\na\r\n$1\r\nb\r\n"},
	{[]string{"*4\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n", "*3\r\n$4\r\nLPOP\r\n$3\r\nkey\r\n$1\r\n3\r\n"}, "*2\r\n$1\r\na\r\n$1\r\nb\r\n"},
	{[]string{"*4\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n", "*2\r\n$4\r\nLLEN\r\n$3\r\nkey\r\n"}, ":2"},
	{[]string{"*2\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n"}, handlers.InvalidArgsNumberError{Command: "RPUSH"}.Error()},
	{[]string{"*5\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n", "*2\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n"}, handlers.InvalidArgsNumberError{Command: "LINDEX"}.Error()},
	{[]string{"*5\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\n3\r\n"}, "_\r\n"},
	{[]string{"*5\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$2\r\n-4\r\n"}, "_\r\n"},
	{[]string{"*5\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\n1\r\n"}, "$1\r\nb\r\n"},
	{[]string{"*5\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$2\r\n-3\r\n"}, "$1\r\na\r\n"},
	{[]string{"*5\r\n$5\r\nRPUSH\r\n$3\r\nkey\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\na\r\n"}, handlers.InvalidDataTypeError{Command: "LINDEX"}.Error()},
	{[]string{"*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$2\r\n-3\r\n"}, "key not found"},
	{[]string{"*3\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\n1\r\n"}, handlers.InvalidArgsNumberError{Command: "LTRIM"}.Error()},
	{[]string{"*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\n1\r\n$2\r\n-1\r\n"}, "key not found"},
	{[]string{"*3\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\b$1\r\na\r\n", "*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\n1\r\n$1\r\na\r\n"}, handlers.InvalidDataTypeError{Command: "LTRIM"}.Error()},
	{[]string{"*3\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\b$1\r\na\r\n", "*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\na\r\n$2\r\n-1\r\n"}, handlers.InvalidDataTypeError{Command: "LTRIM"}.Error()},
	{[]string{"*5\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\b$1\r\nc\r\n$1\r\nb\r\n$1\r\na\r\n", "*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\n1\r\n$2\r\n-1\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\n0\r\n"}, "$1\r\nb\r\n"},
	{[]string{"*5\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\b$1\r\nc\r\n$1\r\nb\r\n$1\r\na\r\n", "*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\n0\r\n$2\r\n-1\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\n0\r\n"}, "$1\r\na\r\n"},
	{[]string{"*5\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\b$1\r\nc\r\n$1\r\nb\r\n$1\r\na\r\n", "*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\n1\r\n$2\r\n-2\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\n0\r\n"}, "$1\r\nb\r\n"},
	{[]string{"*5\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\b$1\r\nc\r\n$1\r\nb\r\n$1\r\na\r\n", "*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$1\r\n2\r\n$2\r\n-1\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\n0\r\n"}, "$1\r\nc\r\n"},
	{[]string{"*5\r\n$5\r\nLPUSH\r\n$3\r\nkey\r\b$1\r\nc\r\n$1\r\nb\r\n$1\r\na\r\n", "*4\r\n$5\r\nLTRIM\r\n$3\r\nkey\r\n$2\r\n10\r\n$2\r\n-2\r\n", "*3\r\n$6\r\nLINDEX\r\n$3\r\nkey\r\n$1\r\n0\r\n"}, "key not found"},
}

func TestCommands(t *testing.T) {
	for _, test := range CommandTableTests {

		var conn ConnectionMock
		for _, command := range test.Commands {
			conn = GetConnectionMock(command)
			HandleConnection(conn)
		}

		if !strings.Contains(conn.String(), test.ExpectedOutput) {
			t.Errorf("expected result '%s' got response of '%s'", test.ExpectedOutput, conn.String())
			for _, command := range test.Commands {
				t.Logf(command)
			}
		}

		handlers.ClearData()
	}
}
