package SkUtils

import (
	"fmt"
	"os"
)

// 用于大量数据处理的优化，动态伸缩，数据安块存放，块大小可以初始化时定义，当一块数据被写满后动态生成新的块
// 应用情景：大文件在内存中写好后一次性写入磁盘，优化字符串拼接速度和文件写速度；大容量数据缓冲池
// 提供更大的缓冲区功能，可以实现大量

type SkBuffer struct {
	buffer_block_size int
	buffer_block_num  int
	buffer_array      map[int]([]byte)

	cur_block_index int
	cur_pos         int
}

// 外部调用1，初始化
func (buffer *SkBuffer) InitBuffer(buffer_block_size int) {

	buffer.buffer_block_size = buffer_block_size

	for _, data_pt := range buffer.buffer_array {

		data_pt = data_pt[0:0]

	}

	buffer.buffer_array = make(map[int]([]byte))

	// 只初始化第一个
	buffer.AddElement(1)

}

// 一般不要外部直接调用
func (buffer *SkBuffer) AddElement(block_num int) {

	for i := 0; i < block_num; i++ {
		buffer.buffer_array[buffer.buffer_block_num+i] = make([]byte, buffer.buffer_block_size)
	}

	buffer.buffer_block_num += block_num
}

// 外部调用2，添加数据
func (buffer *SkBuffer) Append(data_src []byte, data_size int) {

	// 能够放下数据
	left_buffer := buffer.buffer_block_size - buffer.cur_pos

	if left_buffer > data_size {

		cur_block, ok := buffer.buffer_array[buffer.cur_block_index]

		if ok {
			cur_block = append(cur_block[:buffer.cur_pos], data_src...)
		} else {
			fmt.Println("Append Error 1")
		}

		buffer.cur_pos += data_size
	} else {
		// 把剩余空着的填满
		cur_block, ok := buffer.buffer_array[buffer.cur_block_index]

		if ok {
			cur_block = append(cur_block[:buffer.cur_pos], data_src[:left_buffer]...)
		} else {
			fmt.Println("Append Error 2")
		}

		need_append_start := left_buffer

		// 只要一个block就可以存满
	Label_Append:
		if data_size-need_append_start < buffer.buffer_block_size {

			buffer.AddElement(1)
			buffer.cur_block_index = buffer.buffer_block_num - 1
			buffer.cur_pos = 0

			cur_block, ok := buffer.buffer_array[buffer.cur_block_index]

			if ok {
				cur_block = append(cur_block[:buffer.cur_pos], data_src[need_append_start:]...)
			} else {
				fmt.Println("Append Error 3")
			}

			buffer.cur_pos += data_size - need_append_start
			need_append_start = 0

		} else {

			buffer.AddElement(1)
			buffer.cur_block_index = buffer.buffer_block_num - 1
			buffer.cur_pos = 0

			cur_block, ok := buffer.buffer_array[buffer.cur_block_index]

			if ok {
				cur_block = append(cur_block[:buffer.cur_pos], data_src[need_append_start:need_append_start+buffer.buffer_block_size]...)
			} else {
				fmt.Println("Append Error 4")
			}

			buffer.cur_pos += buffer.buffer_block_size
			need_append_start = need_append_start + buffer.buffer_block_size

			goto Label_Append
		}
	}

}

// 外部调用3，写文件
func (buffer *SkBuffer) WriteToFile(fileName string) error {

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	for i := 0; i < buffer.buffer_block_num; i++ {

		block := buffer.buffer_array[i]

		if i == buffer.cur_block_index {
			_, err = f.Write(block[:buffer.cur_pos])
			break
		}
		_, err = f.Write(block)

	}

	if err1 := f.Close(); err == nil {
		err = err1
	}
	return err
}
