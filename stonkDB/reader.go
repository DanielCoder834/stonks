package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

// Ex. response: "$5\r\nValue\r\n"

// Basic type aliasing
const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Value struct {
	typ   string  // For representing types eg. int, string, or...
	str   string  // For representing smaller informational strings eg. \r or \n
	num   int     // Number of characters in a given string or number of arguments
	bulk  string  // The big strings or information
	array []Value // The values
}

type Resp struct {
	reader *bufio.Reader
}

// NewResp Makes a new reader
func NewResp(rd io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(rd)}
}

// Reads and translates the response
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Printf("Unknown type: %v", string(_type))
		return Value{}, nil
	}
}

// If it reaches the end through \r, it stops.
// Reads each new line
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

// Parses the string to an int and returns that value
func (r *Resp) readInteger() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}
	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, n, err
	}
	return int(i64), n, nil
}

// Reads each item in the line and delegates to its respective methods
func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = "array"

	// read length of array
	lengthOfInts, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	// foreach line, parse and read the value
	v.array = make([]Value, 0)
	for i := 0; i < lengthOfInts; i++ {
		val, err := r.Read()
		if err != nil {
			return v, err
		}

		// append parsed value to array
		v.array = append(v.array, val)
	}

	return v, nil
}

// Reads bulk string
func (r *Resp) readBulk() (Value, error) {
	v := Value{}

	v.typ = "bulk"

	lengthOfInts, _, err := r.readInteger()
	if err != nil {
		return v, err
	}

	bulk := make([]byte, lengthOfInts)

	r.reader.Read(bulk)

	v.bulk = string(bulk)

	// Read the trailing CRLF
	r.readLine()

	return v, nil
}
