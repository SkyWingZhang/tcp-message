package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

func main() {

	//阻塞Dial
	/*
	   Dial:
	       在网络network上连接地址address，并返回一个Conn接口。
	       可用的网络类型有："tcp"、"tcp4"、"tcp6"、"udp"、"udp4"、"udp6"、"ip"、"ip4"、"ip6"、"unix"、"unixgram"、"unixpacket"
	       对TCP和UDP网络，地址格式是host:port或[host]:port
	*/
	//conn, err := net.Dial("tcp", "localhost:7777")
	//超时
	conn, err := net.DialTimeout("tcp", "localhost:9999",time.Second*2)
	if err != nil {
		log.Println("链接服务器错误:", err)
		return
	}
	defer conn.Close()
	fmt.Println("已链接服务器")

	strChan := make(chan string)
	structChan := make(chan struct{})
	go inPut(strChan)
	go outPut(conn, structChan)

	for {
		select {
		case str := <- strChan:
			conn.Write([]byte(str))
			if str == "exit" {
				conn.Write([]byte("客户端：断开链接"))
				fmt.Println("已断开")
				return
			}
		case <- structChan:
			fmt.Println(conn.RemoteAddr().String(), "已断开链接")
			return
		}
	}
}

func inPut(c chan string)  {
	fmt.Println("等待输入,输入exit以退出程序")
	str := ""
	for {
		fmt.Scan(&str)
		c <- str
	}
}

func outPut(conn net.Conn, c chan struct{})  {
	for {
		buf := make([]byte,1024)
		if intA,err := conn.Read(buf);err != nil {
			fmt.Println(err)
			c <- struct{}{}
			return
		}else{
			log.Println(string(buf), "--", intA)
			if string(buf) == "exit" {
				c <- struct{}{}
				return
			}
		}

	}
}

