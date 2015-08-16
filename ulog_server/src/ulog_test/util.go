package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"ushare"
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

func make_upload_request(uri, paramName, path string) (*http.Request, error) {
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
