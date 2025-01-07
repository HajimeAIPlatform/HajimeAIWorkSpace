package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/blocto/solana-go-sdk/types"
)

// LoadAccountFromFile loads a Solana account from a private key file.
func LoadAccountFromFile(filePath string) (types.Account, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to read file: %v", err)
	}

	var keypair []byte
	err = json.Unmarshal(data, &keypair)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to unmarshal private key: %v", err)
	}

	account, err := types.AccountFromBytes(keypair)
	if err != nil {
		return types.Account{}, fmt.Errorf("failed to create account from bytes: %v", err)
	}

	return account, nil
}
