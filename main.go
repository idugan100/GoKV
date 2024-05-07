package main

import (
	"fmt"
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

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		r := NewResp(conn)
		value, err := r.Read()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(value)
		if value.typ != "array" {
			fmt.Println("Invalid request expected array")
			//should return error at some point
			continue
		}
		if len(value.array) <= 0 {
			fmt.Println("Invalid request, no args")
			//should return error at some point
			continue
		}

		fmt.Println(string(value.Marshal()))
		handler, ok := Handlers[strings.ToUpper(value.array[0].bulk)]
		if !ok {
			fmt.Println("command not found")
			notFoundVal := Value{typ: "string", str: ""}
			conn.Write(notFoundVal.Marshal())
			//shoudl write error at some point
			continue
		}
		args := value.array[1:]
		rVal := handler(args)
		conn.Write(rVal.Marshal())
	}

}
