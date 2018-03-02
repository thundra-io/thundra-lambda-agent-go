package thundra

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
	"os"
	"thundra-agent-go/constants"
)

func sendReport(msg Message,url string) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}

	if os.Getenv(constants.Thundra_Lambda_Publish_Cloudwatch_Enable) == "true"{
		sendAsync(b)
	} else{
		sendHttpReq(b, msg.ApiKey,url)
	}
}

func sendAsync(msg []byte) {
	fmt.Println("Sending ASYNC request")
	fmt.Println(string(msg))
}

func sendHttpReq(msg []byte, apiKey string, url string) {
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(msg))
	req.Header.Set("Authorization", "ApiKey "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Sending HTTP request")
	fmt.Println(string(msg))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
}
