package handlers

import (
	"container/list"
	"fmt"
	"sync"
	"time"

	"github.com/idugan100/GoKV/resp"
)

type stringDataItem struct {
	str         string
	expiration  time.Time
	will_expire bool
}

var stringData = map[string]stringDataItem{}
var stringMU = sync.RWMutex{}

func getString(key string) (string, bool) {
	// get item
	stringMU.RLock()
	val, ok := stringData[key]
	stringMU.RUnlock()

	// if item is not found
	if !ok {
		return "", false
	}

	is_expired := val.will_expire && (time.Now().Compare(val.expiration) != -1)

	// if item is expired
	if is_expired {
		stringMU.Lock()
		delete(stringData, key)
		stringMU.Unlock()
		return "", false
	}

	// if item is found and not expired
	return val.str, true
}

var hsetData = map[string]map[string]string{}
var hstringMU = sync.RWMutex{}
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
	"RENAMENX":  renamenx,
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
	"LREM":      lrem,
	"DBSIZE":    dbsize,
	"MGET":      mget,
	"TTL":       ttl,
	"PTTL":      pttl,
	"EXPIRE":    expire,
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
	stringData = map[string]stringDataItem{}
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
