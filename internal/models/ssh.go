package models

type SshCreds struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"pass"`
}

type SshDeployAppReq struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type SshDeployAppResp struct {
	Log string `json:"log"`
}

type SshDeployAppErrorResp struct {
	Log   string `json:"log"`
	Error string `json:"error"`
}
