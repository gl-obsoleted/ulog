package ushow

import (
	"encoding/json"
	"net/http"
)

func WriteJsonResponse(w http.ResponseWriter, object interface{}) {
	js, err := json.Marshal(object)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(js)
}
