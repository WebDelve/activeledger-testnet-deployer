package files

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"

	"dynamicledger.com/testnet-deployer/helper"
	"dynamicledger.com/testnet-deployer/structs"
	alsdk "github.com/activeledger/SDK-Golang/v2"
)

func SaveSetupData(data *structs.SetupData, contractData []structs.ContractStore, path string) {
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

	keyData := structs.KeyStore{
		PrivatePem:  prvKey,
		PublicPem:   pubKey,
		PrivateHash: prvHash,
		PublicHash:  pubHash,
	}

	toStore := structs.SetupStore{
		KeyData:      keyData,
		ContractData: contractData,
		Identity:     string(data.Identity),
		Namespace:    data.Namespace,
	}

	bData, err := json.Marshal(toStore)
	if err != nil {
		helper.HandleError(err, "Error marshalling data to store")
	}

	WriteFile(path, bData)

}

func ReadSetupData(config *structs.Config) structs.SetupData {
	bSetup := ReadFile(config.SetupDataSaveFile)

	var setupStore structs.SetupStore
	if err := json.Unmarshal(bSetup, &setupStore); err != nil {
		helper.HandleError(err, "Error unmarshalling setup data")
	}

	keyHan, err := alsdk.SetKey(setupStore.KeyData.PrivatePem, alsdk.RSA)
	if err != nil {
		helper.HandleError(err, "Error setting up key handler")
	}

	connection := alsdk.Connection{
		Protocol: alsdk.HTTP,
		Url:      "localhost",
		Port:     "5260",
	}

	setup := structs.SetupData{
		Folder:     config.TestnetFolder,
		Identity:   alsdk.StreamID(setupStore.Identity),
		Namespace:  setupStore.Namespace,
		KeyHandler: keyHan,
		Conn:       connection,
		Contracts:  setupStore.ContractData,
	}

	return setup
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
