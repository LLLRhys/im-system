package main

import (
	"flag"
	"fmt"
	"net"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	Conn       net.Conn
	flag       int //当前客户端的状态
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       666,
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

func (c *Client) menu() bool {
	var flag int

	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 4 {
		c.flag = flag
		return true
	} else {
		fmt.Println("请输入合法范围内的数字")
		return false
	}
}

func (c *Client) run() {
	for c.flag != 0 {
		for c.menu() != true {
			//空循环，直到menu()返回正确的状态
		}

		//根据不同的模式处理不同的业务
		switch c.flag {
		case 1:
			//公聊模式
			fmt.Println("选择了公聊模式")

		case 2:
			//私聊模式
			fmt.Println("选择了私聊模式")

		case 3:
			//更新用户名
			fmt.Println("选择了更新用户名")
		}

	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "Ip", "127.0.0.1", "设置服务器Ip地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "Port", 8888, "设置服务器端口（默认是8888）")
}

func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败......")
		return
	}

	fmt.Println("连接服务器成功......")

	//启动客户端业务
	client.run()
}
