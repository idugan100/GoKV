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

func lindex(args []resp.Serializable) resp.Serializable {
	if len(args) != 2 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LINDEX"}.Error()}
	}
	listMU.RLock()
	defer listMU.RUnlock()
	l, ok := listData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "error", Str: "key not found"}
	}
	index, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LINDEX"}.Error()}
	}

	//check if index is out of range
	if (index > 0 && index > (l.Len()-1)) || (index < 0 && (l.Len()+index) < 0) {
		return resp.Serializable{Typ: "null"}
	}

	var e *list.Element
	if index >= 0 {

		e = l.Front()
		for range index {
			e = e.Next()
		}

	} else {
		e = l.Back()
		for range (index + 1) * -1 {
			e = e.Prev()
		}
	}

	return resp.Serializable{Typ: "bulk", Bulk: e.Value.(string)}
}

func ltrim(args []resp.Serializable) resp.Serializable {
	if len(args) != 3 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LTRIM"}.Error()}
	}

	listMU.Lock()
	defer listMU.Unlock()

	l, ok := listData[args[0].Bulk]
	if !ok {
		return resp.Serializable{Typ: "error", Str: "key not found"}
	}

	size := l.Len()

	startindex, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LTRIM"}.Error()}
	}
	startindex, _ = normalize_index(startindex, size)

	endindex, err := strconv.Atoi(args[2].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LTRIM"}.Error()}
	}
	endindex, _ = normalize_index(endindex, size)

	if startindex > endindex {
		delete(listData, args[0].Bulk)
		return resp.Serializable{Typ: "string", Str: "OK"}
	}

	for i := 0; i < startindex; i++ {
		l.Remove(l.Front())
	}
	for i := endindex; i < (size - 1); i++ {
		l.Remove(l.Back())
	}
	listData[args[0].Bulk] = l

	return resp.Serializable{Typ: "string", Str: "OK"}
}

func rpop(args []resp.Serializable) resp.Serializable {
	if len(args) < 1 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "RPOP"}.Error()}
	}

	listMU.Lock()
	defer listMU.Unlock()

	l, ok := listData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "null"}
	}

	if len(args) == 1 {
		val := l.Remove(l.Back())
		if l.Len() == 0 {
			delete(listData, args[0].Bulk)
		}
		return resp.Serializable{Typ: "bulk", Bulk: val.(string)}
	} else {
		number, err := strconv.Atoi(args[1].Bulk)
		if err != nil {
			return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "RPOP"}.Error()}
		}
		var results []resp.Serializable
		for range number {
			if l.Len() == 0 {
				delete(listData, args[0].Bulk)
				break
			}
			val := l.Remove(l.Back())
			results = append(results, resp.Serializable{Typ: "bulk", Bulk: val.(string)})
		}
		listData[args[0].Bulk] = l

		return resp.Serializable{Typ: "array", Array: results}
	}

}

func lrange(args []resp.Serializable) resp.Serializable {
	if len(args) != 3 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LRANGE"}.Error()}
	}
	listMU.RLock()
	defer listMU.RUnlock()

	l, ok := listData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "array"}
	}
	size := l.Len()

	startindex, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LRANGE"}.Error()}
	}
	startindex, _ = normalize_index(startindex, size)

	endindex, err := strconv.Atoi(args[2].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LRANGE"}.Error()}
	}
	endindex, _ = normalize_index(endindex, size)

	var results []resp.Serializable
	item := l.Front()
	for i := 0; i < startindex; i++ {
		item = item.Next()
	}

	for i := startindex; i <= endindex; i++ {
		results = append(results, resp.Serializable{Typ: "bulk", Bulk: item.Value.(string)})
		item = item.Next()
	}

	return resp.Serializable{Typ: "array", Array: results}

}

func lset(args []resp.Serializable) resp.Serializable {
	if len(args) != 3 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LSET"}.Error()}
	}
	listMU.Lock()
	defer listMU.Unlock()
	l, ok := listData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "error", Str: "key not found"}
	}
	index, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LSET"}.Error()}
	}
	normalized_index, ok := normalize_index(index, l.Len())

	if !ok {
		return resp.Serializable{Typ: "error", Str: "index out of bounds"}
	}

	element := l.Front()

	for range normalized_index {
		element = element.Next()
	}
	element.Value = args[2].Bulk
	return resp.Serializable{Typ: "string", Str: "OK"}
}

func lrem(args []resp.Serializable) resp.Serializable {
	if len(args) != 3 {
		return resp.Serializable{Typ: "error", Str: InvalidArgsNumberError{Command: "LREM"}.Error()}
	}

	count, err := strconv.Atoi(args[1].Bulk)
	if err != nil {
		return resp.Serializable{Typ: "error", Str: InvalidDataTypeError{Command: "LREM"}.Error()}
	}

	listMU.RLock()
	defer listMU.RUnlock()
	l, ok := listData[args[0].Bulk]

	if !ok {
		return resp.Serializable{Typ: "integer", Num: 0}
	}

	element_for_removal := args[2].Bulk
	num_removed := 0
	if count == 0 {
		element := l.Front()
		for range l.Len() {
			next := element.Next()
			if element.Value.(string) == element_for_removal {

				l.Remove(element)
				num_removed++
			}
			element = next
		}

	} else if count > 0 {
		element := l.Front()
		for range l.Len() {
			next := element.Next()
			if element.Value.(string) == element_for_removal {
				l.Remove(element)
				num_removed++
				if num_removed == count {
					break
				}
			}
			element = next
		}
	} else if count < 0 {
		element := l.Back()
		for range l.Len() {
			next := element.Prev()
			if element.Value.(string) == element_for_removal {

				l.Remove(element)
				num_removed++
				if num_removed == count*-1 {
					break
				}
			}
			element = next
		}
	}
	if l.Len() == 0 {
		delete(listData, args[0].Bulk)
	}
	return resp.Serializable{Typ: "integer", Num: num_removed}
}
