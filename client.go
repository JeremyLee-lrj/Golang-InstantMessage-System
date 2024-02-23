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
		// Message format: rename|newname
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
	} else if len(msg) > 4 && msg[:3] == "to|" {
		// Message format： to|somebody|message content
		contents := strings.Split(msg, "|")
		remoteName := contents[1]
		if remoteName == "" {
			user.SendMsg("Wrong format! Please use format like this: to|somebody|information\n")
			return
		} else {
			remoteUser, ok := user.server.OnlineMap[remoteName]
			if !ok {
				user.SendMsg("User is no existing!\n")
				return
			}

			message := contents[2]
			if message == "" {
				user.SendMsg("No information, please resend\n")
			} else {
				remoteUser.SendMsg(user.Name + ": " + message + "\n")
			}
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
