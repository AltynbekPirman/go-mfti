package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	//tcpSrve()
	//httpServeWithFunc()
	go httpServeWithStruct()
	go getParams()
	go login()
	newHTTPServer()
}

func tcpSrve() {
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	name := conn.RemoteAddr().String()
	fmt.Println("connected: ", name)
	defer conn.Close()

	scanner := bufio.NewScanner(conn)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "Exit" {
			_, err := conn.Write([]byte("Bye\n"))
			if err != nil {
				panic(err)
			}
			break
		} else if text != "" {
			_, err := conn.Write([]byte("Entered: " + text + "\n"))
			if err != nil {
				panic(err)
			}
		}
	}
}
