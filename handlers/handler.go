package handlers

import (
	"sync"

	"github.com/idugan100/GoKV/resp"
)

var setData = map[string]string{}
var setMU = sync.RWMutex{}
var HsetData = map[string]map[string]string{}
var HsetMU = sync.RWMutex{}
var Handlers = map[string]func([]resp.Serializable) resp.Serializable{
	"PING":      ping,
	"SET":       set,
	"GET":       get,
	"DEL":       del,
	"RANDOMKEY": randkey,
	"EXISTS":    exists,
	"STRLEN":    strlen,
	"SETNX":     setnx,
	"INCR":      incr,
	"DECR":      decr,
	"LOLWUT":    lolwut,
	"FLUSHALL":  flushall,
	"GETSET":    getset,
	"HSET":      hset,
	"HGET":      hget,
	"HEXISTS":   hexists,
	"HSTRLEN":   hstrlen,
	"HLEN":      hlen,
	"HGETALL":   hgetall,
	"HSETNX":    hsetnx,
	"HDEL":      hdel,
}

func ClearData() {
	setData = map[string]string{}
	HsetData = map[string]map[string]string{}
}
