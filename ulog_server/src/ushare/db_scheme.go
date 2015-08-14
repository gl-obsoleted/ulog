package ushare

import (
	"time"
)

const DBCName_SharingSession = "sharing_sessions"

type DBSharingSession struct {
	SessionKey string

	RegionName string
	ServerName string
	UserName   string
	RoleName   string

	SharedTime  time.Time
	SharedFiles []string
	SharedWords string
}

const DBCName_UserEvents = "user_events"

const (
	EVT_OK    = 0
	EVT_Error = 1
)

const (
	EVT_TicketAcquired = 0
	EVT_ImageUploaded  = 1
)

type DBUserEvent struct {
	Time   time.Time
	Ticket string
	ID     int
	Result int
	Desc   string
}

func DBLogNewUserEvent(ticket string, ID int, result int, desc string) error {
	// 记录到数据库
	var evt DBUserEvent = DBUserEvent{}
	evt.Time = time.Now()
	evt.Ticket = ticket
	evt.ID = ID
	evt.Result = result
	evt.Desc = desc
	return GDatabase.Insert(DBCName_UserEvents, &evt)
}
