package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"ushare/core"
)

func query_ticket(serverName string, userName string) (string, error) {
	user_info := fmt.Sprintf("%s|%s|3|4", serverName, userName)
	encoded := base64.StdEncoding.EncodeToString([]byte(user_info))
	resp, err := http.PostForm(build_url("query_ticket"), url.Values{"user_info": {encoded}})
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

func validate_logs(ticket string) ([]string, error) {
	// find local logs
	logs := find_logs()
	log.Printf("files: \n")
	for i := range logs {
		log.Println(logs[i])
	}

	url_values := url.Values{}
	for i := range logs {
		url_values.Set(logs[i].Path, logs[i].Digest)
	}

	resp, err := http.PostForm(build_url_t("validate_files", ticket), url_values)
	if err != nil {
		core.LogErrorln("PostForm failed!", err)
		return []string{}, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		core.LogErrorln("response is not readable!", err)
		return []string{}, err
	}

	log.Printf("validate_logs(): %s\n", resp.Status)
	if len(content) == 0 {
		return []string{}, nil
	}

	return strings.Split(string(content), "|"), nil
}

func upload_logs(ticket string, logs []string) error {
	for i, logfile := range logs {
		if len(logfile) == 0 {
			continue
		}

		log.Printf("uploading (%d/%d) %s...\n", i+1, len(logs), logfile)

		request, err := make_upload_request(build_url_t("upload_resource", ticket), "file_list", logfile)
		if err != nil {
			log.Fatal(err)
		}
		client := &http.Client{}
		resp, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		} else {
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			resp.Body.Close()
			fmt.Println(resp.StatusCode)
			fmt.Println(body)
		}
	}
	return nil
}

func main() {
	// argument parsing
	flag.Parse()
	if *GLogServerAddr == "" {
		core.LogErrorln("Log Server not specified! (-addr)")
		os.Exit(-1)
	}

	// get ticket
	ticket, err := query_ticket("测试服务器名", "测试账号名")
	if err != nil {
		core.LogErrorln("query_ticket() failed!", err)
		os.Exit(-1)
	}
	log.Printf("ticket: %s\n", ticket)

	// validate logs
	validated, err := validate_logs(ticket)
	if err != nil {
		core.LogErrorln("validate_logs() failed!", err)
		os.Exit(-1)
	}
	if len(validated) == 0 {
		log.Printf("No file needs to be uploaded.\n")
		os.Exit(0)
	}

	log.Printf("validated files (%d): \n", len(validated))
	for i := range validated {
		log.Println(validated[i])
	}

	// upload logs
	err = upload_logs(ticket, validated)
	if err != nil {
		core.LogErrorln("upload_logs() failed!", err)
		os.Exit(-1)
	}
}
