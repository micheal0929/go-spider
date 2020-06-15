package main

import (
	"fmt"
	"go_spider/view"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	Url = "https://www.zhipin.com/c101280100-p100199/?ka=sel-city-101280100"
	TestUrl = "http://12"
)


func main() {

	spider := view.NewPIGSpider(map[string]string{

		"User-Agent": "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	})

	spider.SetRootUrl(Url)

	spider.SetPrefix("https://www.zhipin.com")
	spider.Start()

}

func ProxyTest() {
	proxyAddr := "http://157.255.144.77:80"

	httpUrl := "http://127.0.0.1:9091/test"

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



