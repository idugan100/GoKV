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

var TableTests = []TableTest{
	{[]string{"*1\r\n$3\r\nset\r\n"}, "incorrect number"},
	{[]string{"*1\r\n$3\r\nget\r\n"}, "incorrect number"},
	{[]string{"*2\r\n$3\r\nget\r\n$3\r\nmissing\r\n"}, "_\r\n"},
	{[]string{"*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n"}, "OK"},
	{[]string{"*3\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n$6\r\nnewval\r\n"}, ":1"},
	{[]string{"*1\r\n$4\r\nINCR\r\b"}, "incorrect number of arguments"},
	{[]string{"*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n"}, ":1"},
	{[]string{"*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n"}, ":-1"},
	{[]string{"*1\r\n$4\r\nDECR\r\b"}, "incorrect number of arguments"},
	{[]string{"*1\r\n$6\r\nDECRBY\r\b"}, "incorrect number of arguments"},
	{[]string{"*3\r\n$6\r\nDECRBY\r\n$3\r\nkey\r\n$3\r\nval\r\n"}, "incorrect data type"},
	{[]string{"*3\r\n$6\r\nDECRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n"}, ":-2"},
	{[]string{"*3\r\n$3\r\nSET\r\n$3\r\nnum\r\n$1\r\n5\r\n", "*3\r\n$6\r\nDECRBY\r\n$3\r\nnum\r\n$1\r\n2\r\n"}, ":3"},
	{[]string{"*2\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n"}, "incorrect number of arguments"},
	{[]string{"*3\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n$6\r\nnewkey\r\n"}, "key to be renamed not found"},
	{[]string{"*3\r\n$3\r\nset\r\n$5\r\nhello\r\n$5\r\nworld\r\n", "*2\r\n$3\r\nget\r\n$5\r\nhello\r\n"}, "world"},
	{[]string{"*3\r\n$3\r\nSET\r\n$10\r\ndeletedkey\r\n$13\r\ndeletedvalue\r\n", "*2\r\n$3\r\nDEL\r\n$10\r\ndeletedkey\r\n"}, ":1"},
	{[]string{"*2\r\n$6\r\nEXISTS\r\n$7\r\nmissing\r\n"}, ":0"},
	{[]string{"*3\r\n$3\r\nset\r\n$6\r\nthekey\r\n$5\r\nfound\r\n", "*2\r\n$6\r\nexists\r\n$6\r\nthekey\r\n"}, ":1"},
	{[]string{"*3\r\n$3\r\nset\r\n$6\r\nlength\r\n$4\r\nfour\r\n", "*2\r\n$6\r\nstrlen\r\n$6\r\nlength\r\n"}, ":4"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nkey\r\n$3\r\nval\r\n", "*3\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n$6\r\nnewval\r\n"}, ":0"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$1\r\n1\r\n", "*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n"}, ":2"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$3\r\none\r\n", "*2\r\n$4\r\nINCR\r\n$3\r\nnum\r\n"}, "incorrect data type"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$1\r\n1\r\n", "*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n"}, ":0"},
	{[]string{"*3\r\n$3\r\nset\r\n$3\r\nnum\r\n$3\r\none\r\n", "*2\r\n$4\r\nDECR\r\n$3\r\nnum\r\n"}, "incorrect data type"},
	{[]string{"*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$3\r\nval\r\n", "*3\r\n$6\r\nRENAME\r\n$3\r\nkey\r\n$6\r\nnewkey\r\n", "*2\r\n$3\r\nGET\r\n$6\r\nnewkey\r\n"}, "val"},
	{[]string{"*2\r\n$5\r\nsetnx\r\n$3\r\nkey\r\n"}, "incorrect number of arguments"},
}

func TestSetCommands(t *testing.T) {
	for _, test := range TableTests {

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
