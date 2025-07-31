package main

import "net"

type User struct {
	Name string
	Addr string
	Ch   chan string
	Conn net.Conn
}

func NewUser(conn net.Conn) *User {
	username := conn.RemoteAddr().String()

	user := &User{
		Name: 	username,
		Addr: 	username,
		Ch: 	make(chan string),
		Conn: 	conn,
	}

	go user.ListenMessage()

	return user
}

func (u *User) ListenMessage() {
	for {
		mesg := <-u.Ch	//一直接收消息，或阻塞

		u.Conn.Write([]byte(mesg+"\n"))  //写到客户端
	}
}