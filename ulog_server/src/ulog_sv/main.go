package main

import (
	"net/http"
	"userve"
	"ushare"
	"ushare/core"
	"ushow"
)

func main() {
	err := ushare.InitEnv()
	defer ushare.DestroyEnv()
	if err != nil {
		core.LogFatalError("Application initializing failed!", err)
		return
	}

	if err := userve.InitTicketSys(); err != nil {
		core.LogFatalError("TicketSys initializing failed!", err)
		return
	}
	if err := userve.Init_GenShareLink(); err != nil {
		core.LogFatalError("'GenShareLink' initializing failed!", err)
		return
	}

	port, err := ushare.GConfig.LocateString("listen_port")
	if err != nil {
		core.LogFatalError("'listen_port' not set in config!", err)
		return
	}

	// 授权服务，用户提交内容相关，需要 ticket
	http.HandleFunc("/query_ticket", userve.QueryTicket)
	http.HandleFunc("/validate_files", userve.ValidateFiles)
	http.HandleFunc("/upload_resource", userve.UploadImage)
	http.HandleFunc("/gen_share_link", userve.GenShareLink)

	// ------ 分享服务的展示部分 ------
	// 静态文件服务
	http.Handle("/", http.FileServer(http.Dir("web")))
	http.Handle("/ugc/", http.StripPrefix("/ugc/", http.FileServer(http.Dir("ugc"))))
	// 请求一个 session 的相关信息
	http.HandleFunc("/query_session", ushow.QuerySession)

	http.ListenAndServe(":"+port, nil)
}
