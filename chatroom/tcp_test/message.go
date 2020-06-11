package tcp_test

import (
	//"C"
	"fmt"
	"strconv"
	"unsafe"
)

// 支持的数据类型及常量
const (
	BaseTypeBool = iota
	BaseTypeByte // 8位数据
	BaseTypeInt8
	BaseTypeUint8

	BaseTypeInt16 // 16位数据
	BaseTypeUint16

	BaseTypeInt32 // 32位数据 rune\uint\int\uintptr这几个类型大小根据实际32bit或是64bit而定因此在网络传输中就不要再用了
	BaseTypeUint32
	BaseTypeFloat32

	BaseTypeInt64 // 64位数据
	BaseTypeUint64
	BaseTypeFloat64
	BaseTypeComplex64

	BaseTypeComplex128 // 128位数据

	BaseTypeString // 不定长数据 额外添加4位长度
	BaseTypeBinary
	BaseTypeMax

	NetMessageMaxLength   = 1024 * 10 // 单个消息最大字节长度
	NetMessageHeaderSize  = 5         // 消息头在消息串中所占的实际字节数，不是MessageHeader结构体内存字节数，结构体内存大小因为有对齐问题因此不同。如果要变更结构体需要重新取值
	NetMessageElementSize = 5         // 每个参数结构体在消息串中所占字节数，不是ParamElement结构体内存字节数，结构体内存大小因为有对齐问题因此不同。如果要变更结构体需要重新取值
)

// 网络消息头结构
// |------------------------|
// |    MessageHeader       |
// |------------------------|
// |    Data                |
// |------------------------|
type MessageHeader struct {
	MsgID uint16 // 消息类型

	// timestamp 时间戳
	// session  连接唯一标识Session之类
	param_num uint8 // 参数数量0~255

	data_length uint16 // 数据结构体长度
}

// // 消息头内存大小
// func (HD *MessageHeader) GetSize() int {
// 	return unsafe.Sizeof(MessageHeader)
// }

func (HD *MessageHeader) ToBytes() []byte {

	var data_bytes [5]byte

	value_ptr := (*uint16)(unsafe.Pointer(&data_bytes[0]))
	*value_ptr = HD.MsgID

	data_bytes[2] = HD.param_num

	value2_ptr := (*uint16)(unsafe.Pointer(&data_bytes[3]))
	*value2_ptr = HD.data_length

	return data_bytes[0:5]
}

func (HD *MessageHeader) FromBytes(data_bytes []byte) {
	HD.MsgID = *(*uint16)(unsafe.Pointer(&data_bytes[0]))
	HD.param_num = data_bytes[2]
	HD.data_length = *(*uint16)(unsafe.Pointer(&data_bytes[3]))
}

// 网络消息参数结构体
type ParamElement struct {
	param_type       byte
	data_start_index uint16
	data_size        uint16 // max data length 65535 byte 最大数据长度，以后如果不够可以改
}

// 获取内存大小
// func (PE *ParamElement) GetSize() int {
// 	return unsafe.Sizeof(ParamElement)
// }

func (PE *ParamElement) ToBytes() []byte {

	var data_bytes [5]byte

	data_bytes[0] = PE.param_type

	value_ptr := (*uint16)(unsafe.Pointer(&data_bytes[1]))
	*value_ptr = PE.data_start_index

	value2_ptr := (*uint16)(unsafe.Pointer(&data_bytes[3]))
	*value2_ptr = PE.data_size

	return data_bytes[0:5]
}

func (PE *ParamElement) FromBytes(data_bytes []byte) {
	PE.param_type = data_bytes[0]
	PE.data_start_index = *(*uint16)(unsafe.Pointer(&data_bytes[1]))
	PE.data_size = *(*uint16)(unsafe.Pointer(&data_bytes[3]))
}

// 完整消息结构
type MessageStruct struct {
	header      MessageHeader // 一个消息头结构体
	TempParam   ParamElement  // 用于临时参数解析的缓存
	CurDataSize uint16        // 当前有效数据长度，主要用在接收数据时记录当前实际接收到的数据长度
	data_buffer []byte        // 消息完整数据BUFFER，实际发送和接收的数据内容
}

