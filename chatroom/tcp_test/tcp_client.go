package tcp_test

import (
	"fmt"
	"net"
	"strconv"
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

	go handleServerConnection(conn)

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

func handleServerConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("connection success")
	fmt.Println("server address: ", conn.RemoteAddr())

	data_buffer := make([]byte, NetMessageMaxLength+1)

	ServerMessage := MessageStruct{}

	for {
		recvLen, err := conn.Read(data_buffer)
		if err != nil {
			fmt.Println("Read error", err)
			return
		}

		used_num := ServerMessage.AppendMsgData(data_buffer, recvLen)

	TestReadyProcess:

		// 如果ready,则处理
		if ServerMessage.IsMsgReady() {
			ClientProcessNetMessage(conn, &ServerMessage)
		}

		if used_num < recvLen {
			// 说明前面已经拼接完成处理过了，则重新初始化新消息处理
			ServerMessage.Init()
			recvLen = recvLen - used_num
			used_num = ServerMessage.AppendMsgData(data_buffer[used_num:recvLen], recvLen)

			goto TestReadyProcess
		}

	}

}

// 处理网络消息
func ClientProcessNetMessage(conn net.Conn, ClientMessage *MessageStruct) bool {

	// strBuffer := string(buffer[:recvLen])
	fmt.Println("Message: ", ClientMessage.GetData())
	fmt.Println("Message len :", ClientMessage.CurDataSize)
	// time.Sleep(time.Second * 1)                                                  //等一秒钟，可以看出client里面的read函数有阻塞效果
	// sendLen, err := conn.Write([]byte("I am server, you message :" + strBuffer)) //将client发过来的消息原样发送回去
	// if err != nil {
	// 	fmt.Println("send message error", err)
	// }

	fmt.Println("MessageID: ", ClientMessage.header.MsgID)

	for i := int(0); i < int(ClientMessage.header.param_num); i++ {

		this_param := ParamElement{}
		this_param.FromBytes(ClientMessage.data_buffer[NetMessageHeaderSize+NetMessageElementSize*i : NetMessageHeaderSize+NetMessageElementSize*(i+1)])

		switch this_param.param_type {

		case BaseTypeBool:
			fmt.Println(ClientMessage.GetBool(i))
			break
		case BaseTypeByte:
			fmt.Println(ClientMessage.GetByte(i))
			break
		case BaseTypeInt8:
			fmt.Println(ClientMessage.GetInt8(i))
			break
		case BaseTypeUint8:
			fmt.Println(ClientMessage.GetUint8(i))
			break

		case BaseTypeInt16:
			fmt.Println(ClientMessage.GetInt16(i))
			break
		case BaseTypeUint16:
			fmt.Println(ClientMessage.GetUint16(i))
			break

		case BaseTypeInt32:
			fmt.Println(ClientMessage.GetInt32(i))
			break
		case BaseTypeUint32:
			fmt.Println(ClientMessage.GetUint32(i))
			break
		case BaseTypeFloat32:
			fmt.Println(ClientMessage.GetFloat32(i))
			break

		case BaseTypeInt64:
			fmt.Println(ClientMessage.GetInt64(i))
			break
		case BaseTypeUint64:
			fmt.Println(ClientMessage.GetUint64(i))
			break
		case BaseTypeFloat64:
			fmt.Println(ClientMessage.GetFloat64(i))
			break
		case BaseTypeComplex64:
			fmt.Println(ClientMessage.GetComplex64(i))
			break

		case BaseTypeComplex128:
			fmt.Println(ClientMessage.GetComplex128(i))
			break

		case BaseTypeString:
			fmt.Println(ClientMessage.GetString(i))
			break
		case BaseTypeBinary:
			fmt.Println(ClientMessage.GetBinary(i))
			break
		default:
			fmt.Println("UnKnow" + strconv.Itoa(int(this_param.param_type)))
			break
		}

	}

	// ToClientMessage := MessageStruct{}
	// ToClientMessage.MakeMessage(65, "Hello client I`m server, welcome back!")

	// _, err := conn.Write(ToClientMessage.GetData())

	// if err != nil {
	// 	return false
	// }

	// fmt.Println("send message success")
	// //fmt.Println("send message len；", sendLen)

	return true
}
