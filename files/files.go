package files

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"dynamicledger.com/testnet-deployer/helper"
)

type SetupStore struct {
	Identity     string          `json:"identity"`
	Namespace    string          `json:"namespace"`
	KeyData      KeyStore        `json:"keyData"`
	ContractData []ContractStore `json:"contractData"`
}

type KeyStore struct {
	PublicPem   string `json:"publicPem"`
	PublicHash  string `json:"publicHash"`
	PrivatePem  string `json:"privatePem"`
	PrivateHash string `json:"privateHash"`
}

type ContractStore struct {
	Name string `json:"name"`
	ID   string `json:"id"`
	Hash string `json:"hash"`
}

func SaveSetupData(data *helper.SetupData, contractData *[]ContractStore, path string) {
	prvKey := data.KeyHandler.GetPrivatePEM()
	pubKey := data.KeyHandler.GetPublicPEM()

	h := sha256.New()
	h.Write([]byte(prvKey))

	prvHashBytes := h.Sum(nil)

	h.Reset()
	h.Write([]byte(pubKey))
	pubHashBytes := h.Sum(nil)

	prvHash := fmt.Sprintf("%x", prvHashBytes)
	pubHash := fmt.Sprintf("%x", pubHashBytes)

	keyData := KeyStore{
		PrivatePem:  prvKey,
		PublicPem:   pubKey,
		PrivateHash: prvHash,
		PublicHash:  pubHash,
	}

	toStore := SetupStore{
		KeyData:      keyData,
		ContractData: *contractData,
		Identity:     string(data.Identity),
		Namespace:    data.Namespace,
	}

	bData, err := json.Marshal(toStore)
	if err != nil {
		helper.HandleError(err, "Error marshalling data to store")
	}

	WriteFile(path, bData)

}

func ReadFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		helper.HandleError(err, fmt.Sprintf("Error reading file with path %s", path))
	}

	return data
}

func WriteFile(path string, data []byte) {
	if err := os.WriteFile(path, data, 0644); err != nil {
		helper.HandleError(err, fmt.Sprintf("Error writing data to file \"%s\"\n", path))
	}
}
