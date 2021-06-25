//provides the server of a simple chatroom

package main

import (
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"
)

const (
	thousand = 1000
	lostRemoteError = ".wsarecv: An existing connection was forcibly closed by the remote host"
	list = "LS"
)

var account = make(map[string]net.Conn)
var mesBase = make(chan [2]string,thousand)


//Record keep a record of the login info
func Record(conn net.Conn){
	account[conn.RemoteAddr().String()] = conn
}



//handleConnection receives a connection and reads the information it carries
func handleConnection(conn net.Conn) {

	text := make([] byte, thousand)

	defer conn.Close()

	for {
		size, err := conn.Read(text)
		uid := conn.RemoteAddr().String()
		if err == nil {
			var info [2] string = [2] string {uid,string(text[0:size])}
			mesBase <- info
		}else if err == io.EOF {
			continue
		}else if lostRemote,_ := regexp.MatchString(lostRemoteError,err.Error());lostRemote {
			fmt.Printf("Remote user %s forcibly closed the connection\n",uid)
			delete(account,uid)
			return
		}else {
			panic(err)
		}



	}

}

//MassMessage send message to all the user in the queue
func MassMessage(text string) error{
	for _,con := range account {

		_, err := con.Write([]byte(text))

		if err != nil {
			 return fmt.Errorf("MassMessage Failure, at user %s.\n",con.RemoteAddr().String())
		}
	}
	return nil
}

//DistributeMes if there is an untouched message, distribute it to the intended user.
func DistributeMes(){
	for {
		select {
			case info := <-mesBase:
				num := strings.Count(info[1], ">")
				if num == 0 {
					err := MassMessage(info[1])
					fmt.Println(info[1])
					if err != nil {
						fmt.Println(err)
					}

				} else {
					contents := strings.SplitN(info[1], ">", 2)
					uid := contents[0]
					message := contents[1]
					if user, ok := account[uid]; ok {
						_, err := user.Write([]byte(message))

						if err != nil {
							fmt.Errorf("Failure to write to user %s.\n", uid)
						}

					}
				}
		}
	}
}

func main(){

	// deploy a socket which listens to clients' requests

	listen_socket, err := net.Listen("tcp","127.0.0.1:8080")

	if err != nil{
		fmt.Errorf("connection failure")
	}

	defer listen_socket.Close()

	go DistributeMes()

	//once receiving, start a goroutine to process the information
	for{
		conn, err := listen_socket.Accept()
		if err != nil{
			fmt.Errorf(err.Error())
		}

		Record(conn)
		go handleConnection(conn)
	}
}