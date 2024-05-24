package main

import (
	"sync"
)

var setData = map[string]string{}
var setMU = sync.RWMutex{}
var HsetData = map[string]map[string]string{}
var HsetMU = sync.RWMutex{}
var Handlers = map[string]func([]Serializable) Serializable{
	"PING":      ping,
	"SET":       set,
	"GET":       get,
	"DEL":       del,
	"RANDOMKEY": randkey,
	"EXISTS":    exists,
	"STRLEN":    strlen,
	"LOLWUT":    lolwut,
	"FLUSHALL":  flushall,
	"GETSET":    getset,
	"HSET":      hset,
	"HGET":      hget,
	"HEXISTS":   hexists,
	"HSTRLEN":   hstrlen,
}
