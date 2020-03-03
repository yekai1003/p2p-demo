//client.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

//解析地址函数，格式为（ip:port）
func parseAddr(addr string) net.UDPAddr {

	t := strings.Split(addr, ":")
	port, _ := strconv.Atoi(t[1])
	return net.UDPAddr{
		IP:   net.ParseIP(t[0]),
		Port: port,
	}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("./client tag remoteIP remotePort port")
		return
	}
	port, _ := strconv.Atoi(os.Args[4])
	tag := os.Args[1]
	remoteIP := os.Args[2]
	remotePort, _ := strconv.Atoi(os.Args[3])

	//一定要绑定固定端口，否则介绍人不好介绍
	localAddr := net.UDPAddr{Port: port}

	//与服务器建立联系（严格意义上，UDP不能叫连接）
	conn, err := net.DialUDP("udp", &localAddr, &net.UDPAddr{IP: net.ParseIP(remoteIP), Port: remotePort})
	if err != nil {
		log.Panic("Failed ot DialUDP", err)
	}

	//自我介绍，亮明身份，但其实说啥都行
	conn.Write([]byte("我是peer:" + tag))

	buf := make([]byte, 256)
	//从服务器获得目标地址
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		log.Panic("Failed to ReadFromUDP", err)
	}
	conn.Close()
	toAddr := parseAddr(string(buf[:n]))
	fmt.Println("获得对象地址:", toAddr)
	//两个人建立P2P通信
	p2p(&localAddr, &toAddr)

}

func p2p(srcAddr *net.UDPAddr, dstAddr *net.UDPAddr) {
	//请求与对方建立联系
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println("Failed to DialUDP", err)
	}
	fmt.Println("发送打洞数据")
	//defer conn.Close()
	if _, err := conn.Write([]byte("我要打洞\n")); err != nil {
		fmt.Println("send msg err", err)
	}
	//启动一个goroutine监控标准输入
	go func() {
		buf := make([]byte, 256)
		for {
			n, _, err := conn.ReadFromUDP(buf)
			if err != nil {
				//fmt.Println("Failed to ReadFromUDP", err)
				//break
			}
			if n > 0 {
				fmt.Printf("收到消息:%sp2p>", buf[:n])
			}

		}
	}()
	//接下来监控标准输入，发送给对方
	reader := bufio.NewReader(os.Stdin)
	//buf := make([]byte, 256)
	for {
		fmt.Printf("p2p>") //n, err := os.Stdin.Read(buf)
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Panic("Failed to read from std", err)
		}
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println("Failed to Write", err)
			continue
		}
	}

}
