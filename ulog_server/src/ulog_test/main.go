package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"ushare/core"
)

var GLogServerAddr = flag.String("addr", "", "log server address")
var GLogServerPort = flag.Int("port", 13080, "log server port")

func query_ticket(serverName string, userName string) (string, error) {
	addr_query := fmt.Sprintf("http://%s:%v/query_ticket", *GLogServerAddr, *GLogServerPort)
	user_info := fmt.Sprintf("%s|%s|3|4", serverName, userName)
	encoded := base64.StdEncoding.EncodeToString([]byte(user_info))
	resp, err := http.PostForm(addr_query, url.Values{"user_info": {encoded}})
	if err != nil {
		core.LogErrorln("PostForm failed!", err)
		return "", err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		core.LogErrorln("response is not readable!", err)
		return "", err
	}

	log.Printf("query_ticket(): %s\n", resp.Status)
	return string(content), nil
}

func main() {
	flag.Parse()

	if *GLogServerAddr == "" {
		core.LogErrorln("Log Server not specified! (-addr)")
		os.Exit(-1)
	}

	ticket, err := query_ticket("测试服务器名", "测试账号名")
	if err != nil {
		core.LogErrorln("Query ticket failed!", err)
		os.Exit(-1)
	}

	log.Printf("ticket: %s\n", ticket)
}
