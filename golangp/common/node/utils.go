package node

import (
	"crypto/rand"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/blocto/solana-go-sdk/types"
	"github.com/mr-tron/base58"
	"github.com/reiver/go-telnet"
)

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

func SolCreateAccount() (account types.Account, publicKey string, privateKey string) {
	account = types.NewAccount()
	publicKey = account.PublicKey.ToBase58()
	privateKey = base58.Encode(account.PrivateKey)
	return
}

func SolRestoreAccount(privateKey string) (account types.Account, publicKey string, err error) {
	account, err = types.AccountFromBase58(privateKey)
	if err != nil {
		return
	}
	publicKey = account.PublicKey.ToBase58()
	return
}

func ExecSystemCmd(cmdStr string, block bool) {
	log.Printf("ExecSystemCmd: %v", cmdStr)
	cmdParts := strings.Split(cmdStr, " ")
	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Printf("ExecSystemCmd error:\n%s\n", err)
	}
	if block {
		cmd.Wait()
	}
}

func Telnet(srvAddr string) error {
	caller := telnet.StandardCaller
	return telnet.DialToAndCall(srvAddr, caller)
}
