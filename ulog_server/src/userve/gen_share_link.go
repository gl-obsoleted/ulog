package userve

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
	"ushare"
	"ushare/core"
)

var SharingSessionSalt ushare.CryptoSalt

func Init_GenShareLink() error {
	salt, err := ushare.NewSalt()
	if err != nil {
		return core.NewStdErr(core.ERR_InitUServModuleFailed, "Init_GenShareLink() : generating salt failed : "+err.Error())
	}
	SharingSessionSalt = salt
	return nil
}

func GenShareLink(w http.ResponseWriter, r *http.Request) {
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

	var files []string
	var player_words string
	for tag, tag_val := range r.PostForm {
		if tag == "file" {
			for _, val := range tag_val {

				strArray := strings.Split(val, "|")
				if len(strArray) != 2 {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				filename := filepath.Base(strArray[0])
				sv_file_path, err := ushare.LocateDataFile(us, filename)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if !core.FileExists(sv_file_path) {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				files = append(files, filename+"|"+strArray[1])
			}
		}

		if tag == "player_words" {
			player_words = tag_val[0]
		}
	}

	sharingSessionKey, err := GenerateSharingDescriptorFile(us, files, player_words)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(sharingSessionKey)) // 直接把 share key 返回给客户端
}

func GenerateSharingDescriptorFile(us *ushare.UserSession, files []string, player_words string) (sharingSessionKey string, err error) {

	// 1. 汇总以下的信息
	//      内含玩家信息（服务器/账号/角色），此次分享的时间戳，此次分享包含的文件，此次分享包含的玩家的话
	// 2. 生成 session key
	//      使用 1 中的所有信息，加盐生成 session key
	var session ushare.DBSharingSession = ushare.DBSharingSession{}
	session.RegionName = us.RegionName
	session.ServerName = us.ServerName
	session.UserName = us.UserName
	session.RoleName = us.RoleName
	session.SharedTime = time.Now()
	session.SharedFiles = files
	session.SharedWords = player_words

	h := sha256.New()
	io.WriteString(h, session.RegionName)
	io.WriteString(h, session.ServerName)
	io.WriteString(h, session.UserName)
	io.WriteString(h, session.RoleName)
	io.WriteString(h, session.SharedTime.String())
	for _, f := range session.SharedFiles {
		io.WriteString(h, f)
	}
	io.WriteString(h, session.SharedWords)
	io.WriteString(h, string(SharingSessionSalt)) // apply salt, which varies from each server session
	session.SessionKey = fmt.Sprintf("%x", h.Sum(nil))

	if err := ushare.GDatabase.Insert(ushare.DBCName_SharingSession, &session); err != nil {
		return "", err
	}

	return session.SessionKey, nil
}
