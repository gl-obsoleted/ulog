package ushare

import (
	"crypto/sha256"
	"fmt"
	"io"
	"time"
	"ushare/core"
)

type UserSession struct {
	RegionName string
	ServerName string
	UserName   string
	RoleName   string
	Timestamp  time.Time
	Ticket     string
}

type activeUserLut map[string]UserSession

var GActiveUsers activeUserLut = activeUserLut{}

func NewUser(info []string, salt []byte) (string, error) {
	if len(info) != 4 || len(salt) == 0 {
		return "", core.NewStdErr(core.ERR_NewUserFailed, "invalid user info.")
	}

	us := UserSession{info[0], info[1], info[2], info[3], time.Time{}, ""}

	for {
		t := time.Now()
		h := sha256.New()
		io.WriteString(h, us.RegionName)
		io.WriteString(h, us.ServerName)
		io.WriteString(h, us.UserName)
		io.WriteString(h, us.RoleName)
		io.WriteString(h, t.String())
		io.WriteString(h, string(salt)) // apply salt, which varies from each server session
		ticket := fmt.Sprintf("%x", h.Sum(nil))

		if _, exists := GActiveUsers[ticket]; !exists {
			us.Timestamp = t
			us.Ticket = ticket
			break
		}
	}

	GActiveUsers[us.Ticket] = us
	return us.Ticket, nil
}

func FindUser(ticket string) *UserSession {
	if us, exists := GActiveUsers[ticket]; exists {
		return &us
	}

	return nil
}
