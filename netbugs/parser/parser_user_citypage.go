package parser

// 用于解析用户页面

import (
	"regexp"
	"strings"
)

// 个人基本信息
const userPageAddress = `<a href="(http://album.zhenai.com[^"]+)" target`
const userPhoto = `<img src="([^"]+)" alt`
const userName = `<a href="http://album.zhenai.com/u/[^>]*" target="_blank">([^<]+)</a>`
const userID = `http://album.zhenai.com/u/([^>]*)`

// 性别，居住地，年龄，月薪，婚况，身 高
const userInfos = `</span>([^<]+)</td>`
const userInstroduction = `<div class="introduce">([^<]+)</div>`

type UserData struct {
	ID           string
	Name         string
	City         string
	Gender       string
	Age          string
	Marry        string
	LiveAddr     string
	Salary       string
	Height       string
	Introduction string
	SelfPage     string
	Image        string
	Education    string
}

var userIDRp *regexp.Regexp
var userPageRp *regexp.Regexp
var userPhotoRp *regexp.Regexp
var userNameRp *regexp.Regexp
var userInfosRp *regexp.Regexp
var userInstroductionRp *regexp.Regexp

var GlobalUserData map[string]UserData

func ParseUserCityPage(content []byte, cityName string, regionName string) {
	pro := UserData{}
	pro.City = cityName
	pro.LiveAddr = regionName

	if userIDRp == nil {
		userIDRp = regexp.MustCompile(userID)
	}

	if userPageRp == nil {
		userPageRp = regexp.MustCompile(userPageAddress)
	}
	if userPhotoRp == nil {
		userPhotoRp = regexp.MustCompile(userPhoto)
	}
	if userNameRp == nil {
		userNameRp = regexp.MustCompile(userName)
	}
	if userInfosRp == nil {
		userInfosRp = regexp.MustCompile(userInfos)
	}
	if userInstroductionRp == nil {
		userInstroductionRp = regexp.MustCompile(userInstroduction)
	}

	userPageMacth := userPageRp.FindSubmatch(content)
	userPhotoMacth := userPhotoRp.FindSubmatch(content)
	userNameMacth := userNameRp.FindSubmatch(content)
	userInfosMacth := userInfosRp.FindAllSubmatch(content, -1)
	userInstroductionMatch := userInstroductionRp.FindSubmatch(content)
	userIDMatch := userIDRp.FindSubmatch(userPageMacth[1])
	// for _, userInfo := range userInfosMacth {
	// 	text := string(userInfo[1])
	// 	fmt.Println(text)
	// }

	pro.ID = string(userIDMatch[1])
	pro.Name = string(userNameMacth[1])
	pro.SelfPage = string(userPageMacth[1])
	pro.Image = string(userPhotoMacth[1])
	pro.Introduction = string(userInstroductionMatch[1])
	pro.Introduction = strings.Replace(pro.Introduction, "\"", "“", -1)
	//性别，居住地，年龄，月薪，婚况，身 高
	pro.Gender = string(userInfosMacth[0][1])
	//pro.LiveAddr = string(userInfosMacth[1][1])
	pro.Age = string(userInfosMacth[2][1])

	if userInfosMacth[3][1][0] >= '0' && userInfosMacth[3][1][0] <= '9' {
		pro.Salary = string(userInfosMacth[3][1])
		pro.Education = "未填"
	} else {
		pro.Salary = string("未填")
		pro.Education = string(userInfosMacth[3][1])
	}

	pro.Marry = string(userInfosMacth[4][1])
	pro.Height = string(userInfosMacth[5][1])

	//fmt.Println(pro)
	GlobalUserData[pro.SelfPage] = pro
}
