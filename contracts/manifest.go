package contracts

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"

	"dynamicledger.com/testnet-deployer/files"
	"dynamicledger.com/testnet-deployer/structs"
)

func (ch *ContractHandler) getManifest() {
	var man structs.ContractManifest

	fHan := files.GetFileHandler(ch.Logger)
	data := fHan.ReadFile(ch.Config.ContractManifest)
	if err := json.Unmarshal(data, &man); err != nil {
		ch.Logger.Fatal(err, "Error parsing contract manifHan file")
	}

	ch.Manifest = man
}

func (ch *ContractHandler) addIDToManifest(contractName string, ID string) {

	for i, cMeta := range ch.Manifest.Contracts {
		if cMeta.Name == contractName {
			ch.Manifest.Contracts[i].ID = ID
		}
	}

	ch.storeManifest()

}

func (ch *ContractHandler) updateOnboardedStatus(ID string) {
	for i, cMeta := range ch.Manifest.Contracts {
		if cMeta.ID == ID {
			ch.Manifest.Contracts[i].Onboarded = true
		}
	}

	ch.storeManifest()
}

func (ch *ContractHandler) updateVersion(ID string, version string) {
	for i, cMeta := range ch.Manifest.Contracts {
		if cMeta.ID == ID {
			ch.Manifest.Contracts[i].Version = version
			break
		}
	}

	ch.storeManifest()
}

func (ch *ContractHandler) setHashes(missingCheck bool) {
	var blank string
	hasChanges := false

	for i, cMeta := range ch.Manifest.Contracts {

		c := ch.readContract(cMeta)

		if missingCheck && cMeta.Hash == blank {
			hash := getContractHash(c)
			ch.Manifest.Contracts[i].Hash = hash
			hasChanges = true
		}

		if !missingCheck {
			hash := getContractHash(c)

			if ch.Manifest.Contracts[i].Hash != hash {
				ch.Manifest.Contracts[i].Hash = hash
				hasChanges = true
			}
		}

	}

	if hasChanges {
		ch.storeManifest()
	}
}

func (ch *ContractHandler) storeManifest() {
	bMan, err := json.MarshalIndent(ch.Manifest, "", "  ")
	if err != nil {
		ch.Logger.Fatal(err, "Error marshalling manifest data")
	}

	fHan := files.GetFileHandler(ch.Logger)
	fHan.WriteFile(ch.Config.ContractManifest, bMan)
}

func getContractHash(contract structs.Contract) string {
	hasher := sha256.New()
	hasher.Write([]byte(contract.Data))

	hash := hasher.Sum(nil)

	encoded := base64.StdEncoding.EncodeToString(hash)

	return encoded
}
