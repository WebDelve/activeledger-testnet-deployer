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
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	alsdk "github.com/activeledger/SDK-Golang/v2"
)

type SetupData struct {
	folder     string
	identity   alsdk.StreamID
	namespace  string
	keyHandler alsdk.KeyHandler
	conn       alsdk.Connection
}

func main() {
	config := loadConfig()

	setupData := SetupData{}

	setupData.conn = alsdk.Connection{
		Protocol: alsdk.HTTP,
		Url:      "localhost",
		Port:     "5260",
	}

	setupData.namespace = config.DefaultNamespace

	setupData.folder = bootstrap(config)

	createTestnet(&setupData)
	createIden(&setupData, config)
	createNamespace(&setupData)
	onboardSmartContracts(&setupData, config)
	saveSetupData(&setupData, config.SetupDataSaveFile)

	fmt.Printf("\n\nBootstrapping complete\n\n")
}

func bootstrap(config *Config) string {

	fmt.Printf("Input base folder name, leave blank for default (default=%s): ", config.TestnetFolder)

	var folder string
	fmt.Scanln(&folder)

	var blank string
	if folder == blank {
		folder = config.TestnetFolder
	}

	if err := os.Mkdir(folder, 0755); err != nil {
		handleError(err, "Error creating base folder")
	}

	return folder
}

func createTestnet(setupData *SetupData) {

	cmd := exec.Command("activeledger", "--testnet")
	cmd.Dir = setupData.folder

	fmt.Println("Creating testnet")
	_, err := cmd.Output()
	if err != nil {
		handleError(err, "Error running activeledger --testnet, is activeledger installed?")
	}
	fmt.Println("Testnet created")

	fmt.Println("Leave this process running and start the testnet, navigate to the created folder and run 'node testnet'. When done return here and press enter to continue.")
	fmt.Scanln()
}

func createIden(setupData *SetupData, config *Config) {
	fmt.Println("Creating and onboarding identity...")

	keyhandler, err := alsdk.GenerateRSA()
	if err != nil {
		handleError(err, "Error generating RSA key")
	}

	setupData.keyHandler = keyhandler

	pubpem := keyhandler.GetPublicPEM()

	input := alsdk.DataWrapper{
		"publicKey": pubpem,
		"type":      "rsa",
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  alsdk.StreamID(config.DefaultIdentity),
		Contract:  "onboard",
		Key:       keyhandler,
		Namespace: "default",
		SelfSign:  true,
		Input:     input,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		handleError(err, "Error building onboarding transaction")
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(*tx, setupData.conn)
	if err != nil {
		handleALError(err, resp, "Error sending identity onboarding transaction")
	}

	var streamId string
	for _, v := range resp.Streams.New {
		if v.ID != streamId {
			streamId = v.ID
			break
		}
	}

	var blank string
	if streamId == blank {
		handleError(errors.New("Blank streamID"), "Identity transaction didn't error, but we got no stream ID")
	}

	setupData.identity = alsdk.StreamID(streamId)
}

func createNamespace(setupData *SetupData) {
	fmt.Println("Creating Namespace...")

	input := alsdk.DataWrapper{
		"namespace": setupData.namespace,
	}

	txOpts := alsdk.TransactionOpts{
		StreamID:  setupData.identity,
		Contract:  "namespace",
		Namespace: "default",
		Input:     input,
		Key:       setupData.keyHandler,
	}

	txHan, _, err := alsdk.BuildTransaction(txOpts)
	if err != nil {
		handleError(err, "Error building namespace transaction")
	}

	tx := txHan.GetTransaction()

	resp, err := alsdk.Send(*tx, setupData.conn)
	if err != nil {
		handleALError(err, resp, "Error sending namespace transaction")
	}

	fmt.Printf("Namespace '%s' created\n", setupData.namespace)
}

func onboardSmartContracts(setupData *SetupData, config *Config) {

	fmt.Println("Onboarding contracts...")

	onboardContracts(setupData, config)

}

func handleError(e error, note string) {
	log.Println(note)
	log.Fatalln(e)
}

func handleALError(e error, resp alsdk.Response, note string) {
	if len(resp.Summary.Errors) > 0 {
		for i, e := range resp.Summary.Errors {
			log.Printf("Error %d: %s\n\n", i, e)
		}
	}

	handleError(e, note)
}
