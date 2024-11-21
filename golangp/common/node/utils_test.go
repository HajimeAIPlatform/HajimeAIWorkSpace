package node

import (
	"testing"
)

func TestSolCreateAccount(t *testing.T) {
	account, publicKey, privateKey := SolCreateAccount()
	t.Logf("account: %v\n", account)
	t.Logf("publicKey: %v\n", publicKey)
	t.Logf("privateKey: %v\n", privateKey)

	newAccount, newPublicKey, err := SolRestoreAccount(privateKey)
	if err != nil {
		t.Errorf("SolRestoreAccount failed: %v\n", err)
		return
	}

	t.Logf("newAccount: %v\n", newAccount)
	t.Logf("newPublicKey: %v\n", newPublicKey)

	t.Log("Test_Sol_Create_Account ok")
}

func TestSystemCmd(t *testing.T) {
	cmdStr := "edge -c hajime -k hajimegogogo! -a 10.10.11.64/16 -f -l n2n.dorylus.chat:6777 -d n2n1"
	ExecSystemCmd(cmdStr, true)
	t.Log("TestSystemCmd ok")
}

func TestTelnet(t *testing.T) {
	err := Telnet("10.10.10.12:11434")
	if err != nil {
		t.Fatalf("TestTelnet err: %v\n", err)
	}
	t.Logf("TestTelnet ok\n")
}
