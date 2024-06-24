package handlers

import (
	"container/list"
	"strconv"

	"github.com/idugan100/GoKV/resp"
)

func lpush(args []resp.Serializable) resp.Serializable {
	if len(args) < 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LPUSH"}.Error()}
	}
	listMU.Lock()
	defer listMU.Unlock()

	l, ok := listData[args[0].Bulk]
	if !ok {
		l = list.New()
	}
	for _, item := range args[1:] {
		l.PushFront(item.Bulk)
	}
	listData[args[0].Bulk] = l
	return resp.Serializable{Typ: "integer", Num: l.Len()}
}

func rpush(args []resp.Serializable) resp.Serializable {
	if len(args) < 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "RPUSH"}.Error()}
	}
	listMU.Lock()
	defer listMU.Unlock()

	l, ok := listData[args[0].Bulk]
	if !ok {
		l = list.New()
	}
	for _, item := range args[1:] {
		l.PushBack(item.Bulk)
	}
	listData[args[0].Bulk] = l
	return resp.Serializable{Typ: "integer", Num: l.Len()}
}

func lpop(args []resp.Serializable) resp.Serializable {
	if len(args) < 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LPOP"}.Error()}
	}
	listMU.Lock()
	defer listMU.Unlock()

	l, ok := listData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "null"}
	}

	if len(args) == 1 {
		val := l.Remove(l.Front())
		if l.Len() == 0 {
			delete(listData, args[0].Bulk)
		}
		return resp.Serializable{Typ: "bulk", Bulk: val.(string)}
	} else {
		number, err := strconv.Atoi(args[1].Bulk)
		if err != nil {
			return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LPOP"}.Error()}
		}
		var results []resp.Serializable
		for range number {
			if l.Len() == 0 {
				delete(listData, args[0].Bulk)
				break
			}
			val := l.Remove(l.Front())
			results = append(results, resp.Serializable{Typ: "bulk", Bulk: val.(string)})
		}
		listData[args[0].Bulk] = l

		return resp.Serializable{Typ: "array", Array: results}
	}

}

func llen(args []resp.Serializable) resp.Serializable {
	if len(args) != 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LLEN"}.Error()}
	}
	listMU.RLock()
	defer listMU.RUnlock()

	l, ok := listData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}
	return resp.Serializable{Typ: "integer", Num: l.Len()}
}
