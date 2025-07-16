package main

import "fmt" 
import "net"

type Server struct {
	IP string
	Port int
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		IP: ip,
		Port: port,
	}
	return server
}

func (s *Server) Handler(conn net.Conn) {
	//...当前业务的链接
	fmt.Println("链接建立成功")
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