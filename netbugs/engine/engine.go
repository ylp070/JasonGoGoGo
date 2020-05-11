package engine

import (
	"fmt"

	"../fetcher"
)

type Request struct {
	Name     string
	Url      string
	ParseFun func(content []byte) ParseResult
}

type ParseResult struct {
	Req   []Request
	Items []interface{}
}

func NilParse(content []byte) ParseResult {
	return ParseResult{}
}

func Run(seeds ...Request) {

	var que []Request

	for _, seed := range seeds {
		que = append(que, seed)
	}

	for len(que) > 0 {
		cur := que[0]
		que = que[1:]

		//logs.Info("fetch url:", cur.Url)
		fmt.Printf("fetch url:", cur.Url)
		cont, e := fetcher.Fetch(cur.Url)
		if e != nil {
			//logs.Info("解析页面异常 url:", cur.Url)
			fmt.Printf("解析页面异常 url:", cur.Url)
			continue
		}

		resultParse := cur.ParseFun(cont)
		que = append(que, resultParse.Req...)

		for _, item := range resultParse.Items {
			fmt.Printf("内容项: %s \n", item)
		}
	}
}
