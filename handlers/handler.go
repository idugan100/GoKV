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
var listData = map[string]*list.List{}
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
	"LPOP":      lpop,
	"RPOP":      rpop,
	"LLEN":      llen,
	"RPUSH":     rpush,
	"LINDEX":    lindex,
	"LTRIM":     ltrim,
	"LRANGE":    lrange,
	"LSET":      lset,
	"DBSIZE":    dbsize,
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
	listData = map[string]*list.List{}

}

// returns the normal zero based index and a bool of if the index was out of bounds and was converted to the closest valid index
func normalize_index(index int, size int) (int, bool) {
	in_range := true
	if index < 0 {
		index = size + index
		if index < 0 {
			index = 0
			in_range = false
		}
	}
	if index >= size {
		index = size - 1
		in_range = false
	}
	return index, in_range
}
