package tcp_test

import (
	"fmt"
	"net"
	"strconv"
)

func TcpServerStart() {

	fmt.Println("hello world")

	lner, err := net.Listen("tcp", "localhost:8888")

	if err != nil {
		fmt.Println("listener creat error", err)
	}

	fmt.Println("waiting for client")

	for {
		// 等待连接，阻塞函数
		conn, err := lner.Accept()
		if err != nil {
			fmt.Println("accept error", err)
		}
		go handleClientConnection(conn)
	}

}

func handleClientConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("connection success")
	fmt.Println("client address: ", conn.RemoteAddr())

	data_buffer := make([]byte, NetMessageMaxLength+1)

	ClientMessage := MessageStruct{}

	for {
		recvLen, err := conn.Read(data_buffer)
		if err != nil {
			fmt.Println("Read error", err)
			return
		}

		used_num := ClientMessage.AppendMsgData(data_buffer, recvLen)

	TestReadyProcess:

		// 如果ready,则处理
		if ClientMessage.IsMsgReady() {
			ProcessNetMessage(conn, &ClientMessage)
		}

		if used_num < recvLen {
			// 说明前面已经拼接完成处理过了，则重新初始化新消息处理
			ClientMessage.Init()
			recvLen = recvLen - used_num
			used_num = ClientMessage.AppendMsgData(data_buffer[used_num:recvLen], recvLen)

			goto TestReadyProcess
		}

	}

}

// 处理网络消息
func ProcessNetMessage(conn net.Conn, ClientMessage *MessageStruct) bool {

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

	ToClientMessage := MessageStruct{}
	ToClientMessage.MakeMessage(65, "Hello client I`m server, welcome back!")

	_, err := conn.Write(ToClientMessage.GetData())

	if err != nil {
		return false
	}

	fmt.Println("send message success")
	//fmt.Println("send message len；", sendLen)

	return true
}
