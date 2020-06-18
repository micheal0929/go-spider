package main

import (
	"encoding/json"
	"fmt"
	_ "go_spider/src/common/version"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type workRsp struct {
	City string `json:"city"`
	Company string `json:"company"`
	Exp string `json:"exp"`
	Link string `json:"link"`
	JobName string `json:"job_name"`
	Salary string `json:"salary"`
}
type MsgRsp struct {
	Data []*workRsp `json:"data"`
}

type workReq struct {
	City      string `json:"city"`
	Company   string `json:"company"`
	Job       string `json:"job"`
	Desc      string `json:"desc"`
	SalaryBot int    `json:"salary_bot"`
	SalaryTop int    `json:"salary_top"`
}

type elkType struct {
	Query string `json:"query"`
}
type ColData struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type elkData struct {
	Columns []ColData `json:"columns"`
	Rows [][]interface{} `json:"rows"`
}
const (
	workSql = "SELECT city, company, exp, link, salary, job_name from job where %s"
)
var ColName []string
func main() {

	ColName = append(ColName, []string{"city", "company", "exp", "link", "salary"}...)
	http.HandleFunc("/work", WorkQuery)

	if err := http.ListenAndServe(":9091", nil); err != nil {
		fmt.Println("start err ", err)
	}
}

func WorkQuery(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("content-type", "application/json")
	byteData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("err :", err)
	}
	req := workReq{}
	err = json.Unmarshal(byteData, &req)
	if err != nil {
		fmt.Println("err : ", err)
	}
	fmt.Println(req)
	retStr := make([]string, 0)
	if req.City != "" {
		retStr = append(retStr, "city = '" + req.City+ "'")
	}
	if req.Company != "" {
		retStr = append(retStr, "company = '" + req.Company+ "'")
	}
	if req.Desc != "" {
		retStr = append(retStr, "desc = '" + req.Desc+ "'")
	}
	if req.Job != "" {
		retStr = append(retStr, "job_type = '" + req.Job+ "'")
	}
	if req.SalaryBot != 0 {
		retStr = append(retStr, "salary_bot >='" + strconv.Itoa(req.SalaryBot) + "'")
	}
	if req.SalaryTop != 0 {
		retStr = append(retStr, "salary_bot <= '" + strconv.Itoa(req.SalaryTop) + "'")
	}

	sql := fmt.Sprintf(workSql, strings.Join(retStr, " and "))
	fmt.Println(sql)
	secReq := elkType{}
	secReq.Query = sql
	data, _ := json.Marshal(secReq)
	resp, err := http.Post("http://127.0.0.1:9200/_xpack/sql", "application/json", strings.NewReader(string(data)))
	if err != nil {
		fmt.Println("err :", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	respData := elkData{}

	err = json.Unmarshal(body, &respData)
	if err != nil {
		fmt.Println("err :", err)
	}

	//for _, col := range respData.Columns {
	//	fmt.Println(col.Name, col.Type)
	//
	//}
	//
	rsp := MsgRsp{}

	for _, col := range respData.Rows {
		tmp := workRsp{}
		tmp.City = col[0].(string)
		tmp.Company = col[1].(string)
		tmp.Exp = col[2].(string)
		tmp.Link = col[3].(string)
		tmp.Salary = col[4].(string)
		tmp.JobName = col[5].(string)

		rsp.Data = append(rsp.Data, &tmp)
	}
	//rsp.Data = "hello"
	send, _ := json.Marshal(rsp)
	w.Write(send)
}