//
func MinInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// 初始化重置结构体
func (MsgPtr *MessageStruct) Init() {
	MsgPtr.header.data_length = 0
	MsgPtr.CurDataSize = 0
}

// 组装消息，返回长度
func (MsgPtr *MessageStruct) MakeMessage(msgid uint16, params ...interface{}) int {

	index := 0
	num := 0

	param_num := len(params)

	// 移到实际数据开始
	index += NetMessageHeaderSize + NetMessageElementSize*param_num

	MsgPtr.header = MessageHeader{msgid, uint8(param_num), uint16(0)} //长度暂时不填

	ParamArray := make([]ParamElement, param_num)

	if MsgPtr.data_buffer == nil {
		MsgPtr.data_buffer = make([]byte, NetMessageMaxLength+1)
	}

	for _, arg := range params {
		switch arg.(type) {

		case bool:
			// 填写参数体
			ParamArray[num].param_type = byte(BaseTypeBool)
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(1)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			if arg.(bool) {
				MsgPtr.data_buffer[index] = byte(1)
			} else {
				MsgPtr.data_buffer[index] = byte(0)
			}

			index++ // 移到下个空位
			num++   // 参数数量+1
		case byte: // uint8
			// 填写参数体
			ParamArray[num].param_type = byte(BaseTypeByte)
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(1)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			MsgPtr.data_buffer[index] = byte(arg.(byte))

			index++ // 移到下个空位
			num++   // 参数数量+1

		case int8:
			// 填写参数体
			ParamArray[num].param_type = byte(BaseTypeInt8)
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(1)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			MsgPtr.data_buffer[index] = byte(arg.(int8))

			index++ // 移到下个空位
			num++   // 参数数量+1

		case int16:
			value := arg.(int16)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeInt16
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(2)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 2)
			bytes := (*[2]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:2]...)
			target_bytes := (*[2]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 2 // 移到下个空位
			num++      // 参数数量+1

		case uint16:
			value := arg.(uint16)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeUint16
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(2)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 2)
			bytes := (*[2]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:2]...)
			target_bytes := (*[2]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 2 // 移到下个空位
			num++      // 参数数量+1
		case int:
			value := int32(arg.(int))

			// 填写参数体
			ParamArray[num].param_type = BaseTypeInt32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&value), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], bytes[0:4]...)
			target_bytes := (*[4]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 4 // 移到下个空位
			num++      // 参数数量+1
		case int32:
			value := arg.(int32)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeInt32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&value), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:4]...)
			target_bytes := (*[4]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 4 // 移到下个空位
			num++      // 参数数量+1

		case uint32:
			value := arg.(uint32)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeUint32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:4]...)
			target_bytes := (*[4]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 4 // 移到下个空位
			num++      // 参数数量+1

		case float32:
			value := arg.(float32)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeFloat32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:4]...)
			target_bytes := (*[4]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 4 // 移到下个空位
			num++      // 参数数量+1

		case int64:
			value := arg.(int64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeInt64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:8]...)
			target_bytes := (*[8]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case uint64:
			value := arg.(uint64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeUint64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:8]...)
			target_bytes := (*[8]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case float64:
			value := arg.(float64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeFloat64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:8]...)
			target_bytes := (*[8]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case complex64:
			value := arg.(complex64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeComplex64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:8]...)
			target_bytes := (*[8]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case complex128:
			value := arg.(complex128)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeComplex128
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(16)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 16)
			bytes := (*[16]byte)(unsafe.Pointer(&value))
			//MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], (*bytes)[0:16]...)
			target_bytes := (*[16]byte)(unsafe.Pointer(&MsgPtr.data_buffer[index]))

			*target_bytes = *bytes

			index += 16 // 移到下个空位
			num++       // 参数数量+1

		case string:

			data_len := len(arg.(string))
			bytes := []byte(arg.(string)) //先把字符串转为byte切片
			// 填写参数体
			ParamArray[num].param_type = BaseTypeString
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(data_len)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), data_len)
			//bytes := (*[]byte)(unsafe.Pointer(&(arg.(string)[0])))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], bytes[0:data_len]...)

			index += data_len // 移到下个空位
			num++             // 参数数量+1

		case []byte:

			data_len := len(arg.([]byte))
			bytes := []byte(arg.(string)) //先把字符串转为byte切片

			// 填写参数体
			ParamArray[num].param_type = BaseTypeBinary
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(data_len)

			if int(ParamArray[num].data_size)+index >= NetMessageMaxLength {
				// 超出网络数据最大场度
				fmt.Println("message.go MakeMessage ERROR buffer is full, just %d parameters has be prepared, others all should be lost.", num)
				goto EndMakeMessage
			}

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), data_len)
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], bytes[0:data_len]...)

			index += data_len // 移到下个空位
			num++             // 参数数量+1
		default:
			fmt.Println("message.go ERROR MakeMessage Unknow data type =", arg)
			break
		}

	}

