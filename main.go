package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

var (
	isconnexion = make(map[net.Conn]bool)
	connexion   = make(chan net.Conn)
	port        string
)

func main() {

	if len(os.Args) <= 2 {
		switch len(os.Args) {
		case 2:
			arg := os.Args[1]
			val, err := strconv.Atoi(arg)
			if err != nil || val <= 0 {
				fmt.Println("Please check your port")
				return
			} else {
				port = arg
			}
		case 1:
			port = "8989"
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	fmt.Println("Listening on :" + port)

	ltn, err := net.Listen("tcp", ":"+port)
	LogError(err)
	defer ltn.Close()

	go func() {

		for {
			conn, er := ltn.Accept()
			LogError(er)
			isconnexion[conn] = true
			connexion <- conn
		}
	}()
	for {
		go Client(<-connexion)
	}
}
