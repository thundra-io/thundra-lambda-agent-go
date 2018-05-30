package thundra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
)

type reporter interface {
	Collect(messages []interface{})
	Report(apiKey string)
	Clear()
}

type reporterImpl struct {
	messageQueue []interface{}
}

var shouldSendAsync string
var collectorUrl string
var mutex = &sync.Mutex{}

func init() {
	shouldSendAsync = os.Getenv(thundraLambdaPublishCloudwatchEnable)
	if url := os.Getenv(thundraLambdaPublishRestBaseUrl); url != "" {
		collectorUrl = url
	} else {
		collectorUrl = defaultCollectorUrl
	}
}

func (c *reporterImpl) Collect(messages []interface{}) {
	defer mutex.Unlock()
	mutex.Lock()
	if shouldSendAsync == "true" {
		sendAsync(messages)
		return
	}
	c.messageQueue = append(c.messageQueue, messages...)
}

func (c *reporterImpl) Report(apiKey string) {
	if shouldSendAsync == "false" || shouldSendAsync == "" {
		sendHttpReq(c.messageQueue, apiKey)
	}
}

func (c *reporterImpl) Clear() {
	c.messageQueue = c.messageQueue[:0]
}

func sendAsync(msg interface{}) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sending ASYNC request to Thundra collector")
	fmt.Println(string(b))
}

func sendHttpReq(mesageQueue []interface{}, apiKey string) {
	b, err := json.Marshal(&mesageQueue)
	if err != nil {
		fmt.Println(err)
	}

	targetURL := collectorUrl + monitorDatas
	if debugEnabled {
		fmt.Println("Sending HTTP request to Thundra collector: " + targetURL)
		fmt.Println(string(b))
	}

	req, _ := http.NewRequest("POST", targetURL, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "ApiKey "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if debugEnabled {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(body))
	}
	resp.Body.Close()
}
