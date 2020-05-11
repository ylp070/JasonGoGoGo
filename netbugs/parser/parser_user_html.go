package parser

// 用于解析用户页面

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// 个人基本信息
const userDataExp = `<tr>(.*?)<\/tr>`

// <td>2</td><td><img src="https://photo.zastatic.com/images/photo/421512/1686047733/4278633680083755.png?scrop=1&amp;crop=1&amp;w=140&amp;h=140&amp;cpos=north">
// </td><td><a href="http://album.zhenai.com/u/1686047733" target="_blank">等你爱我</a></td><td>女士</td><td>24</td><td>邹平</td><td>未婚</td><td>中专</td>
// <td>被伤过痛过哭过也笑过疯过，心地善良正直青春，就是害怕失去不珍惜我</td>
// <tr><td>46</td><td><img src="https://photo.zastatic.com/images/photo/408798/1635189633/9527391599808141.jpg?scrop=1&amp;crop=1&amp;w=140&amp;h=140&amp;cpos=north">
// </td><td><a href="http://album.zhenai.com/u/1635189633" target="_blank">觅爱</a></td><td>男士</td><td>33</td><td>淮安</td><td>离异</td><td>8001-12000元</td><td>我以前相信同甘共苦，现在我已经相信只同甘不共苦</td></tr><tr><td>47</td><td><img src="https://photo.zastatic.com/images/photo/256512/1026047557/19607008883668142.jpg?scrop=1&amp;crop=1&amp;w=140&amp;h=140&amp;cpos=north"></td><td><a href="http://album.zhenai.com/u/1026047557" target="_blank">我是风</a></td><td>男士</td><td>27</td><td>英德</td><td>未婚</td><td>3001-5000元</td><td>简 单 实 在 的 你 在 哪 里 ……</td></tr><tr><td>48</td><td><img src="https://photo.zastatic.com/images/photo/263850/1055398445/6411082152063069.jpg?scrop=1&amp;crop=1&amp;w=140&amp;h=140&amp;cpos=north"></td><td><a href="http://album.zhenai.com/u/1055398445" target="_blank">做你一生的女人</a></td><td>女士</td><td>49</td><td>武安</td><td>丧偶</td><td>高中及以下</td><td>找一个爱我的和我爱的人共度余生，不求大富大贵，只求两厢情愿，只要爱上对方就必须要包容对方的缺点，简简单单的度过后半生。</td></tr>
const userDetialExp = `<td>.*</td><td><img src="(https://photo.zastatic.com/images/photo/[^"]+)"></td><td><a href="(http://album.zhenai.com/u/[0-9]+)" target="_blank">([^<]+)</a></td><td>([^<]+)</td><td>([^<]+)</td><td>([^<]+)</td><td>([^<]+)</td><td>([^<]+)</td><td>([^<]+)</td>`
const userIDExp = `http://album.zhenai.com/u/([^>]*)`

var userIDExpRp *regexp.Regexp
var userDataRp *regexp.Regexp
var userDetailRp *regexp.Regexp

func ParseUserHtmlPage(content []byte) {

	pro := UserData{}

	if userIDExpRp == nil {
		userIDExpRp = regexp.MustCompile(userIDExp)
	}

	if userDataRp == nil {
		userDataRp = regexp.MustCompile(userDataExp)
	}

	if userDetailRp == nil {
		userDetailRp = regexp.MustCompile(userDetialExp)
	}

	userDataMacth := userDataRp.FindAllSubmatch(content, -1)

	for index, userdata := range userDataMacth {

		// if ioutil.WriteFile("test.txt", userdata[0], 0644) == nil {
		// 	fmt.Println("写入文件成功:" + "test.txt")
		// }
		fmt.Println("处理第:" + strconv.Itoa(index))

		userDetialMacth := userDetailRp.FindAllSubmatch(userdata[0], -1)

		if userDetialMacth == nil {
			// if ioutil.WriteFile("test.txt", userdata[0], 0644) == nil {
			// 	fmt.Println("写入文件成功:" + "test.txt")
			// }
			continue
		}

		pro.Image = string(userDetialMacth[0][1])
		pro.SelfPage = string(userDetialMacth[0][2])
		userIDMatchThis := userIDExpRp.FindSubmatch(userDetialMacth[0][2])
		pro.ID = string(userIDMatchThis[1])

		pro.Name = string(userDetialMacth[0][3])
		pro.Gender = string(userDetialMacth[0][4])
		pro.Age = string(userDetialMacth[0][5])

		pro.City = string(userDetialMacth[0][6])
		pro.Marry = string(userDetialMacth[0][7])
		pro.Education = string(userDetialMacth[0][8])

		if userDetialMacth[0][8][0] >= '0' && userDetialMacth[0][8][0] <= '9' {
			pro.Salary = string(userDetialMacth[0][8])
			pro.Education = "未填"
		} else {
			pro.Salary = string("未填")
			pro.Education = string(userDetialMacth[0][8])
		}

		pro.Introduction = string(userDetialMacth[0][9])
		pro.Introduction = strings.Replace(pro.Introduction, "\"", "“", -1)

		pro.LiveAddr = pro.City

		GlobalUserData[pro.SelfPage] = pro
	}

	//fmt.Println(pro)

}
