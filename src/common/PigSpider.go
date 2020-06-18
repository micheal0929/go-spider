package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/anaskhan96/soup"
	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

//大类
const (
	Server     = "100199" //后端开发
	MoveClient = "100299" //移动开发
	Test       = "100399" //测试
	Support    = "100499" //运维、技术支持
	PageClient = "100999" //前端开发
	AI         = "101399" //人工智能
)

//细分
const (
	//后端
	FullStack = "100123" //全栈
	GIS       = "100124" //GIS
	BackEnd   = "100199" //后端开发
	Java      = "100101" //Java
	Cplusplus = "100102" //C++
	Php       = "100103" //PHP
	C         = "100105" //C
	CSharp    = "100106" //C#
	NET       = "100107" //.Net
	Hadoop    =  "100108" //Hadoop
	Python    = "100109"//Python
	Delphi	="100110" //Delphi
	VB     = "100111" //VB
	Perl  = ""
)

const (
	JobUrl = "https://www.zhipin.com/job_detail/?query=&city=101280100&industry=&position=100508"
	Prefix = "https://www.zhipin.com"
)

type soupHeader map[string]string

type SalarySpider struct {
	pageQueue []string
	pageIndex int
	urlPrefix string
	urlBase   string
	//lastCity  string
	//lastJob   string
	cities    map[string]string
	esClient  *elasticsearch.Client
	counter   int64
}

type Job struct {
	JobType   string `json:"job_type"`   //岗位类型
	City      string `json:"city"`       //城市
	Company   string `json:"company"`    //公司名称
	JobName   string `json:"job_name"`   //职位名称
	Salary    string `json:"salary"`     //薪水范围
	SalaryBot int    `json:"salary_bot"` //最低
	SalaryTop int    `json:"salary_top"` //最高
	Desc      string `json:"desc"`       //职位描述
	Exp       string `json:"exp"`        //所需经验
	Link      string `json:"link"`       //职位连接
}


type FileData struct {
	Position string `json:"position"`
	Id string `json:"id"`
}

var TargetJobs map[string]string
var Data []FileData


