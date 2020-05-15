package tcp_test

import (
	//"C"
	"fmt"
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

	Maxparam_num          = 128  // 消息参数最大参数数量128个
	Maxdata_buffer        = 1024 // 单个消息最大消息长度1024个字节
	MessageHeaderDataSize = 5    // 消息头暂时占用12个字节
	ElementDataSize       = 5    // 每个参数结构体占5个字节
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

// 网络消息参数结构体
type ParamElement struct {
	param_type       byte
	data_start_index uint16
	data_size        uint16 // max data length 65535 byte 最大数据长度，以后如果不够可以改
}

// 完整消息结构
type MessageStruct struct {
	head_ptr    *MessageHeader  // 一个消息头结构体
	param_ptr   *[]ParamElement // 多个消息参数列表
	data_buffer []byte          // 消息完整数据BUFFER
}

//
func MinInt(a int, b int) int {
	if a < b {
		return a
	}

	return b
}

// 解码data部分

// data []byte

func (MsgPtr *MessageStruct) UnPackMessage() int {

	index := 0

	// 解出文件头
	MsgPtr.head_ptr = (*MessageHeader)(unsafe.Pointer(&MsgPtr.data_buffer[0]))

	// 移到参数数据开始
	index += MessageHeaderDataSize

	param_num := MsgPtr.head_ptr.param_num

	if param_num > Maxparam_num {
		param_num = uint8(Maxparam_num)
	}

	// 解出参数列表
	param_bytes := MsgPtr.data_buffer[index : index+int(param_num)*ElementDataSize]
	MsgPtr.param_ptr = (*[]ParamElement)(unsafe.Pointer(&param_bytes))

	// 移到实际数据数据头
	index += int(param_num) * ElementDataSize

	return int(param_num)

	// 尝试一下看看不是是不用再手动解析消息了
	// num := 0
	// InData := &MsgPtr.data_buffer

	// max_data_len := MsgPtr.header.data_length

	// for i := 0; index < max_data_len && i < Maxparam_num; {
	// 	// 取第一个字节，判断参数类型
	// 	switch param_type := data[index]{
	// 	case BaseType_bool:
	// 		index++
	// 		new_element := ParamElement{BaseType_bool, &data[index], 1}
	// 		index++
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_int8:
	// 		index++
	// 		new_element := ParamElement{BaseType_int8, &data[index], 1}
	// 		index++
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_uint8:
	// 		index++
	// 		new_element := ParamElement{BaseType_uint8, &data[index], 1}
	// 		index++
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_int16:
	// 		index++
	// 		new_element := ParamElement{BaseType_int16, &data[index], 2}
	// 		index += 2
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_uint16:
	// 		index++
	// 		new_element := ParamElement{BaseType_uint16, &data[index], 2}
	// 		index += 2
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_int32:
	// 		index++
	// 		new_element := ParamElement{BaseType_int32, &data[index], 4}
	// 		index += 4
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_uint32:
	// 		index++
	// 		new_element := ParamElement{BaseType_uint32, &data[index], 4}
	// 		index += 4
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_float32:
	// 		index++
	// 		new_element := ParamElement{BaseType_float32, &data[index], 4}
	// 		index += 4
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_int64:
	// 		index++
	// 		new_element := ParamElement{BaseType_int64, &data[index], 4}
	// 		index += 8
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_uint64:
	// 		index++
	// 		new_element := ParamElement{BaseType_uint64, &data[index], 4}
	// 		index += 4
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_float64:
	// 		index++
	// 		new_element := ParamElement{BaseType_float64, &data[index], 8}
	// 		index += 8
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_complex64:
	// 		index++
	// 		new_element := ParamElement{BaseType_complex64, &data[index], 8}
	// 		index += 8
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_complex128:
	// 		index++
	// 		new_element := ParamElement{BaseType_complex128, &data[index], 16}
	// 		index += 16
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_string:
	// 		index++
	// 		len_ptr unsafe.Pointer = &data[index]
	// 		data_size := *(*uint16(len_ptr))
	// 		index += 2
	// 		data_size := data[index]	// 取4字节大小
	// 		new_element := ParamElement{BaseType_int64, &data[index], data_size}
	// 		index += data_size
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	case BaseType_binary:
	// 		index++
	// 		len_ptr unsafe.Pointer = &data[index]   // 取出16位uint16长度
	// 		data_size := *(*uint16(len_ptr))
	// 		index += 2 // 移到数据位
	// 		new_element := ParamElement{BaseType_int64, &data[index], data_size}
	// 		index += data_size // 移到数据尾
	// 		num++
	// 		OutParams = append(OutParams, new_element)
	// 	default:
	// 		fmt.Println("ERROR UnPackMessage Unknow data type =", param_type,  " src_data=", InData)
	// 		break
	// 	}
	// }

	//return num
}

// 组装消息，返回长度
func (MsgPtr *MessageStruct) MakeMessage(msgid uint16, params ...interface{}) int {

	index := 0
	num := 0

	param_num := len(params)

	// 移到实际数据开始
	index += MessageHeaderDataSize + ElementDataSize*param_num

	header := MessageHeader{msgid, uint8(param_num), uint16(0)} //长度暂时不填

	ParamArray := make([]ParamElement, param_num)

	if MsgPtr.data_buffer == nil {
		MsgPtr.data_buffer = make([]byte, 1024*10)
	}

	for _, arg := range params {
		switch arg.(type) {

		case bool:
			// 填写参数体
			ParamArray[num].param_type = byte(BaseTypeBool)
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(1)

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

			// 实际数值写入buffer
			MsgPtr.data_buffer[index] = byte(arg.(byte))

			index++ // 移到下个空位
			num++   // 参数数量+1

		case int8:
			// 填写参数体
			ParamArray[num].param_type = byte(BaseTypeInt8)
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(1)

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

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 2)
			bytes := (*[2]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:2]...)

			index += 2 // 移到下个空位
			num++      // 参数数量+1

		case uint16:
			value := arg.(uint16)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeUint16
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(2)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 2)
			bytes := (*[2]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:2]...)

			index += 2 // 移到下个空位
			num++      // 参数数量+1
		case int:
			value := int32(arg.(int))

			// 填写参数体
			ParamArray[num].param_type = BaseTypeInt32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&value), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[:index], bytes[0:4]...)

			index += 4 // 移到下个空位
			num++      // 参数数量+1
		case int32:
			value := arg.(int32)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeInt32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&value), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:4]...)

			index += 4 // 移到下个空位
			num++      // 参数数量+1

		case uint32:
			value := arg.(uint32)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeUint32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:4]...)

			index += 4 // 移到下个空位
			num++      // 参数数量+1

		case float32:
			value := arg.(float32)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeFloat32
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(4)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 4)
			bytes := (*[4]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:4]...)

			index += 4 // 移到下个空位
			num++      // 参数数量+1

		case int64:
			value := arg.(int64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeInt64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:8]...)

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case uint64:
			value := arg.(uint64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeUint64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:8]...)

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case float64:
			value := arg.(float64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeFloat64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:8]...)

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case complex64:
			value := arg.(complex64)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeComplex64
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(8)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 8)
			bytes := (*[8]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:8]...)

			index += 8 // 移到下个空位
			num++      // 参数数量+1

		case complex128:
			value := arg.(complex128)
			// 填写参数体
			ParamArray[num].param_type = BaseTypeComplex128
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(16)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), 16)
			bytes := (*[16]byte)(unsafe.Pointer(&value))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], (*bytes)[0:16]...)

			index += 16 // 移到下个空位
			num++       // 参数数量+1

		case string:

			data_len := len(arg.(string))
			bytes := []byte(arg.(string)) //先把字符串转为byte切片
			// 填写参数体
			ParamArray[num].param_type = BaseTypeString
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(data_len)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), data_len)
			//bytes := (*[]byte)(unsafe.Pointer(&(arg.(string)[0])))
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], bytes[0:data_len]...)

			index += data_len // 移到下个空位
			num++             // 参数数量+1

		case []byte:

			data_len := len(arg.([]byte))
			bytes := []byte(arg.(string)) //先把字符串转为byte切片

			// 填写参数体
			ParamArray[num].param_type = BaseTypeBinary
			ParamArray[num].data_start_index = uint16(index)
			ParamArray[num].data_size = uint16(data_len)

			// 实际数值写入buffer
			// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[index]), unsafe.Pointer(&arg), data_len)
			MsgPtr.data_buffer = append(MsgPtr.data_buffer[0:index], bytes[0:data_len]...)

			index += data_len // 移到下个空位
			num++             // 参数数量+1
		default:
			fmt.Println("ERROR MakeMessage Unknow data type =", arg)
			break
		}

	}

	// 写入消息头
	header.data_length = uint16(index)
	// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[0]), unsafe.Pointer(&header), MessageHeaderDataSize)
	MsgPtr.head_ptr = (*MessageHeader)(unsafe.Pointer(&MsgPtr.data_buffer[0]))
	*(MsgPtr.head_ptr) = header
	// 写入参数结构
	// C.memcpy(unsafe.Pointer(&MsgPtr.data_buffer[MessageHeaderDataSize]), unsafe.Pointer(&ParamArray[0]), ElementDataSize*param_num)
	param_info_array := MsgPtr.data_buffer[MessageHeaderDataSize : MessageHeaderDataSize+ElementDataSize*param_num]

	MsgPtr.param_ptr = (*[]ParamElement)(unsafe.Pointer(&param_info_array))
	copy((*MsgPtr.param_ptr), ParamArray)

	//MsgPtr.data_buffer[index] = 0

	fmt.Println("MakeMessage data size =", index)

	return index
}

// 组装好，或者解析以后反回的实际数据
func (MsgPtr *MessageStruct) GetData() []byte {
	return MsgPtr.data_buffer[:MsgPtr.head_ptr.data_length]
}

func (MsgPtr *MessageStruct) GetDataBuffer() []byte {

	if MsgPtr.data_buffer == nil {
		MsgPtr.data_buffer = make([]byte, 1024*10)
	}

	return MsgPtr.data_buffer
}
