package main

import "net"

type User struct {
	Name string
	Addr string
	Ch   chan string
	Conn net.Conn

	Server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	username := conn.RemoteAddr().String()

	user := &User{
		Name:   username,
		Addr:   username,
		Ch:     make(chan string),
		Conn:   conn,
		Server: server,
	}

	go user.ListenMessage()

	return user
}

func (u *User) Online() {
	//用户上线，添加到用户表中
	u.Server.Maplock.Lock()
	u.Server.OnlineMap[u.Name] = u
	u.Server.Maplock.Unlock()

	//广播当前用户上线信息
	u.Server.BroadCast(u, "已上线")
}

func (u *User) Offline() {
	//用户上线，添加到用户表中
	u.Server.Maplock.Lock()
	delete(u.Server.OnlineMap, u.Name)
	u.Server.Maplock.Unlock()

	//广播当前用户上线信息
	u.Server.BroadCast(u, "已下线")
}

func (u *User) DoMessage(mes string) {
	u.Server.BroadCast(u, mes)
}

func (u *User) ListenMessage() {
	for {
		mesg := <-u.Ch //一直接收消息，或阻塞

		u.Conn.Write([]byte(mesg + "\n")) //写到客户端
	}
}
