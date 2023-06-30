/*
* MIT License (MIT)
* Copyright (c) 2023 WebDelve Ltd
*
* Permission is hereby granted, free of charge, to any person obtaining a copy
* of this software and associated documentation files (the "Software"), to deal
* in the Software without restriction, including without limitation the rights
* to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
* copies of the Software, and to permit persons to whom the Software is
* furnished to do so, subject to the following conditions:
*
* The above copyright notice and this permission notice shall be included in all
* copies or substantial portions of the Software.
*
* THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
* IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
* FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
* AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
* LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
* OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
* SOFTWARE.
 */

package main

import (
	b64 "encoding/base64"
	"encoding/json"
	"fmt"

	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type Contract struct {
	name    string
	data    string
	version string
}

type manifest struct {
	Contracts []contractMetadata `json:"contracts"`
}

type contractMetadata struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Version string `json:"version"`
	Exclude bool   `json:"exclude"`
}

func onboardContracts(setup *SetupData, config *Config) {
	fmt.Println("Loading contracts...")
	contracts := loadContracts(config)
	fmt.Println("Contracts loaded, onboarding...")

	for _, c := range contracts {
		onboardContract(c, setup)
	}

}

func onboardContract(contract Contract, setupData *SetupData) {
	fmt.Printf("\nOnboarding %s contract...\n", contract.name)

	input := alsdk.DataWrapper{
		"version":   contract.version,
		"namespace": setupData.namespace,
		"name":      contract.name,
		"contract":  contract.data,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  setupData.identity,
		Contract:  "contract",
		Namespace: "default",
		Input:     input,
		Key:       setupData.keyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		handleError(err, fmt.Sprintf("Error building contract onboard transaction for contract %s", contract.name))
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(*tx, setupData.conn)
	if err != nil {
		handleALError(err, resp, fmt.Sprintf("Error running onboarding transaction for contract %s", contract.name))
	}

	fmt.Printf("%s contract onboarded.\n", contract.name)
}

func loadContracts(config *Config) []Contract {
	contracts := []Contract{}

	contractManifest := readManifest(config)

	for _, c := range contractManifest.Contracts {

		if c.Exclude {
			continue
		}

		con := readContract(c)
		contracts = append(contracts, con)
	}

	return contracts
}

func readContract(contractMeta contractMetadata) Contract {
	c := Contract{}
	c.name = contractMeta.Name
	c.version = contractMeta.Version

	data := readFile(contractMeta.Path)

	c.data = b64.StdEncoding.EncodeToString(data)

	return c
}

func readManifest(config *Config) manifest {
	// man := manifest{}
	var man manifest

	data := readFile(config.ContractManifest)

	if err := json.Unmarshal(data, &man); err != nil {
		handleError(err, "Error parsing contract manifest file")
	}

	return man
}
