package main

import (
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/idugan100/GoKV/resp"
)

func main() {
	startServer()
}

func startServer() {
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
		r := resp.NewResp(conn)

		serializable, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
		}
		if serializable.Typ != "array" {
			fmt.Println("Invalid request expected array")
			v := resp.Serializable{Typ: "error", Str: "invalid input type expected array"}
			conn.Write(v.Marshal())
			continue
		}
		if len(serializable.Array) <= 0 {
			fmt.Println("Invalid request, no args")
			v := resp.Serializable{Typ: "error", Str: "invalid input request no args"}
			conn.Write(v.Marshal())
			continue
		}

		handler, ok := Handlers[strings.ToUpper(serializable.Array[0].Bulk)]

		if !ok {
			notFoundVal := resp.Serializable{Typ: "error", Str: "command not found"}
			conn.Write(notFoundVal.Marshal())
			continue
		}
		args := serializable.Array[1:]
		rVal := handler(args)
		conn.Write(rVal.Marshal())
	}
}
