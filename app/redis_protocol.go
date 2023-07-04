package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type RespValue byte

const (
	SimpleString RespValue = '+'
	BulkString   RespValue = '$'
	Array        RespValue = '*'
)

type Value struct {
	typ    RespValue
	String string
	array  []Value
}

func (v Value) Array() []Value {
	if v.typ == Array {
		return v.array
	}

	return []Value{}
}

func DecodeRESP(reader *bufio.Reader) (Value, error) {
	//Read First Char
	dataTypeByte, err := reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch string(dataTypeByte) {
	case "+":
		return decodeSimpleString(reader)
	case "$":
		return decodeBulkString(reader)
	case "*":
		return decodeArray(reader)
	}

	return Value{}, fmt.Errorf("invalid RESP data type byte: %s", string(dataTypeByte))
}

func decodeSimpleString(reader *bufio.Reader) (Value, error) {
	str, err := readUntilCRLF(reader)
	if err != nil {
		return Value{}, err
	}

	return Value{
		typ:    SimpleString,
		String: str,
	}, nil
}

func decodeBulkString(reader *bufio.Reader) (Value, error) {
	str, err := readUntilCRLF(reader)

	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(str)
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	readBytes := make([]byte, count+2)

	if _, err := io.ReadFull(reader, readBytes); err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string contents: %s", err)
	}

	str = string(readBytes[:count])

	return Value{
		typ:    BulkString,
		String: str,
	}, nil
}

func decodeArray(reader *bufio.Reader) (Value, error) {
	readBytesForCount, err := readUntilCRLF(reader)
	if err != nil {
		return Value{}, fmt.Errorf("failed to read bulk string length: %s", err)
	}

	count, err := strconv.Atoi(string(readBytesForCount))
	if err != nil {
		return Value{}, fmt.Errorf("failed to parse bulk string length: %s", err)
	}

	array := []Value{}

	for i := 1; i <= count; i++ {
		value, err := DecodeRESP(reader)
		if err != nil {
			return Value{}, err
		}

		array = append(array, value)
	}

	return Value{
		typ:   Array,
		array: array,
	}, nil
}

func readUntilCRLF(reader *bufio.Reader) (string, error) {
	str, err := reader.ReadString('\r')
	if err != nil {
		return "", err
	}

	// Discard the newline character '\n'
	_, err = reader.ReadByte()
	if err != nil {
		return "", err
	}

	return str[:len(str)-1], nil
}