func init() {
	TargetJobs = make(map[string]string)
	TargetJobs["后端开发"] = Server
	TargetJobs["移动开发"] = MoveClient
	//TargetJobs["测试"] = Test
	//TargetJobs["运维"] = Support
	//TargetJobs["前端开发"] = PageClient
	//TargetJobs["人工智能"] = AI
	if len(os.Args) < 2 {
		os.Exit(1)
	}

	rootPath := os.Args[1]
	fileList := []string{"back.json", "mobile.json", "ai.json", "data.json", "support.json", "spider.json"}
	for i :=0;i<len(fileList);i++ {
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
func NewPIGSpider(header soupHeader) *SalarySpider {
	ret := &SalarySpider{}
	ret.cities = make(map[string]string)
	config := elasticsearch.Config{
		Addresses: []string{"http://122.51.223.225:9200"},
	}
	ret.esClient, _ = elasticsearch.NewClient(config)
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
	index := strings.Index(pSp.urlBase, "city")
	if index == -1 {
		log.Fatalf("no city found!")
		return
	}
	pSp.urlBase = strings.Replace(pSp.urlBase, pSp.urlBase[index+5:index+14], city, 1)

	if len(pSp.pageQueue) == 0 {
		pSp.pageQueue = append(pSp.pageQueue, pSp.urlBase)
	}
}

func (pSp *SalarySpider) replaceJob(job string) {
	index := strings.Index(pSp.urlBase, "position")
	if index == -1 {
		log.Fatalf("no city found!")
		return
	}
	pSp.urlBase = strings.Replace(pSp.urlBase, pSp.urlBase[index+9:index+15], job, 1)
	//pSp.urlBase = strings.Replace(pSp.urlBase, pSp.lastJob, job, 1)
	if len(pSp.pageQueue) == 0 {
		pSp.pageQueue = append(pSp.pageQueue, pSp.urlBase)
	}
	fmt.Println(pSp.urlBase)
}

func (pSp *SalarySpider) reset() {
	pSp.pageIndex = 0
	pSp.pageQueue = make([]string, 0)
}
func (pSp *SalarySpider) parseAllCity() {
	//rsp, err := soup.Get(pSp.pageQueue[pSp.pageIndex])
	//if err != nil {
	//	log.Printf("get page : %s, err : %s", pSp.pageQueue[pSp.pageIndex], err)
	//}
	//doc := soup.HTMLParse(rsp)
	//aLinks := doc.Find("dd", "class", "city-wrapper").FindAll("a")
	//for i := 0; i < len(aLinks); i++ {
	//	cityStr := aLinks[i].Text()
	//	codeStr := aLinks[i].Attrs()["ka"]
	//	codes := strings.Split(codeStr, "-")
	//	if len(codes) > 2 && codes[2] != "100010000" {
	//		pSp.cities[codes[2]] = cityStr
	//	}
	//}
	pSp.cities["101280100"] = "广州"
	//for k, v := range pSp.cities {
	//	fmt.Println(k, v)
	//}
}
func (pSp *SalarySpider) parseOnePage(city, jobType string) bool {
	if pSp.pageIndex < len(pSp.pageQueue) {
		rsp, err := soup.Get(pSp.pageQueue[pSp.pageIndex])
		if err != nil {
			log.Printf("get page : %s, err : %s", pSp.pageQueue[pSp.pageIndex], err)
		}
		doc := soup.HTMLParse(rsp)
		allContents := doc.FindAll("div", "class", "info-primary")
		for i := 0; i < len(allContents); i++ {
			jobName := allContents[i].Find("div", "class", "job-title").Find("a")
			if strings.Index(jobName.Text(), "实习") != -1 {
				continue
			}
			companyName := allContents[i].Find("div", "class", "company-text").Find("a")
			salary := allContents[i].Find("span", "class", "red")
			exp := allContents[i].Find("p")
			endPos := strings.Index(salary.Text(), "K")
			if endPos == -1 {
				continue
			}
			array := strings.Split(salary.Text()[:endPos], "-")
			saBot, _ := strconv.Atoi(array[0])
			saTop, _ := strconv.Atoi(array[1])

			tags := doc.Find("div", "class", "tags").FindAll("span", "class", "tag-item")
			descStr := make([]string, len(tags))
			for i := 0; i < len(tags); i++ {
				descStr[i] = tags[i].Text()
			}
			job := Job{
				JobType:   jobType,
				City:      city,
				Company:   companyName.Text(),
				JobName:   jobName.Text(),
				Salary:    salary.Text(),
				Exp:       exp.Text(),
				SalaryBot: saBot,
				SalaryTop: saTop,
				Desc:      strings.Join(descStr, " "),
				Link:      pSp.urlPrefix + jobName.Attrs()["href"],
			}
			data, _ := json.Marshal(&job)
			fmt.Println(string(data))
			req := esapi.CreateRequest{
				Index:        "job",
				DocumentType: "job_salary",
				DocumentID:   strconv.Itoa(int(pSp.counter)),
				Body:         bytes.NewReader(data),
			}
			res, err := req.Do(context.Background(), pSp.esClient)
			if err != nil {
				log.Println(err)
			}
			atomic.AddInt64(&pSp.counter, 1)
			defer res.Body.Close()

			//pSp.esClient
			//fmt.Println(city, ",", companyName.Text(), ",", jobName.Text(), ",", salary.Text(), ",", exp.Text())
		}
		pageExist := doc.Find("div", "class", "page")
		if pageExist.Pointer == nil {
			pSp.pageIndex++
			return true
		}
		aLinks := pageExist.FindAll("a")
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

func (pSp *SalarySpider) SetCookie(cookie *http.Cookie) {

	soup.Cookie("Name", cookie.Name)
	soup.Cookie("Path", cookie.Path)
	soup.Cookie("Domain", cookie.Domain)
	soup.Cookie("Raw", cookie.Raw)
}
func (pSp *SalarySpider) Start() {
	pSp.SetRootUrl(JobUrl)
	pSp.SetPrefix(Prefix)
	//pSp.lastCity = "101280100"
	//pSp.lastJob = Server
	log.Println("start scrap!..")
	pSp.parseAllCity()


	for _, value := range Data {
		pSp.reset()
		pSp.replaceJob(value.Id)
		for key, val := range pSp.cities {
			pSp.reset()

			pSp.replaceCity(key)
			for {
				ret := pSp.parseOnePage(val, value.Position)
				if !ret {
					fmt.Printf("page : %s, count : %d\n", pSp.urlBase, pSp.pageIndex)
					break
				}
				<-time.After(time.Millisecond * 1000)
			}
		}
	}

	log.Println("scrap over!")

}
