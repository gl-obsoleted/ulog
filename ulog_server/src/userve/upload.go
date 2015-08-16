package userve

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"ushare"
	"ushare/core"
)

const MAX_MEM_CACHE = 1024 * 1024 // 1 MB

const (
	Upload_Done = 1 << iota
	Upload_ExistSkipped
)

func UploadImage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	ticket := r.URL.Query()["ticket"][0]
	if ticket == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("ticket extracted: %v\n", ticket)
	us := ushare.FindUser(ticket)
	if us == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the multipart form
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, "Parsing request failed: "+err.Error(), http.StatusBadRequest)
		return
	}
	m := r.MultipartForm

	// get files
	files := m.File["file_list"]
	if len(files) != 1 {
		err_info := fmt.Sprintf("File count: %v Only one file is allowed for uploading per session.", len(files))
		http.Error(w, err_info, http.StatusBadRequest)
		return
	}

	// check if it's ok for open
	file, err := files[0].Open()
	defer file.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	destPath, err := ushare.LocateDataFile(us, files[0].Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	destDir := filepath.Dir(destPath)
	if err := core.CreateDirIfNotExists(destDir); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//create destination file making sure the path is writeable.
	dst, err := os.Create(destPath)
	defer dst.Close()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//copy the uploaded file to the destination file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("File uploaded. (%s)\n", destPath)

	// 记录到数据库
	if err := ushare.DBLogNewUserEvent(ticket, ushare.EVT_ImageUploaded, ushare.EVT_OK, "file: "+dst.Name()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(string(Upload_Done)))
}
