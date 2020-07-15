package main

import (
	"fmt"
	"log"
	"net"
)


func main() {

	/*
	   Listen: 返回在一个本地网络地址laddr上监听的Listener。网络类型参数net必须是面向流的网络： "tcp"、"tcp4"、"tcp6"、"unix"或"unixpacket"。
	*/
	listener, err := net.Listen("tcp", "localhost:9999")
	if err != nil {
		fmt.Println("服务坚听端口失败:", err)
		return
	}
	defer listener.Close()
	fmt.Println("等待客户端接入")
	for {
		//等待客户端接入
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("客户端链接失败：", err)
			break
		}
		fmt.Println(conn.RemoteAddr().String(), "已链接:")
		// 使用协程
		go handleConnection(conn)
	}
}

// 读取数据
func handleConnection(conn net.Conn) {
	defer conn.Close()
	strChan := make(chan string)
	structChan := make(chan struct{})
	go inPut(strChan)
	go outPut(conn, structChan)

	for {
		select {
		case str := <- strChan:
			conn.Write([]byte(str))
		case <- structChan:
			fmt.Println(conn.RemoteAddr().String(), "已断开链接")
			return
		}
	}

}

func inPut(c chan string)  {
	fmt.Println("等待输入...")
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