package main

import (
	"net"
	"strings"
)

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

func (u *User) SendMes(mes string) {
	u.Conn.Write([]byte(mes))
}

func (u *User) DoMessage(mes string) {
	if mes == "who" {
		u.Server.Maplock.Lock()
		for _, user := range u.Server.OnlineMap {
			onlineMes := "[" + user.Addr + "]" + user.Name + "在线...\n"
			u.SendMes(onlineMes)
		}
		u.Server.Maplock.Unlock()

	} else if len(mes) > 7 && mes[:7] == "rename|" {
		//消息格式：rename|张三
		newName := strings.Split(mes, "|")[1]

		//rename 是否存在
		_, ok := u.Server.OnlineMap[newName]
		if ok {
			u.SendMes("该用户名已被使用\n")
		} else {
			u.Server.Maplock.Lock()
			delete(u.Server.OnlineMap, u.Name)
			u.Server.OnlineMap[newName] = u
			u.Server.Maplock.Unlock()

			u.Name = newName

			u.SendMes("用户名修改成功：" + u.Name + "\n")
		}

	} else if len(mes) > 5 && mes[:3] == "to|" {
		//私聊格式：to|张三|你好啊
		remoteName := strings.Split(mes, "|")[1]
		remoteUser, ok := u.Server.OnlineMap[remoteName]

		if !ok {
			//用户名不存在
			u.SendMes("当前用户名不存在\n")
			return
		}

		remoteMes := strings.Split(mes, "|")[2]
		if remoteMes == "" {
			u.SendMes("无消息内容，清重发\n")
			return
		}

		//把消息私法给remoteUser
		remoteUser.SendMes(u.Name + "向您说：" + remoteMes + "\n")

	} else {
		u.Server.BroadCast(u, mes)
	}

}

func (u *User) ListenMessage() {
	for {
		mes := <-u.Ch //一直接收消息，或阻塞

		u.Conn.Write([]byte(mes + "\n")) //写到客户端
	}
}
