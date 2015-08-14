package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"ushare/core"
)

func test_get() {
	resp, err := http.Get("http://localhost:13080/")
	if err != nil {
		core.LogFatalError("Server not running!", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		core.LogFatalError("response is not readable!", err)
		return
	}

	log.Printf("  %s\n", content)
}

func test_post() {
	var user_info string
	user_info = "2|8|3|4"
	encoded := base64.StdEncoding.EncodeToString([]byte(user_info))

	resp, err := http.PostForm("http://localhost:13080/query_ticket", url.Values{"user_info": {encoded}})
	if err != nil {
		core.LogFatalError("PostForm failed!", err)
		return
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		core.LogFatalError("response is not readable!", err)
		return
	}

	log.Printf("  %s\n", resp.Status)
	log.Printf("  %s\n", content)
}

func main() {
	test_post()
}
