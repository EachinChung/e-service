package model

import "github.com/eachinchung/log"

type Status uint8

const (
	StatusNormal Status = iota
	StatusDelete
	StatusFreeze
	StatusDanger
)

var statusCodeMsgMap = map[Status]string{
	StatusNormal: "正常",
	StatusDelete: "删除",
	StatusFreeze: "冻结",
	StatusDanger: "风控",
}

func (s Status) Msg() string {
	msg, ok := statusCodeMsgMap[s]
	if !ok {
		msg = statusCodeMsgMap[StatusDelete]
		log.Errorf("监测到未定义的状态, Status: %d", s)
	}
	return msg
}
