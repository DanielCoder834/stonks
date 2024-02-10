package main

import (
	"fmt"
	"net"
	"strings"
)

// Server handles connecting, listening and writes to the program
// ERR: If more than one item, print the message
// DEFER: wait til the function is done
func Server() {
	// Listens to new responses
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Append to file only is meant to serve as a back-up
	// Takes every write in the file and reads them and processes those requests
	aof, err := NewAof("database.aof")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer aof.Close()

	// Reads responses to the file
	aof.Read(func(value Value) {
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("here1")
			fmt.Println("Invalid command: ", command)
			return
		}
		handler(args)
	})
	for {
		// Accepts the connection
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		// Handles responses
		resp := NewResp(conn)
		value, err := resp.Read()

		if err != nil && err.Error() != "EOF" {
			fmt.Println(err)
			return
		}

		// If the value is not an array, I have done-goofed
		if value.typ != "array" {
			fmt.Println("Invalid request, expected an array")
			continue
		}

		// If the value has a length of 0, I have done-goofed
		if len(value.array) == 0 {
			fmt.Println("Invalid request, expected array length > 0")
			continue
		}

		// Command eg. SET, GET, HSET, HGET
		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		if command == "EXIT" {
			break
		}

		// Handles writes
		writer := NewWriter(conn)

		handler, ok := Handlers[command]
		if !ok {
			fmt.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		// Send to the aof
		if command == "SET" || command == "HSET" || command == "ARRSET" {
			aof.Write(value)
		}

		result := handler(args)
		// Write back a response
		writer.Write(result)
		conn.Close()
	}
	//defer connect.Close() // close connection once finished
}
