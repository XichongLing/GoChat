package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func sendMes(conn net.Conn){
	var input string

	for {
		reader := bufio.NewReader(os.Stdin)
		data,_,_ := reader.ReadLine()
		input = string(data)

		if strings.Compare(strings.ToUpper(input), "EXIT") == 0 {
			conn.Close()
			break
		}

		_,err := conn.Write([]byte(input))
		if err != nil {
			conn.Close()
			fmt.Println("Connection failure")
			break
		}
	}
}

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	go sendMes(conn)

	buffer := make([]byte,1024)
	for {
		_, err_3 := conn.Read(buffer)
		if err_3 != nil {
			conn.Close()
			break
		}
		fmt.Println(string(buffer))
	}
}

