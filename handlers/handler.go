package handlers

import (
	"fmt"
	"sync"

	"github.com/idugan100/GoKV/resp"
)

var setData = map[string]string{}
var setMU = sync.RWMutex{}
var hsetData = map[string]map[string]string{}
var hsetMU = sync.RWMutex{}
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
	"DECRBY":    decrby,
	"LOLWUT":    lolwut,
	"FLUSHALL":  flushall,
	"GETSET":    getset,
	"RENAME":    rename,
	"HSET":      hset,
	"HGET":      hget,
	"HEXISTS":   hexists,
	"HSTRLEN":   hstrlen,
	"HLEN":      hlen,
	"HGETALL":   hgetall,
	"HSETNX":    hsetnx,
	"HDEL":      hdel,
}

type InvalidArgsNumberError struct {
	Command string
}

func (i InvalidArgsNumberError) Error() string {
	return fmt.Sprintf("Invalid number of arguments for '%s' command", i.Command)
}

type InvalidDataTypeError struct {
	Command string
}

func (i InvalidDataTypeError) Error() string {
	return fmt.Sprintf("Data type for '%s' command", i.Command)
}

func ClearData() {
	setData = map[string]string{}
	hsetData = map[string]map[string]string{}
}