EndMakeMessage:
	// 写入消息头
	MsgPtr.header.data_length = uint16(index)
	copy(MsgPtr.data_buffer, MsgPtr.header.ToBytes())
	new_index := NetMessageHeaderSize
	// 写入参数结构
	// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[NetMessageHeaderSize]), unsafe.Pointer(&ParamArray[0]), NetMessageElementSize*param_num)
	for _, ref := range ParamArray {

		new_clip := MsgPtr.data_buffer[new_index : new_index+NetMessageElementSize]

		copy(new_clip, ref.ToBytes())

		new_index += NetMessageElementSize
	}

	//MsgPtr.data_buffer[index] = 0

	fmt.Println("message.go MakeMessage data size =", index)

	MsgPtr.CurDataSize = uint16(index)
	return index
}

// 组装好，或者解析以后返回实际数据
func (MsgPtr *MessageStruct) GetData() []byte {
	return MsgPtr.data_buffer[:MsgPtr.header.data_length]
}

// 返回原始data_buffer，注意长度不是实际数据
func (MsgPtr *MessageStruct) GetDataBuffer() []byte {

	if MsgPtr.data_buffer == nil {
		MsgPtr.data_buffer = make([]byte, NetMessageMaxLength+1)
	}

	return MsgPtr.data_buffer
}

func (MsgPtr *MessageStruct) IsMsgReady() bool {
	return MsgPtr.header.data_length > 0 && (MsgPtr.CurDataSize == MsgPtr.header.data_length)
}

// 解码data_buffer消息数据，返回参数数量
func (MsgPtr *MessageStruct) AppendMsgData(recived []byte, len int) int {

	if int(MsgPtr.CurDataSize)+len >= NetMessageMaxLength {
		// 超出数据范围不处理
		fmt.Println("message.go AppendMsgData data too big cant be processed cur=%d newaddlen=%d", int(MsgPtr.CurDataSize)+len, len)
		return 0
	}

	NewDataSize := int(MsgPtr.CurDataSize)

	need_append := len

	// 不超范围，则直接拼接上去
	if int(MsgPtr.CurDataSize)+len < NetMessageMaxLength {

		need_append = len
		MsgPtr.data_buffer = append(MsgPtr.data_buffer[:MsgPtr.CurDataSize], recived[0:len]...)

		NewDataSize = int(MsgPtr.CurDataSize) + len

	} else {
		// 超范围则拼接一部分
		need_append = NetMessageMaxLength - int(MsgPtr.CurDataSize)
		MsgPtr.data_buffer = append(MsgPtr.data_buffer[:int(MsgPtr.CurDataSize)], recived[0:NetMessageMaxLength-int(MsgPtr.CurDataSize)]...)

		NewDataSize = int(MsgPtr.CurDataSize) + len
	}

	// 头部数据是否已解出
	if MsgPtr.header.data_length < 1 {

		// 未解出，并且数据够了
		if NewDataSize > NetMessageHeaderSize {
			// 头部数据够了，可以直接先解出头
			MsgPtr.UnPackHeader()
		}
	}

	// 如果头部数据已解出，则看是否数据已填满
	if MsgPtr.header.data_length > 0 {

		if NewDataSize >= int(MsgPtr.header.data_length) {

			// 数据已完整

			// 计算使用掉的数据大小
			used := int(MsgPtr.header.data_length) - int(MsgPtr.CurDataSize)

			// 设置完整大小
			MsgPtr.CurDataSize = MsgPtr.header.data_length

			// 返回本次使用掉的数据
			return used

		} else {

			// 数据仍未完整，则加上拼接的部分，再等下次数据拼接
			MsgPtr.CurDataSize = uint16(int(MsgPtr.CurDataSize) + need_append)
		}
	} else {
		// 头部数据还不够，则加上拼接的部分，再等下次数据拼接
		MsgPtr.CurDataSize = uint16(int(MsgPtr.CurDataSize) + need_append)
	}

	return need_append
}

