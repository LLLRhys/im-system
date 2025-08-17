package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	IP   string
	Port int

	//在线用户的列表
	OnlineMap map[string]*User
	Maplock   sync.RWMutex

	//用于广播的管道
	Ch chan string
}

// 创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Ch:        make(chan string),
	}
	return server
}

// 准备上线信息的方法
func (s *Server) BroadCast(user *User, mes string) {
	sendmes := "[" + user.Addr + "]" + user.Name + ":" + mes

	s.Ch <- sendmes
}

// 监听server的广播消息，有消息就发送上线信息给所有在线的user
func (s *Server) LisenMessage() {
	for {
		mes := <-s.Ch

		//向用户发送信息
		s.Maplock.Lock()
		for _, user := range s.OnlineMap {
			user.Ch <- mes //会阻塞 要单独开个gorutine
		}
		s.Maplock.Unlock()
	}
}

func (s *Server) Handler(conn net.Conn) {
	//...当前业务的链接
	//fmt.Println("链接建立成功")

	user := NewUser(conn, s)

	//用户的上线部分
	user.Online()

	//监听用户是否活跃的channel
	isLive := make(chan bool)

	//接受客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf) //阻塞，直到有数据到来
			if n == 0 {
				//用户下线
				user.Offline()
				return
			}
			//如果是EOF，说明只是文件末尾
			if err != nil && err != io.EOF {
				fmt.Println("conn.Read error:", err)
				return
			}

			//提取用户发送的信息（最后一位是换行‘/n’）
			mes := string(buf[0 : n-1])

			//用户对mes进行消息处理
			user.DoMessage(mes)

			//用户发送了消息，用户活跃
			isLive <- true
		}
	}() //这个圆括号很重要，前面func(){}只是定义了一个匿名函数，后面的()才是调用这个函数。

	//阻塞当前handle
	for {
		select {
		case <-isLive:
			//当前用户活跃，重置定时器
			//不做任何事，就是为了重新触发下面的定时器

		case <-time.After(time.Second * 180):
			//已经超时  将当前User强制关闭
			user.SendMes("超时，您被强制下线了\n")

			//销毁用的资源
			close(user.Ch)

			//关闭连接
			conn.Close() //客户端下线 还会发送一个”0“，触发User的OffLine

			//conn.Close()之后user还在，只是管道被关闭了，此时还和和还能触发user的OffLine方法
			//retuen之后user才会消失（handler中的局部变量）

			//退出当前的handler
			return //runtime.Goexit()
		}

	}

}

// 启动服务器的接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
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
		conn, err := listener.Accept() //不断阻塞接收，直到有客户端连接。接受后返回客户端对象conn。conn可以理解为一个客户端的连接实例。
		if err != nil {
			fmt.Println("Listener.Accept error:", err)
			continue
		}

		// do handler
		go s.Handler(conn) //客户端对象重新开一个线程。
	}
}
