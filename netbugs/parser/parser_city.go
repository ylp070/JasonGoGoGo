package parser

// 用于解析城市页面

import (
	"fmt"
	"io/ioutil"
	"regexp"
)

const cityPersionItem = `<div class="list-item">([^打]*)<div class="item-btn">打招呼</div></div>`

func ParseCity(content []byte, cityName string, regionName string) {

	fmt.Println("City Page Process Start")

	cityRegexp := regexp.MustCompile(cityPersionItem)

	if cityRegexp == nil {

		fmt.Println("ERROR parser.ParseCity 正则表达式错误！")
		return
	}

	subs := cityRegexp.FindAllSubmatch(content, -1)

	if len(subs) == 0 {

		if ioutil.WriteFile("test.html", content, 0644) == nil {
			fmt.Println("写入文件成功:test.html -2")
		}
	}

	for _, sub := range subs {

		//fmt.Println("City User　: " + string(sub[0]))

		ParseUserCityPage(sub[0], cityName, regionName)

	}

	fmt.Println("City Page Process End")
}
