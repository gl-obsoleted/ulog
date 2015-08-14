package ushare

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"ushare/core"
)

const GExecutableDirectory string = "bin"
const GConfigFileName string = "config.json"

var GConfig Config
var GDatabase Database

// ============ Public Functions ============

func InitEnv() error {
	log.Println("init environment: ...")
	if err := initWorkingDirectory(); err != nil {
		return err
	}

	log.Println("init logging: ...")
	if err := core.InitLogging(); err != nil {
		return err
	}

	log.Printf("init GConfig(%s)...", GConfigFileName)
	GConfig = Config{}
	if err := GConfig.Load(GConfigFileName); err != nil {
		return err
	}
	log.Println("GConfig initialized.")

	web_dir, _ := GConfig.LocateString("web_dir")
	if err := validateWebFolder(web_dir); err != nil {
		return err
	}
	log.Println("web directory validated: ", web_dir)

	db_addr, _ := GConfig.LocateString("db.addr")
	db_name, _ := GConfig.LocateString("db.name")
	log.Printf("connecting to db(addr:%s, name:%s)...", db_addr, db_name)
	if err := GDatabase.Connect(db_addr, db_name); err != nil {
		return err
	}
	log.Println("database connected.")

	log.Println("initialization Done.")
	return nil
}

func DestroyEnv() {
	log.Println("destroying environment...")
	GDatabase.Close()
	core.DestroyLogging()
}

func initWorkingDirectory() error {
	var wd string
	var err error
	if wd, err = os.Getwd(); err != nil {
		return core.NewStdErr(core.ERR_EnvFatal, err.Error())
	}

	log.Printf("checking working directory: %s...\n", wd)
	if filepath.Base(wd) != GExecutableDirectory {
		return core.NewStdErr(core.ERR_EnvFatal, fmt.Sprintf("The current running executable is not in the expected directory ('%s').", GExecutableDirectory))
	}

	var new_wd string = filepath.Dir(wd)
	if err = os.Chdir(new_wd); err != nil { // set the working dir to be the root
		return core.NewStdErr(core.ERR_EnvFatal, fmt.Sprintf("Changing working directory failed (from '%s', to '%s').", wd, new_wd))
	}
	log.Println("working directory changed: ", new_wd)

	return nil
}

func validateWebFolder(web string) error {
	web_dir, err := filepath.Abs(web)
	if err != nil {
		return core.NewStdErr(core.ERR_EnvFatal, err.Error())
	}

	if !core.DirExists(web_dir) {
		return core.NewStdErr(core.ERR_EnvFatal, fmt.Sprintf("web directory not found ('%s').", web_dir))
	}
	return nil
}
