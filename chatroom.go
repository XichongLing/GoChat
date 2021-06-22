//provides the server of a simple chatroom

package main

import (
	"fmt"
	"io"
	"net"
)

//CheckError to check if errors occur

func CheckError(err error){
	if err != nil {
		panic(err)
	}
}


//handleConnection recieves a connection and reads the information it carries

func handleConnection(conn net.Conn){

	text := make([] byte, 1024)

	defer conn.Close()

	for {
		_, err := conn.Read(text)

		if err == io.EOF {
			continue
		}
		CheckError(err)
		fmt.Println(string(text))
	}

}


func main(){

	// deploy a socket which listens to clients' requests

	listen_socket, err := net.Listen("tcp","127.0.0.1:8080")
	CheckError(err)
	defer listen_socket.Close()

	//once receiving, start a goroutine to process the information

	for{
		conn, err := listen_socket.Accept()
		CheckError(err)
		go handleConnection(conn)
	}
}