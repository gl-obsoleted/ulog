package main

import (
	"bytes"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"ushare"
	"ushare/core"
)

var GLogServerAddr = flag.String("addr", "", "log server address")
var GLogServerPort = flag.Int("port", 13080, "log server port")

var GTestLogDir = flag.String("testlogdir", `..\..\test_data\testlogs`, "the directory of test logs")

type LogFileInfo struct {
	Path   string
	Digest string
}

func build_url(verb string) string {
	return fmt.Sprintf("http://%s:%v/%s", *GLogServerAddr, *GLogServerPort, verb)
}

func build_url_t(verb string, ticket string) string {
	return build_url(verb) + "?ticket=" + ticket
}

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

func find_logs() []LogFileInfo {
	ret := []LogFileInfo{}
	filepath.Walk(*GTestLogDir, func(path string, info os.FileInfo, err error) error {
		if info.Mode().IsRegular() {
			var file LogFileInfo
			file.Path = path
			if digest, err := ushare.GetDataFileInfo(path); err == nil {
				file.Digest = digest
				ret = append(ret, file)
			}
		}
		return nil
	})
	return ret
}

func validate_logs(ticket string, logs []LogFileInfo) ([]string, error) {
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

func newfileUploadRequest(uri, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func upload_logs(ticket string, logs []string) error {
	for _, logfile := range logs {
		if len(logfile) == 0 {
			continue
		}

		log.Printf("uploading %s...\n", logfile)

		request, err := newfileUploadRequest(build_url_t("upload_resource", ticket), "file_list", logfile)
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

	// find local logs
	files := find_logs()
	log.Printf("files: \n")
	for i := range files {
		log.Println(files[i])
	}

	// validate logs
	validated, err := validate_logs(ticket, files)
	if err != nil {
		core.LogErrorln("validate_logs() failed!", err)
		os.Exit(-1)
	}

	if len(validated) > 0 {
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
}
