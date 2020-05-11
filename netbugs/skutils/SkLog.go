package SkUtils

import (
	"fmt"
	"log"
	"os"
)

// 将系统LOG重新向输出到文件
func RedirectLogToFile(filename string) {

	outfile, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666) //打开文件，若果文件不存在就创建一个同名文件并打开
	if err != nil {
		fmt.Println(*outfile, "open failed")
		os.Exit(1)
	}

	log.SetOutput(outfile)                               //设置log的输出文件，不设置log输出默认为stdout
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile) //设置答应日志每一行前的标志信息，这里设置了日期，打印时间，当前go文件的文件名

	//write log
	log.Printf("---------------Start Crawler----------------") //向日志文件打印日志，可以看到在你设置的输出文件中有输出内容了
}
