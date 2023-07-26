package contracts

import (
	"dynamicledger.com/testnet-deployer/files"

	b64 "encoding/base64"
)

func (ch *ContractHandler) loadContracts() {
	contracts := []Contract{}

	ch.getManifest()

	for k, c := range ch.manifest.Contracts {
		if c.Exclude {
			continue
		}

		con := ch.readContract(c)

		// Only hash if one isn't set, we only want to update hashes after uploading
		// updated contracts
		var blank string
		if c.Hash == blank {
			hash := getContractHash(con)
			ch.manifest.Contracts[k].Hash = hash
		}

		contracts = append(contracts, con)
	}

	// Write the manifest file in case there are new hashes
	ch.storeManifest()

	ch.contracts = contracts
}

func (ch *ContractHandler) readContract(contractMeta contractMetadata) Contract {
	c := Contract{}
	c.name = contractMeta.Name
	c.version = contractMeta.Version

	data := files.ReadFile(contractMeta.Path)

	c.data = b64.StdEncoding.EncodeToString(data)

	return c
}
