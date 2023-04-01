package models

import "net/netip"

type SshCreds struct {
	Addr     netip.AddrPort
	Login    string
	Password string
}

type SshDeployAppReq struct {
	IP       string `json:"ip"`
	Port     uint16 `json:"port"`
	Login    string `json:"user"`
	Password string `json:"password"`
}

type DeployAppStatus string

const (
	STATUS_IN_QUEUE   DeployAppStatus = "in queue"
	STATUS_START      DeployAppStatus = "started"
	STATUS_CONN_ERR   DeployAppStatus = "connection error"
	STATUS_IN_PROCESS DeployAppStatus = "in process"
	STATUS_ERROR      DeployAppStatus = "error"
	STATUS_SUCCESS    DeployAppStatus = "success"
)

type TaskProgressMsg struct {
	Log     string                `json:"log"`
	Percent uint8                 `json:"percent"`
	Error   string                `json:"error"`
	Status  DeployAppStatus       `json:"status"`
	Chan    *chan TaskProgressMsg `json:"-"`
	TaskId  uint64                `json:"-"`
}

type Error struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}
