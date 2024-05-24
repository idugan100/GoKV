package main

import "github.com/idugan100/GoKV/resp"

func ping(args []resp.Serializable) resp.Serializable {
	return resp.Serializable{Typ: "bulk", Bulk: "PONG"}
}

func lolwut(args []resp.Serializable) resp.Serializable {
	return resp.Serializable{Typ: "bulk", Bulk: "GoKV 0.1 :):):)\r\n"}
}

func flushall(args []resp.Serializable) resp.Serializable {
	setMU.Lock()
	setData = map[string]string{}
	setMU.Unlock()
	HsetMU.Lock()
	HsetData = map[string]map[string]string{}
	HsetMU.Unlock()
	return resp.Serializable{Typ: "string", Str: "OK"}
}