// 解码data_buffer消息数据，返回参数数量
func (MsgPtr *MessageStruct) UnPackHeader() int {

	// 解出文件头
	MsgPtr.header.FromBytes(MsgPtr.data_buffer[0:NetMessageHeaderSize])

	return int(MsgPtr.header.param_num)
}

// 获取参数数据,数据接收方需在UnPackMessage之后调用
func (MsgPtr *MessageStruct) GetNum() int {
	return int(MsgPtr.header.param_num)
}

func GetTypeString(t byte) string {

	switch t {
	case BaseTypeBool:
		return "BaseTypeBool"
		break
	case BaseTypeByte:
		return "BaseTypeByte"
		break
	case BaseTypeInt8:
		return "BaseTypeInt8"
		break
	case BaseTypeUint8:
		return "BaseTypeUint8"
		break

	case BaseTypeInt16:
		return "BaseTypeInt16"
		break
	case BaseTypeUint16:
		return "BaseTypeUint16"
		break

	case BaseTypeInt32:
		return "BaseTypeInt32"
		break
	case BaseTypeUint32:
		return "BaseTypeUint32"
		break
	case BaseTypeFloat32:
		return "BaseTypeFloat32"
		break

	case BaseTypeInt64:
		return "BaseTypeInt64"
		break
	case BaseTypeUint64:
		return "BaseTypeUint64"
		break
	case BaseTypeFloat64:
		return "BaseTypeFloat64"
		break
	case BaseTypeComplex64:
		return "BaseTypeComplex64"
		break

	case BaseTypeComplex128:
		return "BaseTypeComplex128"
		break

	case BaseTypeString:
		return "BaseTypeString"
		break
	case BaseTypeBinary:
		return "BaseTypeBinary"
		break
	default:
		return "UnKnow" + strconv.Itoa(int(t))
	}

	return "UnknowNULL"
}

