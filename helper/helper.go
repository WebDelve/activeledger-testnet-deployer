package helper

import (
	"log"

	alsdk "github.com/activeledger/SDK-Golang/v2"
)

func HandleError(e error, note string) {
	log.Println(note)
	log.Fatalln(e)
}

func HandleALError(e error, resp alsdk.Response, note string) {
	if len(resp.Summary.Errors) > 0 {
		for i, e := range resp.Summary.Errors {
			log.Printf("Error %d: %s\n\n", i, e)
		}
	}

	HandleError(e, note)
}
