/*
* Http (curl) request in golang
* @author Shashank Tiwari
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"log"
	"os"
	"encoding/json"
	"strconv"
)

type Config struct {
	Url string
	Authorization string
}
const limit = 15;

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}

func buildRawQuery(req *http.Request,query *url.Values) (rawQuery string) {
	req.ParseForm();
	keyword := req.FormValue("txtKeyword")
	location := req.FormValue("txtLocation")
	sort := req.FormValue("selSort")
	var page int64 = 1
	if s, err := strconv.ParseInt(req.FormValue("page"), 10, 32); err == nil {
		page = s
	}
	if len(location) ==0 {location="Singapore"}
	if len(keyword)>0 { query.Add("term", keyword) }
	query.Add("location", location)
	query.Add("sort_by", sort)
	query.Add("limit", strconv.FormatInt(limit, 10))
	if(page != 1){query.Add("offset", strconv.FormatInt((page*limit), 10))}
	//query.Add("latitude", req.FormValue("txtKeyword"))
	//query.Add("longitude", req.FormValue("txtKeyword"))
	rawQuery = query.Encode()
	return
}

func queryYelp(req *http.Request) (resp *http.Response) {
	config := getConfig()
	curl, err := http.NewRequest(http.MethodGet, config.Url, nil)
	if err != nil {
		panic(err)
	}
	query := curl.URL.Query()
	curl.URL.RawQuery = buildRawQuery(req,&query)
	curl.Header.Set("Authorization", config.Authorization)
	writeLog(curl.URL.String())
	resp, err = http.DefaultClient.Do(curl)
	if err != nil {
		panic(err)
	}
	return resp
}

func getHandler(resp http.ResponseWriter, req *http.Request) {
	yelpResp := queryYelp(req)
	defer yelpResp.Body.Close()
	yelpData, err := ioutil.ReadAll(yelpResp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(resp, string(yelpData))
}

func getConfig() (config Config) {
	file, err := os.Open("config.json")
	defer file.Close()
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		panic(err)
	}
	return config
}

func writeLog(v ...interface{}){
	f, err := os.OpenFile("log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Print(v)
		return
}
