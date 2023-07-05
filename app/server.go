package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleClient(conn net.Conn, store *Redis) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		values, err := DecodeRESP(reader)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Println("Error reading command: ", err.Error())
			return
		}

		command := values.Array()[0].String
		args := values.Array()[1:]

		switch strings.ToUpper(command) {
		case "SET":
			handleSet(writer, args, store)
		case "GET":
			handleGet(writer, args, store)
		case "PING":
			handlePing(writer)
		case "ECHO":
			handleEcho(writer, args)
		default:
			handleError(writer, "ERR unknown command '"+command+"'")
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing writer: ", err.Error())
			break
		}

	}
}

func handleSet(writer *bufio.Writer, args []Value, store *Redis) {
	if len(args) < 2 {
		handleError(writer, "Arguments are less")
		return
	}
	key := args[0].String
	value := args[1].String

	err := store.Set(key, value)
	if err != nil {
		handleError(writer, "no")
		return
	}

	writer.WriteString("+ok\r\n")
}

func handleGet(writer *bufio.Writer, args []Value, store *Redis) {
	if len(args) < 1 {
		handleError(writer, "Arguments are less")
		return
	}
	key := args[0].String

	value, err := store.Get(key)
	if err != nil {
		writer.WriteString("-" + err.Error() + "\r\n")
		writer.Flush()
		return
	}

	if value == "" {
		writer.WriteString("$-1\r\n")
	} else {
		writer.WriteString("+" + value + "\r\n")
	}
	writer.Flush()
}

func handlePing(writer *bufio.Writer) {
	writer.WriteString("+PONG\r\n")
}

func handleEcho(writer *bufio.Writer, args []Value) {
	var messages []string

	for _, arg := range args {
		messages = append(messages, arg.String)
	}

	response := strings.Join(messages, " ")
	writer.WriteString("$" + strconv.Itoa(len(response)) + "\r\n" + response + "\r\n")
}

func handleError(writer *bufio.Writer, errMessage string) {
	writer.WriteString("-" + errMessage + "\r\n")
}

func main() {

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Listening on :6379 ....")
	defer l.Close()

	store := NewRedis()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}

		go handleClient(conn, store)
	}
}
