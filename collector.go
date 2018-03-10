package thundra

//Collector is responsible for collecting and sending the data

import (
	"sync"
	"os"
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
)

var ShouldSendAsync string

func init() {
	ShouldSendAsync = os.Getenv(ThundraLambdaPublishCloudwatchEnable)
}

type collector interface {
	collect(msg Message)
	report()
	clear()
}

type collectorImpl struct {
	msgQueue []Message
}

var mutex = &sync.Mutex{}

func (c *collectorImpl) collect(msg Message) {
	defer mutex.Unlock()
	mutex.Lock()
	if ShouldSendAsync == "true" {
		sendAsync(msg)
		return
	}
	c.msgQueue = append(c.msgQueue, msg)
}

func (c *collectorImpl) report() {
	if ShouldSendAsync == "false" {
		sendHttpReq(c.msgQueue)
	}
}

func (c *collectorImpl) clear() {
	c.msgQueue = c.msgQueue[:0]
}

func sendAsync(msg Message) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Sending ASYNC request")
	fmt.Println(string(b))
}

func sendHttpReq(msg []Message) {
	b, _ := json.Marshal(&msg)
	req, _ := http.NewRequest("POST", collectorUrl, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "ApiKey "+ApiKey)
	req.Header.Set("Content-Type", "application/json")

	fmt.Println("Sending HTTP request")
	fmt.Println(msg)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	fmt.Println("response Status:", resp.Status)
	//TODO if resp.status == 401 unauthorized : ApiKey is missing
	fmt.Println("ApiKey:", ApiKey)
	resp.Body.Close()
}
