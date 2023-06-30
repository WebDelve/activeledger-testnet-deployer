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
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DefaultIdentity   string `json:"identity"`
	DefaultNamespace  string `json:"namespace"`
	ContractFolder    string `json:"contractDir"`
	ContractManifest  string `json:"contractManifest"`
	SetupDataSaveFile string `json:"setupDataSaveFile"`
	TestnetFolder     string `json:"testnetFolder"`
}

type SetupStore struct {
	PrivatePem string `json:"privatePem"`
	Identity   string `json:"identity"`
	Namespace  string `json:"namespace"`
}

func loadConfig() *Config {
	var c Config

	data := readFile("./config.json")

	if err := json.Unmarshal(data, &c); err != nil {
		handleError(err, "Error unmarshalling content of config file")
	}

	return &c
}

func saveSetupData(data *SetupData, path string) {

	toStore := SetupStore{
		PrivatePem: data.keyHandler.GetPrivatePEM(),
		Identity:   string(data.identity),
		Namespace:  data.namespace,
	}

	bData, err := json.Marshal(toStore)
	if err != nil {
		handleError(err, "Error marshalling data to store")
	}

	if err = os.WriteFile(path, bData, 0644); err != nil {
		handleError(err, fmt.Sprintf("Error writing data to file \"%s\"\n", path))
	}
}

func readFile(path string) []byte {
	data, err := os.ReadFile(path)
	if err != nil {
		handleError(err, fmt.Sprintf("Error reading file with path %s", path))
	}

	return data
}
