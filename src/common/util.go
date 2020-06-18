package common

import (
	"github.com/anaskhan96/soup"
	"log"
)

const IpProxy = "http://www.xicidaili.com/nn/"

func GetIPProxyList() []string {
	soup.Headers = map[string]string{"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/52.0.2743.116 Safari/537.36"}
	rsp, err := soup.Get(IpProxy)
	if err != nil {
		log.Printf("err get ip proxy : %v", err)
	}

	retIps := make([]string, 0)

	doc := soup.HTMLParse(rsp)

	ips := doc.FindAll("tr", "class", "odd")
	for index, _ := range ips {
		tds := ips[index].FindAll("td")
		ip, port := tds[1].Text(), tds[2].Text()
		retIps = append(retIps, "http://" + ip + ":" + port)
	}
	//for index, _ := range retIps {
	//	retIps[index]
	//}
	return retIps
}



