package main

import (
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
	}

	//连接server
	// 此方式合并Ip和端口，兼容IPv4和IPv6
	address := net.JoinHostPort(serverIp, fmt.Sprintf("%d", serverPort))
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("net.Dial error:", err)
		return nil
	}

	client.Conn = conn

	return client
}

func main() {
	client := NewClient("127.0.0.1", 8888)
	if client == nil {
		fmt.Println("连接服务器失败......")
		return
	}

	fmt.Println("连接服务器成功......")

	//启动客户端业务
	select {}
}
