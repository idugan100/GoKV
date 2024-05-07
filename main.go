package main

import (
	"fmt"
	"net"
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
		conn.Write([]byte("+OK\r\n"))
	}

}
