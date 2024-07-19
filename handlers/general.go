package handlers

import (
	"container/list"

	"github.com/idugan100/GoKV/resp"
)

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
	hsetMU.Lock()
	hsetData = map[string]map[string]string{}
	hsetMU.Unlock()
	listMU.Lock()
	listData = map[string]*list.List{}
	listMU.Unlock()
	return resp.Serializable{Typ: "string", Str: "OK"}
}

func dbsize(args []resp.Serializable) resp.Serializable {
	setMU.RLock()
	defer setMU.RUnlock()
	listMU.RLock()
	defer listMU.RUnlock()
	hsetMU.RLock()
	defer hsetMU.RUnlock()
	total_keys := len(setData) + len(listData) + len(hsetData)

	return resp.Serializable{Typ: "integer", Num: total_keys}
}
