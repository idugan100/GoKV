package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	fmt.Println("starting GoKV server ...")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn io.ReadWriteCloser) {
	defer conn.Close()

	for {
		r := NewResp(conn)
		value, err := r.Read()

		if err != nil {
			if err == io.EOF {
				return
			}
		}
		if value.typ != "array" {
			fmt.Println("Invalid request expected array")
			v := Value{typ: "error", str: "invalid input type expected array"}
			conn.Write(v.Marshal())
			continue
		}
		if len(value.array) <= 0 {
			fmt.Println("Invalid request, no args")
			v := Value{typ: "error", str: "invalid input request no args"}
			conn.Write(v.Marshal())
			continue
		}

		handler, ok := Handlers[strings.ToUpper(value.array[0].bulk)]
		if !ok {
			notFoundVal := Value{typ: "error", str: "command not found"}
			conn.Write(notFoundVal.Marshal())
			continue
		}
		args := value.array[1:]
		rVal := handler(args)
		conn.Write(rVal.Marshal())
	}
}
