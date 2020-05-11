package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"./parser"

	"./fetcher"
	SkUtils "./skutils"
)

func main() {

	PullDownZhenaiWebData()
}

// 读取本地文件
func PullFromLocalFile() {

	filedata, err := ioutil.ReadFile("test20200428.html")

	if err != nil {
		return
	}

	parser.GlobalUserData = make(map[string]parser.UserData)

	parser.ParseUserHtmlPage(filedata)

	fmt.Println("总共发现用户 :", len(parser.GlobalUserData))

	fmt.Println("开始ExtBufferSaveTabFile导出用户-----------------------")
	// 新的输出数据格式
	ExtBufferSaveTabFile("tabtable.txt")

	fmt.Println("开始ExtBufferSaveFile导出用户-----------------------")
	// 新的输出HTML
	ExtBufferSaveFile("extoutput.html")

	// fmt.Println("开始HtmlSaveFile导出用户-----------------------")
	// // 输出HTML
	// HtmlSaveFile("output.html")

	fmt.Println("总共导出用户 :", len(parser.GlobalUserData))

}

// 从网上直接爬取
func PullDownZhenaiWebData() {

	// 第一步, 通过url抓取页面
	url := "http://www.zhenai.com"

	//logs.Info("fetch url:", cur.Url)
	fmt.Println("fetch url:", url)

	content, e := fetcher.Fetch(url)

	if e != nil {
		fmt.Printf("解析页面异常 url:", url)
		return
	}

	// 清空原来的数据
	parser.GlobalUserData = make(map[string]parser.UserData)
	parser.NeedProcessList = make(map[string]string)
	parser.CityMap = make(map[string]string)

	drag_start_time := time.Now().Unix()

	parser.ParseCityList(content)

	drag_stop_time := time.Now().Unix()

	fmt.Println("爬取所有内容用时:" + strconv.Itoa(int(drag_stop_time-drag_start_time)))

	fmt.Println("总共发现用户 :", len(parser.GlobalUserData))

	fmt.Println("开始ExtBufferSaveTabFile导出用户-----------------------")
	// 新的输出数据格式
	ExtBufferSaveTabFile("tabtable.txt")

	fmt.Println("开始ExtBufferSaveFile导出用户-----------------------")
	// 新的输出HTML
	ExtBufferSaveFile("extoutput.html")

	fmt.Println("开始 输出城市列表文件 ----------------------")
	parser.SabeCityUrlList("city_url_list.txt")

	fmt.Println("总共导出用户 :", len(parser.GlobalUserData))

}

// 保存成数据表格
func ExtBufferSaveTabFile(fileName string) {

	mix_start_time := time.Now().Unix()

	file_content := []byte("ID	照片	姓名	性别	年龄	婚姻	工资	学历	区	城	省	国家	介绍	主页\r\n")

	const buffer_size int = 1024 * 1024 * 16 // 16MB
	var buffer SkUtils.SkBuffer

	buffer.InitBuffer(1024 * 1024 * 8)

	buffer.Append(file_content, len(file_content))

	i := 1
	for _, data := range parser.GlobalUserData {

		City1 := data.City
		City2 := ""
		City3 := ""
		City4 := ""

		City2In, ok := parser.CityMap[City1]

		if ok {
			City2 = City2In
			City3In, ok := parser.CityMap[City2]

			if ok {
				City3 = City3In
				City4In, ok := parser.CityMap[City3]

				if !ok {
					City4 = "中国"
					City3 = City2
					City2 = City1
				} else {
					City4 = City4In
				}

			} else {
				City4 = "中国"
				City3 = City2
				City2 = City1
			}

		} else {
			City2 = City1
			City3 = City1
			City4 = "中国"
		}

		element_content := data.ID + "	" + data.Image + "	" + data.Name + "	" + data.Gender + "	" +
			data.Age + "	" + data.Marry + "	" + data.Salary + "	" + data.Education + "	" + City1 + "	" + City2 + "	" + City3 + "	" + City4 + "	" + data.Introduction + "	" + data.SelfPage + "\r\n"

		file_content = []byte(element_content)
		buffer.Append(file_content, len(file_content))

		i++
	}

	mix_end_time := time.Now().Unix()

	fmt.Println("ExtBuffer合并TabTable文件内容用时:" + strconv.Itoa(int(mix_end_time-mix_start_time)))

	buffer.WriteToFile(fileName)

	savefile_time := time.Now().Unix()

	fmt.Println("ExtBuffer写入文件用时:" + strconv.Itoa(int(savefile_time-mix_end_time)))
}

