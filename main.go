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
		serializable, err := r.Read()

		if err != nil {
			if err == io.EOF {
				return
			}
		}
		if serializable.typ != "array" {
			fmt.Println("Invalid request expected array")
			v := Serializable{typ: "error", str: "invalid input type expected array"}
			conn.Write(v.Marshal())
			continue
		}
		if len(serializable.array) <= 0 {
			fmt.Println("Invalid request, no args")
			v := Serializable{typ: "error", str: "invalid input request no args"}
			conn.Write(v.Marshal())
			continue
		}

		handler, ok := Handlers[strings.ToUpper(serializable.array[0].bulk)]
		if !ok {
			notFoundVal := Serializable{typ: "error", str: "command not found"}
			conn.Write(notFoundVal.Marshal())
			continue
		}
		args := serializable.array[1:]
		rVal := handler(args)
		conn.Write(rVal.Marshal())
	}
}
