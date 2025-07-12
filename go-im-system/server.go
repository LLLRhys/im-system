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
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Listener.Accept error:", err)
			continue
		}

		// do handler
		go s.Handler(conn)
	}
}