package parser

// 用于解析用户页面

import (
	"regexp"
	"strconv"
	"strings"

	"../engine"
)

// 个人基本信息
const userRegexpSelfPage = `<div[^class]*class="m-btn purple"[^>]*>([^<]+)</div>`

// 个人隐私信息
const userPrivateRegexp = `<div [^>]*class="m-btn pink"[^>]*>([^<]+)</div>`

// 择偶条件
const userPartRegexp = `<div [^>]*class="m-btn"[^>]*>([^<]+)</div>`

type Profile struct {
	Name      string
	Marry     string
	Age       int
	Xingzuo   string
	Height    int
	Weight    int
	WorkAddr  string
	Salary    string
	Occuption string
	Education string
}

func ParseUser(content []byte, name string) engine.ParseResult {
	pro := Profile{}
	pro.Name = name
	// 获取用户的年龄
	userCompile := regexp.MustCompile(userRegexpSelfPage)
	usermatch := userCompile.FindAllSubmatch(content, -1)

	pr := engine.ParseResult{}
	for i, userInfo := range usermatch {
		text := string(userInfo[1])
		if i == 0 {
			pro.Marry = text
			continue
		}
		if strings.Contains(text, "岁") {
			age, _ := strconv.Atoi(strings.Split(text, "岁")[0])
			pro.Age = age
			continue
		}
		if strings.Contains(text, "座") {
			pro.Xingzuo = text
			continue
		}
		if strings.Contains(text, "cm") {
			height, _ := strconv.Atoi(strings.Split(text, "cm")[0])
			pro.Height = height
			continue
		}

		if strings.Contains(text, "kg") {
			weight, _ := strconv.Atoi(strings.Split(text, "kg")[0])
			pro.Weight = weight
			continue
		}

		if strings.Contains(text, "工作地:") {
			workaddr := strings.Split(text, "工作地:")[1]
			pro.WorkAddr = workaddr
			continue
		}

		if strings.Contains(text, "月收入:") {
			salary := strings.Split(text, "月收入:")[1]
			pro.Salary = salary
			continue
		}

		if i == 7 {
			pro.Occuption = text
			continue
		}

		if i == 8 {
			pro.Education = text
			continue
		}
	}
	pr.Items = append(pr.Items, pro)

	return pr
}
