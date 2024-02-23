package main

import "net"

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	// listen user's channel
	go user.ListenMessage()
	return user
}
func (user *User) Online() {
	user.server.myLock.Lock()
	user.server.OnlineMap[user.Name] = user
	user.server.myLock.Unlock()

	user.server.BroadCast(user, "Online")
}

func (user *User) Offline() {
	user.server.myLock.Lock()
	delete(user.server.OnlineMap, user.Name)
	user.server.myLock.Unlock()

	user.server.BroadCast(user, "Offline")
}

func (user *User) DoMessage(msg string) {
	user.server.BroadCast(user, msg)
}
func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
