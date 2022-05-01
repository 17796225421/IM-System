package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip   string
	Port int

	// 在线用户列表
	OnlineMap map[string]*User
	mapLock   sync.RWMutex

	// 消息广播的channel
	Message chan string
}

// 构造
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听server channel，广播给所有user
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message

		// 将msg发送给全部在线User
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息
func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	// fmt.Println("连接建立成功")

	// 用户上线，加入OnlineMap
	user := NewUser(conn, this)

	user.Online()

	isLive := make(chan bool)

	// 接收客户端发送的消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}

			if err != nil && err != io.EOF {
				fmt.Println("Conn.Read err:", err)
			}

			// 提取用户消息，去除\n
			msg := string(buf[:n-1])

			// 将得到的消息广播
			user.DoMessage(msg)

			isLive <- true
		}
	}()

	for {
		select {
		// 用一个case来从channel取出数据，一旦取出数据，就说明活跃，本次select就结束，开始下一次for，也就是重置计时器
		case <-isLive:
		case <-time.After(time.Second * 10):
			// 一个case来每10s触发一次超时踢出
			user.SendMsg("长时间不活跃，你被踢了")

			close(user.C)
			conn.Close()
			return
		}
	}
}

// 启动
func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}

	defer listener.Close()

	// 启动server的channel
	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.Accept err:", err)
			continue
		}

		// 启动server对每个user的监听
		go this.Handler(conn)
	}
}
