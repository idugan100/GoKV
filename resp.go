package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Serializable struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Serializable
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInteger() (x int, n int, err error) {
	ln, n, err := r.readLine()
	if err != nil {
		return 0, 0, nil
	}
	num, err := strconv.ParseInt(string(ln), 10, 64)
	if err != nil {
		return 0, n, nil
	}

	return int(num), n, nil
}

func (r *Resp) Read() (Serializable, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Serializable{}, err
	}
	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("unknown type %s", string(_type))
		return Serializable{}, nil
	}
}

func (r *Resp) readArray() (Serializable, error) {
	v := Serializable{}
	v.typ = "array"
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	v.array = make([]Serializable, 0)

	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}
		v.array = append(v.array, val)
	}
	return v, nil
}

func (r *Resp) readBulk() (Serializable, error) {
	v := Serializable{}
	v.typ = "bulk"
	length, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, length)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	r.readLine()

	return v, nil

}

func (v Serializable) Marshal() []byte {
	switch v.typ {
	case "array":
		return v.marshalArray()
	case "bulk":
		return v.marshalBulk()
	case "string":
		return v.marshalString()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	case "integer":
		return v.marshalInt()
	default:
		return []byte{}
	}
}

func (v Serializable) marshalInt() []byte {
	var bytes []byte
	bytes = append(bytes, INTEGER)
	bytes = append(bytes, strconv.Itoa(v.num)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Serializable) marshalBulk() []byte {
	var bytes []byte
	bytes = append(bytes, BULK)
	bytes = append(bytes, strconv.Itoa(len(v.bulk))...)
	bytes = append(bytes, '\r', '\n')
	bytes = append(bytes, v.bulk...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Serializable) marshalString() []byte {
	var bytes []byte
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Serializable) marshalArray() []byte {
	var bytes []byte
	length := len(v.array)
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(length)...)
	bytes = append(bytes, '\r', '\n')

	for i := 0; i < length; i++ {
		bytes = append(bytes, v.array[i].Marshal()...)
	}
	return bytes

}

func (v Serializable) marshalError() []byte {
	var bytes []byte
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Serializable) marshalNull() []byte {
	return []byte("_\r\n")
}
