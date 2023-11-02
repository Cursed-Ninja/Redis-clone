package main

import (
	"log"
	"net"
	"strings"
)

func main() {
	// setting logger flags
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// setting up the listner
	lsntr, err := net.Listen("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	// accepting connection requests
	conn, err := lsntr.Accept()
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	log.Println("Connection established on port 6379")

	for {
		// read message from client
		resp := NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			log.Println("Error reading from client: ", err.Error())
			continue
		}

		log.Println(value)

		if value.typ != "array" {
			log.Println("Invalid request, expected array")
			continue
		}

		if len(value.array) == 0 {
			log.Println("Invalid request, expected array length > 0")
			continue
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		writer := NewWriter(conn)
		handler, ok := Handlers[command]
		if !ok {
			log.Println("Invalid command: ", command)
			writer.Write(Value{typ: "string", str: ""})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
