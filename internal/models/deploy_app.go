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