// 获取参数开始----------------------------------------------------------------
// 获取bool参数
func (MsgPtr *MessageStruct) GetBool(index int) bool {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetBool index error index=%d max=%d", index, MsgPtr.header.param_num)
		return false
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeBool {

		fmt.Println("message.go GetBool type not [bool] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return false
	}

	return MsgPtr.data_buffer[MsgPtr.TempParam.data_start_index] == byte(1)
}

// 获取byte参数
func (MsgPtr *MessageStruct) GetByte(index int) byte {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetByte index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeByte {

		fmt.Println("message.go GetByte type not [byte] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	return MsgPtr.data_buffer[MsgPtr.TempParam.data_start_index]
}

// 获取Int8参数
func (MsgPtr *MessageStruct) GetInt8(index int) int8 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetInt8 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeInt8 {

		fmt.Println("message.go GetInt8 type not [int8] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	return int8(MsgPtr.data_buffer[MsgPtr.TempParam.data_start_index])
}

// 获取Uint8参数
func (MsgPtr *MessageStruct) GetUint8(index int) uint8 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetUint8 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeUint8 {

		fmt.Println("message.go GetUint8 type not [uint8] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	return uint8(MsgPtr.data_buffer[MsgPtr.TempParam.data_start_index])
}

// 获取int16参数
func (MsgPtr *MessageStruct) GetInt16(index int) int16 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go Getint16 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeInt16 {

		fmt.Println("message.go Getint16 type not [int16] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+2]

	valueptr := (*int16)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取uint16参数
func (MsgPtr *MessageStruct) GetUint16(index int) uint16 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetUint16 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeUint16 {

		fmt.Println("message.go GetUint16 type not [uint16] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+2]

	valueptr := (*uint16)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取int32参数
func (MsgPtr *MessageStruct) GetInt32(index int) int32 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetInt32 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeInt32 {

		fmt.Println("message.go GetInt32 type not [int32] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+4]

	valueptr := (*int32)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取float参数
func (MsgPtr *MessageStruct) GetInt(index int) int {
	return int(MsgPtr.GetInt32(index))
}

// 获取uint32参数
func (MsgPtr *MessageStruct) GetUint32(index int) uint32 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetUint32 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeUint32 {

		fmt.Println("message.go GetUint32 type not [uint32] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+4]

	valueptr := (*uint32)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取float参数
func (MsgPtr *MessageStruct) GetFloat(index int) float32 {
	return MsgPtr.GetFloat32(index)
}

// 获取float32参数
func (MsgPtr *MessageStruct) GetFloat32(index int) float32 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetFloat32 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeFloat32 {

		fmt.Println("message.go GetFloat32 type not [float32] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+4]

	valueptr := (*float32)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取int64参数
func (MsgPtr *MessageStruct) GetInt64(index int) int64 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetInt64 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeInt64 {

		fmt.Println("message.go GetInt64 type not [int64] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+8]

	valueptr := (*int64)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取uint64参数
func (MsgPtr *MessageStruct) GetUint64(index int) uint64 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetUint64 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeUint64 {

		fmt.Println("message.go GetUint64 type not [uint64] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+8]

	valueptr := (*uint64)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取float64参数
func (MsgPtr *MessageStruct) GetFloat64(index int) float64 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetFloat64 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeFloat64 {

		fmt.Println("message.go GetFloat64 type not [float64] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+8]

	valueptr := (*float64)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取Complex64参数
func (MsgPtr *MessageStruct) GetComplex64(index int) complex64 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetComplex64 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeComplex64 {

		fmt.Println("message.go GetComplex64 type not [complex64] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+8]

	valueptr := (*complex64)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取Complex128参数
func (MsgPtr *MessageStruct) GetComplex128(index int) complex128 {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetComplex64 index error index=%d max=%d", index, MsgPtr.header.param_num)
		return 0
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeComplex128 {

		fmt.Println("message.go GetComplex64 type not [complex64] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return 0
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+16]

	valueptr := (*complex128)(unsafe.Pointer(&bytes[0]))
	return *valueptr
}

// 获取String参数
func (MsgPtr *MessageStruct) GetString(index int) string {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetString index error index=%d max=%d", index, MsgPtr.header.param_num)
		return ""
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeString {

		fmt.Println("message.go GetString type not [string] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return ""
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	bytes := MsgPtr.data_buffer[start_index : start_index+int(MsgPtr.TempParam.data_size)]

	return string(bytes)
}

// 获取binary []byte参数
func (MsgPtr *MessageStruct) GetBinary(index int) []byte {

	if index >= int(MsgPtr.header.param_num) {
		fmt.Println("message.go GetBinary index error index=%d max=%d", index, MsgPtr.header.param_num)
		return nil
	}

	MsgPtr.TempParam.FromBytes(MsgPtr.data_buffer[NetMessageHeaderSize+NetMessageElementSize*index : NetMessageHeaderSize+NetMessageElementSize*(index+1)])
	if MsgPtr.TempParam.param_type != BaseTypeBinary {

		fmt.Println("message.go GetBinary type not [binary] real type is ", GetTypeString(MsgPtr.TempParam.param_type))
		return nil
	}

	start_index := int(MsgPtr.TempParam.data_start_index)

	return MsgPtr.data_buffer[start_index : start_index+int(MsgPtr.TempParam.data_size)]
}
