package handlers

import (
	"container/list"

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
		l = *list.New()
	}
	for _, item := range args[1:] {
		l.PushFront(item.Bulk)
	}
	listData[args[0].Bulk] = l
	return resp.Serializable{Typ: "integer", Num: l.Len()}
}
