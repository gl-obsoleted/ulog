package ushare

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"ushare/core"
)

type Config struct {
	FilePath   string
	JsonObject map[string]interface{}
}

func (c *Config) IsValid() bool {
	return c.FilePath != "" && c.JsonObject != nil
}

func (c *Config) Load(filepath string) error {
	if c == nil {
		return core.NewStdErr(core.ERR_ConfigFileLoadingFailed, fmt.Sprintf("Trying to load config into 'nil' object. (%s)", filepath))
	}

	c.FilePath = ""
	c.JsonObject = nil

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return core.NewStdErr(core.ERR_ConfigFileLoadingFailed, fmt.Sprintf("Reading json ('%s') failed. (%s)", filepath, err.Error()))
	}

	var cfg interface{}
	err = json.Unmarshal([]byte(content), &cfg)
	if err != nil {
		return core.NewStdErr(core.ERR_ConfigFileLoadingFailed, fmt.Sprintf("Unmarshal() json ('%s') failed. (%s)", filepath, err.Error()))
	}

	var succ bool
	c.JsonObject, succ = cfg.(map[string]interface{})
	if !succ {
		return core.NewStdErr(core.ERR_ConfigFileLoadingFailed, fmt.Sprintf("Bad json object (type mismatch). (Current: %T, Expected: %T)", cfg, c.JsonObject))
	}

	c.FilePath = filepath
	return nil
}

func (c *Config) LocateString(path string) (string, error) {
	if !c.IsValid() {
		return "", core.NewStdErr(core.ERR_ConfigInvalidObject, fmt.Sprintf("Trying to access a 'nil' config object. (%s)", path))
	}

	var segments []string = strings.Split(path, ".")
	if len(segments) == 0 {
		return "", core.NewStdErr(core.ERR_ConfigInvalidPath, fmt.Sprintf("Trying to access the config object with an invalid path. (%s)", c.FilePath))
	}

	var type_conv_succ bool = false
	var ret_value string = ""

	json_obj := c.JsonObject
	for index, seg := range segments {
		if index != len(segments)-1 {
			json_obj, type_conv_succ = json_obj[seg].(map[string]interface{})
			if !type_conv_succ {
				return "", core.NewStdErr(core.ERR_ConfigValueNotFound, fmt.Sprintf("ConfigFile: %s, ConfigPath: %s", c.FilePath, path))
			}
		} else {
			ret_value, type_conv_succ = json_obj[seg].(string)
			if !type_conv_succ {
				return "", core.NewStdErr(core.ERR_ConfigValueNotFound, fmt.Sprintf("ConfigFile: %s, ConfigPath: %s", c.FilePath, path))
			}
			break
		}
	}
	return ret_value, nil
}
