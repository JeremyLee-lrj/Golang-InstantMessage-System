package main

import (
	"net"
	"strings"
)

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

func (user *User) SendMsg(msg string) {
	user.conn.Write([]byte(msg))
}

func (user *User) DoMessage(msg string) {
	if msg == "who" {
		// Query all online users
		user.server.myLock.Lock()
		for _, onlineuser := range user.server.OnlineMap {
			onlineMsg := "[" + onlineuser.Addr + "]" + onlineuser.Name + ": Online\n"
			user.SendMsg(onlineMsg)
		}
		user.server.myLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		// Check if newName is existing
		_, ok := user.server.OnlineMap[newName]
		if ok {
			user.SendMsg("New Name is already using\n")
		} else {
			user.server.myLock.Lock()
			delete(user.server.OnlineMap, user.Name)
			user.server.OnlineMap[newName] = user
			user.server.myLock.Unlock()
			user.Name = newName
			user.SendMsg("Rename to " + newName + " successfully!\n")
		}
	} else {
		user.server.BroadCast(user, msg)
	}
}

func (user *User) ListenMessage() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\n"))
	}
}
