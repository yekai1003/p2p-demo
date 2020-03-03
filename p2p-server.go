package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {
	fmt.Println("begin server")
	//服务器启动侦听
	listener, err := net.ListenUDP("udp", &net.UDPAddr{Port: 9527})
	if err != nil {
		log.Panic("Failed to ListenUDP", err)
	}
	defer listener.Close()
	//定义切片存放2哥udp地址
	peers := make([]*net.UDPAddr, 2, 2)
	buf := make([]byte, 256)
	//接下来从2个UDP消息中获得连接的地址A和B
	n, addr, err := listener.ReadFromUDP(buf)
	if err != nil {
		log.Panic("Failed to ReadFromUDP", err)
	}
	fmt.Printf("read from<%s>:%s\n", addr.String(), buf[:n])
	peers[0] = addr
	n, addr, err = listener.ReadFromUDP(buf)
	if err != nil {
		log.Panic("Failed to ReadFromUDP", err)
	}
	fmt.Printf("read from<%s>:%s\n", addr.String(), buf[:n])
	peers[1] = addr
	fmt.Println("begin nat \n")
	//将A和B分别介绍给彼此
	listener.WriteToUDP([]byte(peers[0].String()), peers[1])
	listener.WriteToUDP([]byte(peers[1].String()), peers[0])
	//睡眠10s确保消息发送完成，可以退出历史舞台
	time.Sleep(time.Second * 10)
}