// 保存成数据表格
func ExtBufferSaveFile(fileName string) {

	mix_start_time := time.Now().Unix()

	file_content := []byte("<html lang=\"zh-cn\"><head><meta charset=\"utf-8\"></head><body><table  border=\"1\">")

	const buffer_size int = 1024 * 1024 * 16 // 16MB
	var buffer SkUtils.SkBuffer

	buffer.InitBuffer(1024 * 1024 * 8)

	buffer.Append(file_content, len(file_content))

	i := 1
	for _, data := range parser.GlobalUserData {

		element_content := "<tr>"

		element_content += "<td style=\"width: 100px;\">" + strconv.Itoa(i+1) + "</td>"
		element_content += "<td><img src=\"" + data.Image + "\"></td>"
		element_content += "<td style=\"width: 200px;\"><a href=\"" + data.SelfPage + "\" target=\"_blank\">" + data.Name + "</a></td>"
		element_content += "<td>" + data.Gender + "</td>"
		element_content += "<td>" + data.Age + "</td>"
		element_content += "<td>" + data.City + "</td>"
		element_content += "<td>" + data.Marry + "</td>"
		element_content += "<td>" + data.Salary + "</td>"
		element_content += "<td>" + data.Education + "</td>"
		element_content += "<td>" + data.Introduction + "</td>"

		element_content += "</tr>\r\n"

		file_content = []byte(element_content)
		buffer.Append(file_content, len(file_content))

		i++
	}

	file_content = []byte("</table></body>")
	buffer.Append(file_content, len(file_content))

	mix_end_time := time.Now().Unix()

	fmt.Println("ExtBuffer合并HTML文件内容用时:" + strconv.Itoa(int(mix_end_time-mix_start_time)))

	buffer.WriteToFile(fileName)

	savefile_time := time.Now().Unix()

	fmt.Println("ExtBuffer写入文件用时:" + strconv.Itoa(int(savefile_time-mix_end_time)))
}

// 保存成HTML 非常慢不再需要这个
func HtmlSaveFile(fileName string) {

	mix_start_time := time.Now().Unix()

	file_content := "<html lang=\"zh-cn\"><head><meta charset=\"utf-8\"></head><body><table border=\"1\">"

	i := 1
	for _, data := range parser.GlobalUserData {

		file_content += "<tr>"

		file_content += "<td>" + strconv.Itoa(i+1) + "</td>"
		file_content += "<td><img src=\"" + data.Image + "\"></td>"
		file_content += "<td><a href=\"" + data.SelfPage + "\" target=\"_blank\">" + data.Name + "</a></td>"
		file_content += "<td>" + data.Gender + "</td>"
		file_content += "<td>" + data.Age + "</td>"
		file_content += "<td>" + data.City + "</td>"
		file_content += "<td>" + data.Marry + "</td>"
		file_content += "<td>" + data.Salary + "</td>"
		file_content += "<td>" + data.Introduction + "</td>"

		file_content += "</tr>\r\n"

		i++
	}

	file_content += "</table></body>"

	mix_end_time := time.Now().Unix()

	fmt.Println("合并HTML文件内容用时:" + strconv.Itoa(int(mix_end_time-mix_start_time)))

	if ioutil.WriteFile(fileName, []byte(file_content), 0644) == nil {
		fmt.Println("写入文件成功:" + fileName)
	}

	savefile_time := time.Now().Unix()

	fmt.Println("写入文件用时:" + strconv.Itoa(int(savefile_time-mix_end_time)))
}

func test() {
	content, err := fetcher.FetchUser("http://album.zhenai.com/u/1853811589")

	if err == nil {

		if ioutil.WriteFile("test.html", content, 0644) == nil {
			fmt.Println("写入文件成功:test.html - 1")
		}
	} else {
		fmt.Println("加载网址错误！-1")
	}

	content, err = fetcher.FetchUser("http://album.zhenai.com/u/1853811589")

	if err == nil {

		if ioutil.WriteFile("test.html", content, 0644) == nil {
			fmt.Println("写入文件成功:test.html -2")
		}
	} else {
		fmt.Println("加载网址错误！-2")
	}

	content, err = fetcher.FetchUser("http://album.zhenai.com/u/1853811589")

	if err == nil {

		if ioutil.WriteFile("test.html", content, 0644) == nil {
			fmt.Println("写入文件成功:test.html -3")
		}
	} else {
		fmt.Println("加载网址错误！-3")
	}
}
