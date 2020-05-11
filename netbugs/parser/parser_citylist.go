package parser

// 用于解析城市列表页面

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"../fetcher"
	SkUtils "../skutils"
)

const cityListRegexp = `<a[^>]*href="(http://www.zhenai.com/zhenghun/[a-z1-9]+)"[^>]*>([^<]+)</a>`

const citysAreaRegexp = `area(.*)</div>`

//const cityRegionRegexp = `<a[^>]*href="(http://www.zhenai.com/zhenghun/[a-z1-9]+)"[^>]*black[^>]*>([^<]+)</a>`
const cityRegionRegexp = `<a[^>]*href="(http://www.zhenai.com/zhenghun/[a-z1-9]+)"[^>]*>([^<]+)</a>`

// 名称到URL
var NeedProcessList map[string]string

var city_re *regexp.Regexp
var area_re *regexp.Regexp
var region_re *regexp.Regexp

// map self parent
var CityMap map[string]string

func ParseCityList(content []byte) {

	city_re = regexp.MustCompile(cityListRegexp)
	area_re = regexp.MustCompile(citysAreaRegexp)

	region_re = regexp.MustCompile(cityRegionRegexp)

	allCityInMainPage := city_re.FindAllSubmatch(content, -1)

	fmt.Println("城市列表：")

	for _, city := range allCityInMainPage {

		fmt.Println("City: " + string(city[2]))
	}

	fmt.Println("开始处理城市列表----------------------")

	citylist_start_time := time.Now().Unix()

	var content_new []byte
	var e error
	var url string

	// 当前主页上的所有城市
	for _, city := range allCityInMainPage {

		// 通过主城页分析区页
		url = string(city[1])

		// 递归查找所有可以处理的分页
		ProcessCityPage(url, string(city[2]))
	}

	// 先输出城市列表文件
	SabeCityUrlList("city_url_list_inner.txt")

	citylist_end_time := time.Now().Unix()
	//
	fmt.Println("递归查找到 " + strconv.Itoa(len(NeedProcessList)) + " 个城市 用时：" + strconv.Itoa(int(citylist_end_time-citylist_start_time)) + "--------------------")

	city_index := 1

	for name, base_url := range NeedProcessList {

		fmt.Println(strconv.Itoa(city_index) + " City: " + name + " process start=================================")

		for index := 1; index < 100; index++ {

			url = base_url + "/" + strconv.Itoa(index)

			fmt.Println("Start Fetch City Page: " + url)

			content_new, e = fetcher.Fetch(url)

			if e != nil || strings.Contains(string(content_new), "404 您访问的页面不存在") {
				// 城市已达最后一页
				fmt.Println("City: " + name + " Region: " + name + " 总共有：" + strconv.Itoa(index-1) + "页")
				break
			}

			ParseCity(content_new, name, name)

		}

		fmt.Println(strconv.Itoa(city_index) + " City: " + name + " process end-------------------")

		city_index++
	}

}

func ProcessCityPage(url string, selfname string) {

	content_new, e := fetcher.Fetch(url)

	if e != nil || strings.Contains(string(content_new), "404 您访问的页面不存在") {
		// 城市已达最后一页
		fmt.Println("City: " + selfname + " 页面不存在:" + url)
		return
	}

	RegionsArea := area_re.FindSubmatch(content_new)

	// if ioutil.WriteFile("test.html", RegionsArea[0], 0644) == nil {
	// 	fmt.Println("写入文件成功:test.html - 1")
	// }

	if RegionsArea != nil {
		//
		AllCityRegions := region_re.FindAllSubmatch(RegionsArea[1], -1)

		if AllCityRegions != nil {

			// 判断最后一个是不是就是当前处理名称
			this_page_parent_city := string(AllCityRegions[len(AllCityRegions)-1][2])

			if this_page_parent_city == selfname {

				// 是的，说明是根节点,需要遍历所有子节点，根节点自身无需处理
				NeedProcessList[selfname] = ""

				CityMap[selfname] = "中国"

				for _, region := range AllCityRegions {

					_, ok := NeedProcessList[string(region[2])]

					if !ok {

						fmt.Println("Find a new City1 : No." + strconv.Itoa(len(CityMap)) + "  " + string(region[2]))

						// 未处理过的，开始进行处理
						// 把自已和父关联
						CityMap[string(region[2])] = this_page_parent_city

						// 通过主城页分析区页
						sub_url := string(region[1])

						ProcessCityPage(sub_url, string(region[2]))

					} else {

						// 已处理过则不再处理
						//fmt.Println("Parent: " + string(selfname) + " Region: " + string(region[2]) + " Is Used Ignore!")

					}
				}
			} else {
				// 当前处理页不是最后一页，则处理自已
				NeedProcessList[selfname] = url

				// 并且仍然要处理其他分页
				for _, region := range AllCityRegions {

					_, ok := NeedProcessList[string(region[2])]

					if !ok {

						fmt.Println("Find a new City2 : No." + strconv.Itoa(len(CityMap)) + "  " + string(region[2]))

						// 把自已和父关联
						CityMap[string(region[2])] = this_page_parent_city

						sub_url := string(region[1])
						ProcessCityPage(sub_url, string(region[2]))

					} else {

						// 已处理过则不再处理
						// fmt.Println("Parent: " + string(selfname) + " Region: " + string(region[2]) + " Is Used Ignore!")

					}
				}
			}

			// 没有找到分页，则处理自已
		} else {

			// 当前处理页不是最后一页，则处理自已
			NeedProcessList[selfname] = url

			fmt.Println("Find a new City3 :" + selfname)
		}
	} else {
		// region为空，则处理自已
		// 当前处理页不是最后一页，则处理自已
		NeedProcessList[selfname] = url
		fmt.Println("Find a new City4 :" + selfname)
	}

}

func SabeCityUrlList(fileName string) {

	fmt.Println("将城市页面列表保存为citylist.txt")

	citylist_start_time := time.Now().Unix()

	const buffer_size int = 1024 * 1024 * 1 // 16MB
	var buffer SkUtils.SkBuffer

	buffer.InitBuffer(1024 * 1024 * 8)

	i := 1
	for city, url := range NeedProcessList {

		element_content := city + "	" + url + "\r\n"

		file_content := []byte(element_content)
		buffer.Append(file_content, len(file_content))

		i++
	}

	citylist_end_time := time.Now().Unix()

	fmt.Println("SabeCityUrlList 合并文件内容用时:" + strconv.Itoa(int(citylist_end_time-citylist_start_time)))

	buffer.WriteToFile(fileName)

	savefile_time := time.Now().Unix()

	fmt.Println("SabeCityUrlList 写入文件用时:" + strconv.Itoa(int(savefile_time-citylist_end_time)))
}
