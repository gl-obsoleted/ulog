package ushare

import (
	"fmt"
	"strconv"
	"strings"
	"ushare/core"
)

func LocateDataFile(us *UserSession, filename string) (string, error) {
	pos := strings.Index(filename, "_")
	if pos == -1 {
		return "", core.NewStdErr(core.ERR_DataFileInvalid, "cannot find date in filename: "+filename)
	}

	date_info := filename[:pos]
	file_path := fmt.Sprintf("ugc/%s/%s/%s/data/%s/%s", us.ServerName, us.UserName, us.RoleName, date_info, filename)
	return file_path, nil
}

func LocateDataFileByDetailedInfo(serverName string, userName string, roleName string, filename string) (string, error) {
	pos := strings.Index(filename, "_")
	if pos == -1 {
		return "", core.NewStdErr(core.ERR_DataFileInvalid, "cannot find date in filename: "+filename)
	}

	date_info := filename[:pos]
	file_path := fmt.Sprintf("ugc/%s/%s/%s/data/%s/%s", serverName, userName, roleName, date_info, filename)
	return file_path, nil
}

func GetDataFileInfo(dfp string) (string, error) {
	if !core.FileExists(dfp) {
		return "", core.NewStdErr(core.ERR_DataFileInvalid, fmt.Sprintf("datafile '%s' not found.", dfp))
	}

	fileSize := core.GetFileSize(dfp)
	if fileSize == 0 {
		return "", core.NewStdErr(core.ERR_DataFileInvalid, fmt.Sprintf("datafile '%s' size is 0.", dfp))
	}

	fileMD5Checksum := core.GetFileMD5(dfp)
	if fileMD5Checksum == "" {
		return "", core.NewStdErr(core.ERR_DataFileInvalid, fmt.Sprintf("datafile '%s' MD5 checksum is invalid.", dfp))
	}

	// log.Println("dfp: ", dfp)
	// log.Println("fileSize: ", strconv.FormatInt(fileSize, 10))
	// log.Println("fileMD5Checksum: ", fileMD5Checksum)
	return strconv.FormatInt(fileSize, 10) + "|" + fileMD5Checksum, nil
}
