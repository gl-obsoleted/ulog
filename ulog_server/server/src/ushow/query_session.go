package ushow

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"time"
	"ushare"
)

type ShareSessionLut map[string]ushare.DBSharingSession

var GCachedSharedSessions ShareSessionLut = ShareSessionLut{}

type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	raw := time.Time(t)
	stamp := fmt.Sprintf("\"%v年%v月%v日 %v:%v:%v\"", raw.Year(), int(raw.Month()), raw.Day(), raw.Hour(), raw.Minute(), raw.Second())
	return []byte(stamp), nil
}

type SharedSessionWithImagePaths struct {
	ServerName string
	RoleName   string

	SharedTime  JSONTime
	SharedWords string

	ImagePaths []string
}

func QuerySession(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	share_id := r.URL.Query()["share_id"][0]
	if share_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("share_id: ", share_id)

	var session ushare.DBSharingSession
	if s, ok := GCachedSharedSessions[share_id]; ok {
		session = s
	} else {
		s := ushare.DBSharingSession{}
		sharing_session_collection := ushare.GDatabase.Database.C(ushare.DBCName_SharingSession)
		if err := sharing_session_collection.Find(bson.M{"sessionkey": share_id}).One(&s); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session = s
		GCachedSharedSessions[share_id] = s
	}

	paths := make([]string, len(session.SharedFiles))
	for i, val := range session.SharedFiles {
		if path, err := ushare.LocateDataFileByDetailedInfo(session.ServerName, session.UserName, session.RoleName, val); err == nil {
			paths[i] = path
		}
	}

	ret_obj := SharedSessionWithImagePaths{session.ServerName, session.RoleName, JSONTime(session.SharedTime), session.SharedWords, paths}
	WriteJsonResponse(w, ret_obj)
}
