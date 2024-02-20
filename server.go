package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// create a API of server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (server *Server) Handler(conn net.Conn) {
	fmt.Println("Successfully connect")
}

func (server *Server) Start() {
	listener, error := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if error != nil {
		fmt.Println("net.listener error:", error)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		go server.Handler(conn)
	}
}
