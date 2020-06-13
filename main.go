package main

import (
	view "go_spider/view"
)

const (
	url = "https://www.zhipin.com/c101280100-p100199/?ka=sel-city-101280100"
)


func main() {


	spider := view.NewPIGSpider(map[string]string{
		"User-Agent": "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)",
	})

	spider.SetRootUrl(url)

	spider.SetPrefix("https://www.zhipin.com")
	spider.Start()

}



