package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleClient(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		command, err := readRESPArray(reader)

		if err != nil {
			fmt.Println("Error reading command: ", err.Error())
			break
		}

		if len(command) == 0 {
			continue
		}

		switch strings.ToUpper(command[0]) {
		case "PING":
			handlePing(writer)
		case "ECHO":
			if len(command) > 1 {
				handleEcho(writer, command[1])
			} else {
				handleError(writer, "ERR wrong number of arguments for 'Echo' command")
			}
		default:
			handleError(writer, "ERR unknown command '"+command[0]+"'")
		}

		err = writer.Flush()
		if err != nil {
			fmt.Println("Error flushing writer: ", err.Error())
			break
		}

	}
}

func readRESPArray(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, fmt.Errorf("Invalid RESP array")
	}

	numArgs, err := strconv.Atoi(strings.TrimSpace(line[1:]))
	if err != nil {
		return nil, err
	}

	args := make([]string, numArgs)
	for i := 0; i < numArgs; i++ {
		line, err = reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if line[0] != '$' {
			return nil, fmt.Errorf("Invalid RESP array")
		}

		argLen, err := strconv.Atoi(strings.TrimSpace(line[1:]))
		if err != nil {
			return nil, err
		}

		arg := make([]byte, argLen)
		_, err = io.ReadFull(reader, arg)
		if err != nil {
			return nil, err
		}

		_, err = reader.Discard(2)
		if err != nil {
			return nil, err
		}

		args[i] = string(arg)
	}
	return args, nil
}

func handlePing(writer *bufio.Writer) {
	writer.WriteString("+PONG\r\n")
}

func handleEcho(writer *bufio.Writer, message string) {
	writer.WriteString("$" + strconv.Itoa(len(message)) + "\r\n" + message + "\r\n")
}

func handleError(writer *bufio.Writer, errMessage string) {
	writer.WriteString("-" + errMessage + "\r\n")
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			break
		}

		go handleClient(conn)
	}
}
