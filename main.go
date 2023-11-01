package main

import (
	"log"
	"net"
	"github.com/Cursed-Ninja/Redis-clone/resp"
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
		resp := resp.NewResp(conn)
		value, err := resp.Read()
		if err != nil {
			log.Fatalln("Error reading from client: ", err.Error())
		}

		log.Println(value)

		// ignore request and send back a PONG
		conn.Write([]byte("+OK\r\n"))
	}
}