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
		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("connection success")
	fmt.Println("client address: ", conn.RemoteAddr())

	data_buffer := make([]byte, NetMessageMaxLength+1)

	NewMessage := MessageStruct{}

	for {
		recvLen, err := conn.Read(data_buffer)
		if err != nil {
			fmt.Println("Read error", err)
			return
		}

		used_num := NewMessage.AppendMsgData(data_buffer, recvLen)

	TestReadyProcess:
		// 如果ready,则处理
		if NewMessage.IsMsgReady() {

			// strBuffer := string(buffer[:recvLen])
			fmt.Println("Message: ", NewMessage.GetData())
			fmt.Println("Message len :", recvLen)
			// time.Sleep(time.Second * 1)                                                  //等一秒钟，可以看出client里面的read函数有阻塞效果
			// sendLen, err := conn.Write([]byte("I am server, you message :" + strBuffer)) //将client发过来的消息原样发送回去
			// if err != nil {
			// 	fmt.Println("send message error", err)
			// }

			fmt.Println("MessageID: ", NewMessage.header.MsgID)

			for i := int(0); i < int(NewMessage.header.param_num); i++ {

				this_param := ParamElement{}
				this_param.FromBytes(NewMessage.data_buffer[NetMessageHeaderSize+NetMessageElementSize*i : NetMessageHeaderSize+NetMessageElementSize*(i+1)])

				switch this_param.param_type {

				case BaseTypeBool:
					fmt.Println(NewMessage.GetBool(i))
					break
				case BaseTypeByte:
					fmt.Println(NewMessage.GetByte(i))
					break
				case BaseTypeInt8:
					fmt.Println(NewMessage.GetInt8(i))
					break
				case BaseTypeUint8:
					fmt.Println(NewMessage.GetUint8(i))
					break

				case BaseTypeInt16:
					fmt.Println(NewMessage.GetInt16(i))
					break
				case BaseTypeUint16:
					fmt.Println(NewMessage.GetUint16(i))
					break

				case BaseTypeInt32:
					fmt.Println(NewMessage.GetInt32(i))
					break
				case BaseTypeUint32:
					fmt.Println(NewMessage.GetUint32(i))
					break
				case BaseTypeFloat32:
					fmt.Println(NewMessage.GetFloat32(i))
					break

				case BaseTypeInt64:
					fmt.Println(NewMessage.GetInt64(i))
					break
				case BaseTypeUint64:
					fmt.Println(NewMessage.GetUint64(i))
					break
				case BaseTypeFloat64:
					fmt.Println(NewMessage.GetFloat64(i))
					break
				case BaseTypeComplex64:
					fmt.Println(NewMessage.GetComplex64(i))
					break

				case BaseTypeComplex128:
					fmt.Println(NewMessage.GetComplex128(i))
					break

				case BaseTypeString:
					fmt.Println(NewMessage.GetString(i))
					break
				case BaseTypeBinary:
					fmt.Println(NewMessage.GetBinary(i))
					break
				default:
					fmt.Println("UnKnow" + strconv.Itoa(int(this_param.param_type)))
					break
				}

			}

			fmt.Println("send message success")
			//fmt.Println("send message len；", sendLen)
		}

		if used_num < recvLen {
			// 说明前面已经拼接完成处理过了，则重新初始化新消息处理
			NewMessage.Init()
			recvLen = recvLen - used_num
			used_num = NewMessage.AppendMsgData(data_buffer[used_num:recvLen], recvLen)

			goto TestReadyProcess
		}

	}

}
