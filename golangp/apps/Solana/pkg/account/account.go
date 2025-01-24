package account

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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

func CreateAccounts(folderPath string, count int) ([]string, error) {
	var publicKeys []string

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return nil, fmt.Errorf("unable to create folder: %v", err)
	}

	for i := 0; i < count; i++ {
		account := types.NewAccount()

		truncatedPublicKey := account.PublicKey.ToBase58()[:10]
		filename := fmt.Sprintf("%s/wallet_%s.json", folderPath, truncatedPublicKey)

		file, err := os.Create(filename)
		if err != nil {
			return nil, fmt.Errorf("unable to create file: %v", err)
		}
		defer file.Close()

		if err := json.NewEncoder(file).Encode(account.PrivateKey); err != nil {
			return nil, fmt.Errorf("unable to write to file: %v", err)
		}

		fmt.Printf("New account created with Public Key: %s, Private Key saved to %s\n", account.PublicKey.ToBase58(), filename)

		publicKeys = append(publicKeys, account.PublicKey.ToBase58())
	}

	return publicKeys, nil
}