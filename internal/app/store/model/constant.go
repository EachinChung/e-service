package model

import "github.com/eachinchung/log"

type Status int16

const (
	StatusNormal Status = iota
	StatusDelete
	StatusDanger
)

var statusCodeMsgMap = map[Status]string{
	StatusNormal: "正常",
	StatusDelete: "注销",
	StatusDanger: "风控",
}

func (s Status) Msg() string {
	msg, ok := statusCodeMsgMap[s]
	if !ok {
		msg = statusCodeMsgMap[StatusDanger]
		log.Errorf("监测到未定义的状态, Status: %d", s)
	}
	return msg
}
