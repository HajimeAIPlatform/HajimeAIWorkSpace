package node

import "github.com/blocto/solana-go-sdk/types"

type VerifierNodeInfo struct {
	account    types.Account
	publicKey  string
	PrivateKey string `json:"private_key"`
	Id         int    `json:"id"`
	Cmd        string `json:"cmd"`
	Ip         string `json:"ip"`
}

type AuthRequest struct {
	PublicKey string `json:"public_key"`
	Sig       string `json:"sig"`
	Ts        int64  `json:"ts"`
}

type MsgBase struct {
	MsgType string `json:"msgType"`
}

type MsgLogin struct {
	MsgBase
	Data struct {
		Id      int    `json:"id"`
		Cmd     string `json:"cmd"`
		Healthy int    `json:"healthy"`
		Type    string `json:"type"`
		Imei    string `json:"imei"`
		Ip      string `json:"ip"`
	} `json:"data"`
}

type AuthResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data string `json:"data"`
}

type CheckTaskData struct {
	NodeIp string `json:"node_ip"`
	Imei   string `json:"imei"`
	TaskId string `json:"task_id"`
}

type CheckTaskResultData struct {
	CheckTaskData
	TaskResult map[string]bool `json:"task_result"`
}

type MsgCheckTask struct {
	MsgBase
	Data CheckTaskData `json:"data"`
}

type MsgCheckTaskResult struct {
	MsgBase
	Data CheckTaskResultData `json:"data"`
}
