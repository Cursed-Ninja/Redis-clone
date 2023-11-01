package resp

import (
	"bufio"
	"io"
	"log"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string
	str   string
	num   int
	bulk  string
	array []Value
}

type Resp struct {
	reader *bufio.Reader
}


func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, size int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		size++
		line = append(line, b)
		if size >= 2 && line[size - 2] == '\r' {
			break
		}
	}

	return line[:size - 2], size, nil
}

func (r *Resp) readInteger() (int, int, error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, nil
	} 

	num, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}

	return int(num), n, nil
}

func (r *Resp) readArray() (Value, error) {
	val := Value{}
	val.typ = "array"

	// read length of array
	len, _, err := r.readInteger()
	if err != nil {
		return val, err
	}

	val.array = make([]Value, len)

	for i := 0; i < len; i++ {
		value, err := r.Read()
		if err != nil {
			return val, err
		}

		val.array[i] = value
	}

	return val, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = "bulk"

	len, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, len)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// Reading the remaining CRLF
	r.readLine()

	return v, nil
}

func (r *Resp) Read() (Value, error) {
	input_type, err := r.reader.ReadByte()

	if err != nil {
		log.Println("Error reading the data: ", err.Error())
		return Value{}, nil
	}

	switch input_type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		log.Println("Unknown input type: ", string(input_type))
		return Value{}, nil
	}
}

