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
	"net/http/httputil"
)

//Config Structure
type Config struct {
	Url string
	Authorization string
	Debug bool
}
//Config for Yelp API return Limit
const limit = 15;
const configFileName = "config.json";
const LogFileName = "log";

func main() {
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.HandleFunc("/get", getHandler)
	http.ListenAndServe(":8080", nil)
}

/**
 * Return Get Request string for yelp API 
 *
 * @param  req  *http.Request	http request pointer of current request
 * @param  query *url.Values	url values pointer from yelp request
 * @return rawQuery	string	query string
 */
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
	rawQuery = query.Encode()
	return
}

/**
 * Parse form data and get data from Yelp API
 *
 * @param  req  *http.Request	http request pointer of current request
 * @return resp	*http.Response	response pointer from Yelp API reponse
 */
func queryYelp(req *http.Request) (resp *http.Response) {
	config := getConfig()
	curl, err := http.NewRequest(http.MethodGet, config.Url, nil)
	if err != nil {
		panic(err)
	}
	query := curl.URL.Query()
	curl.URL.RawQuery = buildRawQuery(req,&query)
	curl.Header.Set("Authorization", config.Authorization)
	if(config.Debug){writeLog(curl.URL.String())}
	resp, err = http.DefaultClient.Do(curl)
	if(config.Debug){
		dump, _ := httputil.DumpResponse(resp, true)
		writeLog(string(dump))
	}
	if err != nil {
		panic(err)
	}
	return resp
}
/**
 * Handler function for /get url to get Data from YELP API
 * and return back to front end
 *
 * @param  req  http.ResponseWriter	Response writer that is used to send response to front end
 * @param  req  http.Request	http request pointer of incomming request
 */
func getHandler(resp http.ResponseWriter, req *http.Request) {
	config := getConfig()
	if(config.Debug){
		dump, _ := httputil.DumpRequest(req, true)
		writeLog(string(dump))
	}
	yelpResp := queryYelp(req)
	defer yelpResp.Body.Close()
	yelpData, err := ioutil.ReadAll(yelpResp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Fprint(resp, string(yelpData))
}
/**
 * Read config from json file and return Config Struct
 *
 * @return  config  Config	Config Struct
 */
func getConfig() (config Config) {
	file, err := os.Open(configFileName)
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
/**
 * Write to log file
 *
 * @param  v  ...interface	interfaces that needed to be write to log file
 */
func writeLog(v ...interface{}){
	f, err := os.OpenFile(LogFileName, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	log.Print(v)
		return
}
