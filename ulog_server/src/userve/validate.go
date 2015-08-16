package userve

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"ushare"
)

func ValidateFiles(w http.ResponseWriter, r *http.Request) {
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

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var neededFiles []string
	for file_path, file_info := range r.PostForm {
		filename := filepath.Base(file_path)
		if sv_file_path, err := ushare.LocateDataFile(us, filename); err == nil {
			log.Printf("file located: %v\n", sv_file_path)
			if sv_file_info, err := ushare.GetDataFileInfo(sv_file_path); err == nil {
				log.Printf("file fingerprint testing: %s, %s\n", sv_file_info, file_info[0])
				if sv_file_info == file_info[0] {
					log.Printf("file already exists on server: %v, %s\n", sv_file_path, sv_file_info)
					continue // this file has been uploaded to server correctly
				}
			} else {
				log.Printf("error: %v\n", err)
			}
		} else {
			log.Printf("error: %v\n", err)
		}

		neededFiles = append(neededFiles, file_path)
		//log.Printf("PostForm: %v, %v\n", file_path, file_info)
		// w.Write([]byte(file_path))
		// w.Write([]byte(file_info[0]))
	}

	if len(neededFiles) > 0 {
		w.Write([]byte(strings.Join(neededFiles, "|")))
	}
}
