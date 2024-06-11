package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/idugan100/GoKV/handlers"
	"github.com/idugan100/GoKV/resp"
)

func main() {
	port := flag.String("port", "6379", "use this flag to set the port the GoKV server will run on")
	flag.Parse()
	startServer(*port)
}

func startServer(port string) {
	fmt.Println("starting GoKV server ...")
	l, err := net.Listen("tcp", ":"+port)
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

		go HandleConnection(conn)
	}
}

func HandleConnection(conn io.ReadWriteCloser) {
	defer conn.Close()

	for {
		r := resp.NewResp(conn)

		serializable, err := r.Read()
		if err != nil {
			if err == io.EOF {
				return
			}
			conn.Write(resp.Serializable{Typ: "error", Str: err.Error()}.Marshal())
			continue
		}

		serializableErr, ok := serializable.ValidateIncoming()
		if !ok {
			conn.Write(serializableErr.Marshal())
		}

		handler, ok := handlers.Handlers[strings.ToUpper(serializable.Array[0].Bulk)]

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
