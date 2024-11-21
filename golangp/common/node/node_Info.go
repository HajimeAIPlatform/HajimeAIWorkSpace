package node

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"hajime/golangp/common/logging"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/carlmjohnson/requests"
)

func LoadNodeInfo() (nodeInfo *VerifierNodeInfo) {
	contents, err := os.ReadFile(NODE_INFO_PATH_FILE)
	if err != nil {
		account, publicKey, privateKey := SolCreateAccount()
		nodeInfo = &VerifierNodeInfo{
			account:    account,
			publicKey:  publicKey,
			PrivateKey: privateKey,
		}
		nodeInfo.Save()

		return
	}

	nodeInfo = &VerifierNodeInfo{}
	json.Unmarshal(contents, &nodeInfo)
	account, publicKey, err := SolRestoreAccount(nodeInfo.PrivateKey)
	if err != nil {
		logging.Warning("SolRestoreAccount ERROR: %v\n", err)
		return
	}

	nodeInfo.account = account
	nodeInfo.publicKey = publicKey
	return
}

func (nodeInfo *VerifierNodeInfo) Save() {
	contentBytes, _ := json.Marshal(nodeInfo)
	err := os.WriteFile(NODE_INFO_PATH_FILE, contentBytes, 0644)
	if err != nil {
		logging.Warning("WriteFile(%v) ERROR: %v\n", NODE_INFO_PATH_FILE, err)
	}
}

func (nodeInfo *VerifierNodeInfo) HandleMsgLogin(msg *MsgLogin) {
	nodeInfo.Id = msg.Data.Id
	nodeInfo.Cmd = msg.Data.Cmd
	nodeInfo.Ip = msg.Data.Ip

	nodeInfo.Save()

	if runtime.GOOS != "windows" {
		ExecSystemCmd(nodeInfo.Cmd, false)
	}
}

func (nodeInfo *VerifierNodeInfo) GetAccessToken(td *CheckTaskData) (token string, err error) {
	err = Telnet(fmt.Sprintf("%v:8080", td.NodeIp))
	if err != nil {
		logging.Warning("GetAccessToken telnet err: %v\n", err)
		return
	}
	now := time.Now().Unix()
	sig := nodeInfo.account.Sign([]byte(strconv.FormatInt(now, 10)))
	sigStr := hex.EncodeToString(sig)
	url := fmt.Sprintf("http://%v:8080/auth", td.NodeIp)
	body := &AuthRequest{
		PublicKey: nodeInfo.publicKey,
		Sig:       sigStr,
		Ts:        now,
	}
	authRsp := new(AuthResponse)
	err = requests.URL(url).BodyJSON(&body).
		ToJSON(&authRsp).Fetch(context.Background())
	if err != nil {
		logging.Warning("GetAccessToken fetch token err: %v\n", err)
		return
	}

	token = authRsp.Data
	return
}
