package userve

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"strings"
	"ushare"
	"ushare/core"
)

const GTicketSaltLength = 128

var GTicketSalt []byte

func InitTicketSys() error {
	saltBuf := make([]byte, GTicketSaltLength)
	if _, err := rand.Read(saltBuf); err != nil {
		return core.NewStdErr(core.ERR_InitUServModuleFailed, "InitTicketSys() : generating salt failed : "+err.Error())
	}

	GTicketSalt = saltBuf

	log.Println("Ticket system initialized.")
	log.Printf("GTicketSalt: %v\n", fmt.Sprintf("%x", GTicketSalt))
	return nil
}

func QueryTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	base := r.PostFormValue("user_info")
	if base == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoded, err := base64.StdEncoding.DecodeString(base)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	strArray := strings.Split(string(decoded), "|")
	if len(strArray) != 4 {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ticket, err := ushare.NewUser(strArray, GTicketSalt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// -- 流程调试 --
	// log.Println("base: ", base)
	// log.Println("decoded: ", string(decoded))
	// log.Println("ticket: ", ticket)

	// 记录到数据库
	if err := ushare.DBLogNewUserEvent(ticket, ushare.EVT_TicketAcquired, ushare.EVT_OK, "info: "+string(decoded)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// http response (目前暂时不使用 json)
	w.Write([]byte(ticket))
	//WriteJsonResponse(w, us.Ticket)

	log.Println("Ticket allocated: ", ticket)
}
