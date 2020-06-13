package view

import (
	"fmt"
	"github.com/anaskhan96/soup"
	"log"
	"strconv"
	"strings"
	"time"
)

type soupHeader map[string]string

type SalarySpider struct {
	pageQueue []string
	pageIndex int
	urlPrefix string
	urlBase   string
	lastCity  string
	cities    map[string]string
}

func NewPIGSpider(header soupHeader) *SalarySpider {
	ret := &SalarySpider{}
	ret.cities = make(map[string]string)
	soup.Headers = header
	return ret
}
func (pSp *SalarySpider) addIndexUrl(url string) {
	pSp.pageQueue = append(pSp.pageQueue, url)
}

func (pSp *SalarySpider) SetRootUrl(url string) {
	pSp.urlBase = url
	pSp.addIndexUrl(url)
}
func (pSp *SalarySpider) SetPrefix(url string) {
	pSp.urlPrefix = url
}

func (pSp *SalarySpider) replaceCity(city string) {
	strings.ReplaceAll(pSp.urlBase, pSp.lastCity, city)
	if len(pSp.pageQueue) == 0 {
		pSp.pageQueue = append(pSp.pageQueue, pSp.urlBase)
	}
}

func (pSp *SalarySpider) reset() {
	pSp.pageIndex = 0
	pSp.pageQueue = pSp.pageQueue[:]
}
func (pSp *SalarySpider) parseAllCity() {
	rsp, err := soup.Get(pSp.pageQueue[pSp.pageIndex])
	if err != nil {
		log.Printf("get page : %s, err : %s", pSp.pageQueue[pSp.pageIndex], err)
	}
	doc := soup.HTMLParse(rsp)
	aLinks := doc.Find("dd", "class", "city-wrapper").FindAll("a")
	for i := 0; i < len(aLinks); i++ {
		cityStr := aLinks[i].Text()
		codeStr := aLinks[i].Attrs()["ka"]
		codes := strings.Split(codeStr, "-")
		if len(codes) > 2 && codes[2] != "100010000" {
			pSp.cities[codes[2]] = cityStr
		}
	}
	//for k, v := range pSp.cities {
	//	fmt.Println(k, v)
	//}
}
func (pSp *SalarySpider) parseOnePage() bool {
	if pSp.pageIndex < len(pSp.pageQueue) {
		fmt.Printf("page %d \n", pSp.pageIndex+1)
		rsp, err := soup.Get(pSp.pageQueue[pSp.pageIndex])
		if err != nil {
			log.Printf("get page : %s, err : %s", pSp.pageQueue[pSp.pageIndex], err)
		}
		doc := soup.HTMLParse(rsp)
		allContents := doc.FindAll("div", "class", "info-primary")
		for i := 0; i < len(allContents); i++ {
			jobName := allContents[i].Find("div", "class", "job-title").Find("a")
			companyName := allContents[i].Find("div", "class", "company-text").Find("a")
			salary := allContents[i].Find("span", "class", "red")
			exp := allContents[i].Find("p")
			fmt.Println(companyName.Text(), "-", jobName.Text(), ":", salary.Text(), ":", exp.Text())
		}
		aLinks := doc.Find("div", "class", "page").FindAll("a")
		maxIndex := 0
		for i := 0; i < len(aLinks); i++ {
			index, _ := strconv.Atoi(aLinks[i].Text())
			if maxIndex < index {
				maxIndex = index
			}
		}
		nextLink := doc.Find("div", "class", "page").FindAll("a", "class", "next")
		if len(nextLink) >= 1 && maxIndex > pSp.pageIndex+1 {
			pSp.pageQueue = append(pSp.pageQueue, pSp.urlPrefix+nextLink[0].Attrs()["href"])
		}
		pSp.pageIndex++
		return true
	}
	return false

}

func (pSp *SalarySpider) Start() {
	log.Println("start scrap!..")
	pSp.parseAllCity()
	for {
		for key, value := range pSp.cities {
			pSp.replaceCity(key)
			fmt.Printf("city : %s\n", value)
			for {
				ret := pSp.parseOnePage()
				if !ret {
					break
				}
				time.After(time.Microsecond * 100)
			}
			pSp.reset()
			pSp.lastCity = key
		}

	}
	log.Println("scrap over!")

}
