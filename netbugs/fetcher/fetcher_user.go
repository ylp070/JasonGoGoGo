package fetcher

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// 抓取器
func FetchUser(url string) ([]byte, error) {

	// 第一步, 通过url抓取页面
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	// 方法1
	request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/66.0.3359.181 Safari/537.36")
	// 方法2 增加header选项
	// request.Header.Add("Cookie", "xxxxxx")
	// request.Header.Add("User-Agent", "xxx")
	// request.Header.Add("X-Requested-With", "xxxx")

	resp, err := client.Do(request)
	// resp, err := http.Get(url) // 会导致403

	if err != nil {
		return nil, fmt.Errorf("http get error :%s", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		if resp.StatusCode == http.StatusAccepted {

			// need retry
			fmt.Errorf("http get error 202 need retry :%s", url)
			fmt.Println("http get error 202 need retry :", url)

			test, _ := ioutil.ReadAll(resp.Body)

			if ioutil.WriteFile("test.html", test, 0644) == nil {
				fmt.Println("写入文件成功:test.html")
			}

			return nil, fmt.Errorf("http get error errCode:%d", resp.StatusCode)

		} else {

			fmt.Println("http get error errCode:%d", resp.StatusCode)

			return nil, fmt.Errorf("http get error errCode:%d", resp.StatusCode)

		}

	}

	// 读取出来body的所有内容
	return ioutil.ReadAll(resp.Body)
}
