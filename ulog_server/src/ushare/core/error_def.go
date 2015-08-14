package core

import (
	"fmt"
	"time"
)

// =========== Public Functions ===========

func NewStdErr(code int, info string) error {
	return &StdError{time.Now(), code, info}
}

// =========== Standard Error Codes ===========

const (
	ERR_EnvFatal = 100 << iota
	ERR_DBNotAvail
	ERR_NewUserFailed
	ERR_InitUServModuleFailed
)

const (
	ERR_ConfigFileLoadingFailed = 10000 << iota
	ERR_ConfigInvalidObject
	ERR_ConfigInvalidPath
	ERR_ConfigValueNotFound

	ERR_DataFileInvalid = 20000 << iota

	ERR_DatabaseInsertionFailed = 30000 << iota
)

// =========== Standard Error Code Descriptions ===========

var ErrorCodeLut map[int]string = map[int]string{
	ERR_EnvFatal:              "Env Fatal Error",
	ERR_DBNotAvail:            "Database not available",
	ERR_NewUserFailed:         "New User Failed.",
	ERR_InitUServModuleFailed: "Init UServ Module Failed.",

	ERR_ConfigFileLoadingFailed: "Config file loading failed.",
	ERR_ConfigInvalidObject:     "Config object is invalid.",
	ERR_ConfigInvalidPath:       "Invalid config accessing path.",
	ERR_ConfigValueNotFound:     "Config not found.",

	ERR_DataFileInvalid: "Data file invalid.",

	ERR_DatabaseInsertionFailed: "Database Insertion Failed",
}

// =========== Standard Error Structure ===========

type StdError struct {
	When time.Time
	Code int
	What string
}

func (e *StdError) Error() string {
	return fmt.Sprintf("\n  <%v> {%v}\n  [%v]\n  %s", e.Code, ErrorCodeLut[e.Code], e.When, e.What)
}
