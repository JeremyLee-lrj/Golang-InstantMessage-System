package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	OnlineMap map[string]*User
	myLock    sync.RWMutex

	Message chan string
}

// create a API of server
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		myLock:    sync.RWMutex{},
		Message:   make(chan string),
	}
	return server
}

func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message
		server.myLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.C <- msg
		}
		server.myLock.Unlock()
	}
}

func (server *Server) BroadCast(user *User, msg string) {
	sendMessage := fmt.Sprintf("[%s]%s: %s", user.Addr, user.Name, msg)
	server.Message <- sendMessage
}
func (server *Server) Handler(conn net.Conn) {
	// fmt.Println("Successfully connect")
	user := NewUser(conn)
	server.myLock.Lock()
	server.OnlineMap[user.Name] = user
	server.myLock.Unlock()

	server.BroadCast(user, "Online")
}

func (server *Server) Start() {
	listener, error := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if error != nil {
		fmt.Println("net.listener error:", error)
		return
	}
	defer listener.Close()

	go server.ListenMessager()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}

		go server.Handler(conn)
	}
}
