package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "go_spider/src/common/version"
)

const (
	Url     = "https://www.zhipin.com/c101280100-p100199/?ka=sel-city-101280100"
	TestUrl = "http://12"
)

var Data []FileData

type FileData struct {
	Position string `json:"position"`
	Id       string `json:"id"`
}

func init() {
	if len(os.Args) < 2 {
		log.Println("need json path")
		os.Exit(1)
	}

	rootPath := os.Args[1]
	fileList := []string{"back.json", "mobile.json", "ai.json", "data.json", "support.json", "spider.json"}
	for i := 0; i < len(fileList); i++ {
		fileName := rootPath + "/" + fileList[i]
		data, err := ioutil.ReadFile(fileName)
		if err != nil {
			log.Println("err", err)
		}
		tmp := make([]FileData, 0)
		err = json.Unmarshal(data, &tmp)
		if err != nil {
			log.Println("err", err)
		}
		Data = append(Data, tmp...)

	}
}
func main() {
	str1 := "asSASA ddd dsjkdsjs dk"
	fmt.Println(&str1)
	// _, _ = &&str1
	//
	// spider := common.NewPIGSpider(map[string]string{
	//
	// 	"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.61 Safari/537.36",
	// })
	//cookis := http.Cookie{
	//	Path:       "=/",
	//	Domain:     "ynuf.aliapp.org",
	//	MaxAge:     31536000,
	//	Secure:     true,
	//	HttpOnly:   false,
	//	Raw:        "umdata_=G144FFA470A5CCBC3E41D6FB981C84C645C29A1",
	//	Unparsed:   nil,
	//}
	//
	//
	////web_service.SetPrefix("https://www.zhipin.com")
	//web_service.SetCookie(&cookis)
	// spider.Start()
	//

}

func ProxyTest() {
	proxyAddr := "http://157.255.144.77:80"

	httpUrl := "http://127.0.0.1:9091/spider"

	poststr := ""

	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		log.Fatal(err)
	}

	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}

	httpClient := http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}

	res, err := http.NewRequest("POST", httpUrl, strings.NewReader(poststr))
	if err != nil {
		log.Println(err)
		return
	}

	res.Header.Add("content-type", "application/x-ndjson")
	_, err = httpClient.Do(res)
	if err != nil {
		fmt.Println(err)
	}

	//
	//defer resp.Body.Close()
	//if resp.StatusCode != http.StatusOK {
	//	log.Println(err)
	//}
	//
	//c, _ := ioutil.ReadAll(resp.Body)
	//
	//fmt.Println(string(c))
}
