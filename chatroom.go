//provides the server of a simple chatroom

package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
)

var account = make(map[string]net.Conn)
var mesBase = make(chan string,1000)
var sentinal []byte= []byte("EXIT")

//CheckError to check if errors occur

func CheckError(err error){
	if err != nil {
		panic(err)
	}
}


//Record keep a record of the login info
func Record(conn net.Conn){
	account[conn.RemoteAddr().String()] = conn
}


//handleConnection receives a connection and reads the information it carries
func handleConnection(conn net.Conn){

	text := make([] byte, 1024)

	defer conn.Close()

	for {
		size, err := conn.Read(text)

		if err == io.EOF {
			continue
		}

		//handle the termination situation

		if bytes.Compare(bytes.ToUpper(text), sentinal) == 0 {
			panic("log out")
		}

		CheckError(err)
		mesBase <- string(text[0:size])

	}

}

//DistributeMes if there is an untouched message, distribute it to the intended user.
func DistributeMes(){
	for {
		select {
		case text := <-mesBase:
			contents := strings.Split(text, ">")
			uid := contents[0]
			message := contents[1]
			fmt.Println(account)
			if user, ok := account[uid]; ok {
				_, err := user.Write([]byte(message))

				CheckError(err)
			}
		}
	}
}

func main(){

	// deploy a socket which listens to clients' requests

	listen_socket, err := net.Listen("tcp","127.0.0.1:8080")
	CheckError(err)
	defer listen_socket.Close()


	go DistributeMes()

	//once receiving, start a goroutine to process the information

	for{
		conn, err := listen_socket.Accept()
		CheckError(err)
		Record(conn)
		go handleConnection(conn)
	}
}