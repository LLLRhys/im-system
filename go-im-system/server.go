package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	IP string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	Maplock sync.RWMutex

	//用于广播的管道
	Ch chan string
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP: 		ip,
		Port: 		port,
		OnlineMap: 	make(map[string]*User),
		Ch: 		make(chan string),
	}
	return server
}

//准备上线信息的方法
func (s *Server) BroadCast(user *User, mes string) {
	sendmes := "[" + user.Addr + "]" + user.Name + ":"+ mes

	s.Ch <- sendmes
}

//监听server的广播消息，有消息就发送上线信息给所有在线的user
func (s *Server) LisenMessage() {
	for {
		mes := <- s.Ch

		//向用户发送信息
		s.Maplock.Lock()
		for _,user := range s.OnlineMap {
			user.Ch <- mes  //会阻塞 要单独开个gorutine
		}
		s.Maplock.Unlock()
	}
}

func (s *Server) Handler(conn net.Conn) {
	//...当前业务的链接
	//fmt.Println("链接建立成功")

	user := NewUser(conn)

	//用户上线，添加到用户表中
	s.Maplock.Lock()
	s.OnlineMap[user.Name] = user
	s.Maplock.Unlock()

	//准备当前用户上线信息
	s.BroadCast(user, "你好")

	//阻塞当前handle
	select {}
}

//启动服务器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp",fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("net.Listen error:", err)
		return
	}
	//close socket listen
	defer listener.Close()

	//启动监听server消息的gorutine
	go s.LisenMessage()

	for {
		//accept
		conn, err := listener.Accept()  //不断阻塞接收，直到有客户端连接。接受后返回客户端对象conn。conn可以理解为一个客户端的连接实例。
		if err != nil {
			fmt.Println("Listener.Accept error:", err)
			continue
		}

		// do handler
		go s.Handler(conn)  //客户端对象重新开一个线程。
	}
}