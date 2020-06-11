package tcp_test

import (
	"fmt"
	"net"
	"time"
)

func TcpClientStart() {
	fmt.Println("client launch")
	serverAddr := "localhost:8888"
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		fmt.Println("Resolve TCPAddr error", err)
	}
	conn, err := net.DialTCP("tcp4", nil, tcpAddr)

	if err != nil {
		fmt.Println("connect server error", err)
		return
	}

	defer conn.Close()

	message := MessageStruct{}
	message.MakeMessage(125, true, byte(2), int8(-120), uint8(255), int16(-25521), uint16(65533), 18, 3.3, 4.4, "5.55555string", string("6.66666string"))

	fmt.Println("Message: ", message.GetDataBuffer())
	fmt.Println("Message len :", message.header.data_length)

	conn.Write(message.GetDataBuffer())

	go recv(conn)
	time.Sleep(100 * time.Second) //等两秒钟，不然还没接收数据，程序就结束了。
}

func recv(conn net.Conn) {
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err == nil {
		fmt.Println("read message from server:" + string(buffer[:n]))
		fmt.Println("Message len:", n)
	}
}
