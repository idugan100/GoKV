package main

import (
	"strings"
	"testing"

	"github.com/idugan100/GoKV/handlers"
)

var CommandTableTests = []TableTest{
	{[]string{"*1\r\n$4\r\nPING\r\n"}, "PONG"},
	{[]string{"*1\r\n$6\r\nLOLWUT\r\n"}, "GoKV 0.1 :):):)"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n"}, ":3"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*6\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nage\r\n$3\r\njob\r\n$3\r\nlol\r\n"}, "isaac"},
	{[]string{"*3\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n"}, "incorrect number"},
	{[]string{"*2\r\n$4\r\nhget\r\n$6\r\nmyinfo\r\n"}, "incorrect number"},
	{[]string{"*2\r\n$7\r\nhexists\r\n$6\r\nmyinfo\r\n"}, "incorrect number"},
	{[]string{"*3\r\n$7\r\nhexists\r\n$6\r\nkey\r\n$5\r\nvalue\r\n"}, ":0"},
	{[]string{"*4\r\n$4\r\nhset\r\n$6\r\nlength\r\n$3\r\nkey\r\n$3\r\nval\r\n", "*3\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n$3\r\nkey\r\n"}, ":3"},
	{[]string{"*2\r\n$7\r\nhstrlen\r\n$6\r\nlength\r\n"}, "incorrect number"},
	{[]string{"*1\r\n$4\r\nHLEN\r\n"}, "incorrect number"},
	{[]string{"*2\r\n$4\r\nHLEN\r\n$6\r\nmissing\r\n"}, ":0"},
	{[]string{"*8\r\n$4\r\nhset\r\n$4\r\ndata\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*2\r\n$4\r\nHLEN\r\n$4\r\ndata\r\n"}, ":3"},
	{[]string{"*1\r\n$7\r\nhgetall\r\n"}, "incorrect number of args"},
	{[]string{"*2\r\n$7\r\nhgetall\r\n$7\r\nmissing\r\n"}, "*0"},
	{[]string{"*8\r\n$4\r\nhset\r\n$4\r\ndata\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*2\r\n$7\r\nhgetall\r\n$4\r\ndata\r\n"}, "swe"},
	{[]string{"*2\r\n$6\r\nHSETNX\r\n$3\r\nkey\r\n"}, "incorrect number of args"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*4\r\n$6\r\nHSETNX\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$3\r\nbob\r\n"}, ":0"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*4\r\n$6\r\nHSETNX\r\n$6\r\nmyinfo\r\n$8\r\nfavcolor\r\n$4\r\nblue\r\n"}, ":1"},
	{[]string{"*2\r\n$4\r\nHDEL\r\ndata\r\n"}, "incorrect number of args"},
	{[]string{"*8\r\n$4\r\nhset\r\n$6\r\nmyinfo\r\n$4\r\nname\r\n$5\r\nisaac\r\n$3\r\nage\r\n$2\r\n20\r\n$3\r\njob\r\n$3\r\nswe\r\n", "*5\r\n$4\r\nHDEL\r\n$6\r\nmyinfo\r\n$8\r\nfavcolor\r\n$4\r\nname\r\n$3\r\nage\r\n"}, ":2"},
	{[]string{"*1\r\n$3\r\n123\r\n"}, "command not found"},
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
