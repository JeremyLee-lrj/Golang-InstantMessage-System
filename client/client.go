package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// create Client
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}

	// link server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial() error: ", err)
		return nil
	}

	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1. Public Chat")
	fmt.Println("2. Private Chat")
	fmt.Println("3. Rename")
	fmt.Println("0. Exit")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Please enter legal number!!!")
		return false
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn Write err:", err)
		return
	}
}

func (client *Client) PrivateChat() {
	var remoteName, chatMsg string
	client.SelectUsers()
	fmt.Println("Please enter username, print exit to terminate")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println("Please enter message, print exit to terminate")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn Write err:", err)
					break
				}
			}

			// chatMsg = ""
			fmt.Println("Please enter message, print exit to terminate")
			fmt.Scanln(&chatMsg)
		}
		fmt.Println("Please enter username, print exit to terminate")
		fmt.Scanln(&remoteName)

	}
}

func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println("Please enter message, print exit to terminate")
	// fmt.Scanln(&chatMsg)
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("Please enter message, print exit to terminate")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("Please enter new name")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err: ", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	for client.flag != 0 {
		for !client.menu() {
		}
		switch client.flag {
		case 1:
			// fmt.Println("Public Chat")
			client.PublicChat()
		case 2:
			// fmt.Println("Private Chat")
			client.PrivateChat()
		case 3:
			// fmt.Println("Rename")
			client.UpdateName()
		}
	}
}

var serverIp string
var serverPort int
var wordPtr *string

func init() {
	wordPtr = flag.String("word", "foo", "a string named word")
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "Set ServerIp(default value is \"127.0.0.1\")")
	flag.IntVar(&serverPort, "port", 8888, "Set ServerPort(default value is 8888)")
}

func main() {
	// parse command-line
	Testflag()
	flag.Parse()
	fmt.Println(flag.Arg(0))
	fmt.Printf("word = %v\n", *wordPtr)
	client := NewClient(serverIp, serverPort)
	// print serverIp and serverPort
	// fmt.Printf("serverIp is %s, and serverPort is %d\n", serverIp, serverPort)
	if client == nil {
		fmt.Println("Link server failed...")
		return
	}
	fmt.Println("Link server successfully...")
	go client.DealResponse()
	client.Run()
}
