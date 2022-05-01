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
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	// 构造，连接服务器

	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       -1,
	}

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn

	return client
}

var serverIp string
var serverPort int

func init() {
	// 解析命令行 ./client -ip 127.0.0.1 -port 8888
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "设置服务器IP地址(默认是127.0.0.1)")
	flag.IntVar(&serverPort, "port", 8888, "设置服务器端口(默认是8888)")
}

// 负责处理服务器消息回应的goroutine，直接打印到屏幕
func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) memu() bool {
	var flag int

	fmt.Println("1.公聊")
	fmt.Println("2.私聊")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")

	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>请输入合法数字<<<<")
		return false
	}
}

func (client *Client) SelectUsers() {
	sendMsg := "who\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err", err)
		return
	}
}

// 私聊，想组装who\n发送给服务器查询所有在线用户，再死循环组装to|zhouzihong|nihao\n发送给服务器
func (client *Client) PrivateChat() {
	var remoteName string
	var chatMsg string

	client.SelectUsers()
	fmt.Println(">>>>请输入聊天对象[用户名]，exit退出：")
	fmt.Scanln(&remoteName)

	for remoteName != "exit" {
		fmt.Println(">>>>请输入消息内容，exit退出：")
		fmt.Scanln(&chatMsg)

		for chatMsg != "exit" {
			if len(chatMsg) != 0 {
				sendMsg := "to|" + remoteName + "|" + chatMsg + "\n\n"
				_, err := client.conn.Write([]byte(sendMsg))
				if err != nil {
					fmt.Println("conn.Write err:", err)
					break
				}
			}
			chatMsg = ""
			fmt.Println(">>>>请输入消息内容，exit退出：")
			fmt.Scanln(&chatMsg)
		}
		client.SelectUsers()
		fmt.Println(">>>>请输入聊天对象[用户名]，exit退出：")
		fmt.Scanln(&remoteName)
	}
}

// 公聊，死循环接收标准输入，只要不是exit，就发送给服务器
func (client *Client) PublicChat() {
	var chatMsg string
	fmt.Println(">>>>输入聊天内容，exit退出")
	fmt.Scanln(&chatMsg)

	for chatMsg != "exit" {
		if len(chatMsg) != 0 {
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("conn.Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println(">>>>输入聊天内容，exit退出")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) updateName() bool {
	// 更新用户名，组装rename|zhouzihong，发送给服务器
	fmt.Println(">>>>请输入用户名")
	fmt.Scanln(&client.Name)

	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("conn.Write err:", err)
		return false
	}

	return true
}

func (client *Client) Run() {
	// 显示菜单，执行相应公聊、私聊、更新用户名业务

	for client.flag != 0 {
		for client.memu() != true {

		}
		switch client.flag {
		case 1:
			// 公聊
			client.PublicChat()
			break
		case 2:
			// 私聊
			client.PrivateChat()
			break
		case 3:
			// 更新用户名
			client.updateName()
			break

		}
	}
}

func main() {
	// 构造客户端连接服务器后，开始业务

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>连接服务器失败")
		return
	}

	go client.DealResponse()

	fmt.Println(">>>>>>连接服务器成功")

	client.Run()
}
