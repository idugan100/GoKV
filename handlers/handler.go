package handlers

import (
	"container/list"
	"fmt"
	"sync"

	"github.com/idugan100/GoKV/resp"
)

var setData = map[string]string{}
var setMU = sync.RWMutex{}
var hsetData = map[string]map[string]string{}
var hsetMU = sync.RWMutex{}
var listData = map[string]list.List{}
var listMU = sync.RWMutex{}
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
	"INCRBY":    incrby,
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
	"LPUSH":     lpush,
}

type InvalidArgsNumberError struct {
	Command string
}

func (i InvalidArgsNumberError) Error() string {
	return fmt.Sprintf("invalid number of arguments for '%s' command", i.Command)
}

type InvalidDataTypeError struct {
	Command string
}

func (i InvalidDataTypeError) Error() string {
	return fmt.Sprintf("data type for '%s' command", i.Command)
}

func ClearData() {
	setData = map[string]string{}
	hsetData = map[string]map[string]string{}
	listData = map[string]list.List{}

}
