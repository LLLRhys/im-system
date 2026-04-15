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

func (c *Client) DealResponse() {
	// 永久阻塞监听。不断的读取conn，有数据就copy到stdout标准输出上。
	io.Copy(os.Stdout, c.Conn)

	/*  上下等效
	for {
		buf := make()
		c.Conn.Read(buf)
		fmt.Println(buf)

	}
	*/
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

func (c *Client) SelectUser() {
	mes := "who\n"
	_, err := c.Conn.Write([]byte(mes))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return
	}
}

func (c *Client) PubilcChat() {
	var mes string
	fmt.Println(">>>>>>请输入消息内容,exit退出")
	fmt.Scanln(&mes)

	for mes != "exit" {
		//消息不为空就发送
		if len(mes) != 0 {
			sendMes := mes + "\n"
			_, err := c.Conn.Write([]byte(sendMes))
			if err != nil {
				fmt.Println("conn.Write err :", err)
				break
			}
		}

		//完成初始化，准备第二次发消息
		mes = ""
		fmt.Println(">>>>>>请输入消息内容,exit退出")
		fmt.Scanln(&mes)
	}
}


func (c *Client) PrivateChat() {
	var remoteName string
	var mes string

	c.SelectUser()
	fmt.Println(">>>>>>请输入聊天对象|用户名|，exit退出")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>>>请输入消息内容，exit退出")
		fmt.Scanln(&mes)

		for mes != "exit" {
			//消息不为空就发送
			if mes != "" {
				sendMes := "to|" + remoteName + "|" + mes + "\n"
				_, err := c.Conn.Write([]byte(sendMes))
				if err != nil {
					fmt.Println("conn.Write err:", err)
				}
			}

			mes = ""
			fmt.Println(">>>>>>请输入消息内容，exit退出")
			fmt.Scanln(&mes)
		}

		remoteName = ""
		fmt.Println(">>>>>>请输入聊天对象|用户名|，exit退出")
		fmt.Scanln(&remoteName)
	}
}

func (c *Client) UpdateName() bool {
	fmt.Println(">>>>>>请输入新的用户名：")
	fmt.Scanln(&c.Name)

	//在server.go中的handler()中，对消息的处理会去掉最后一位，即换行符。所以这里要手动加上换行符
	sendMes := "rename|" + c.Name + "\n"

	_, err := c.Conn.Write([]byte(sendMes))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
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
			c.PubilcChat()
		case 2:
			//私聊模式
			c.PrivateChat()
		case 3:
			//更新用户名
			c.UpdateName()
		}

	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	// 参数定义
	flag.StringVar(&serverIp, "Ip", "127.0.0.1", "设置服务器Ip地址（默认是127.0.0.1）")
	flag.IntVar(&serverPort, "Port", 8888, "设置服务器端口（默认是8888）")
}

func main() {
	//命令行解析，读数据 赋值给 serverIp 和 serverPort
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("连接服务器失败......")
		return
	}

	fmt.Println("连接服务器成功......")

	//新开一个goroutine去处理server的回复消息
	go client.DealResponse()

	//启动客户端业务
	client.run()
}
