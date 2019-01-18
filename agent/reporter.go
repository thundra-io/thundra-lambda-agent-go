package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type reporter interface {
	Collect(messages []plugin.MonitoringDataWrapper)
	Report()
	ClearData()
	Reported() *uint32
	FlushFlag()
}

type reporterImpl struct {
	messageQueue []plugin.MonitoringDataWrapper

	// reported is a flag to prevent system from sending data twice in case of timeout
	reported *uint32
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

// Collect collects the data from plugins. If async is on, it sends the data immediately.
func (r *reporterImpl) Collect(messages []plugin.MonitoringDataWrapper) {
	defer mutex.Unlock()
	mutex.Lock()
	if shouldSendAsync == "true" {
		sendAsync(messages)
		return
	}
	r.messageQueue = append(r.messageQueue, messages...)
}

// Report sends the data to collector
func (r *reporterImpl) Report() {
	if shouldSendAsync == "false" || shouldSendAsync == "" {
		sendHttpReq(r.messageQueue)
	}
	atomic.CompareAndSwapUint32(r.reported, 0, 1)
}

// ClearData clears the reporter data
func (r *reporterImpl) ClearData() {
	r.messageQueue = r.messageQueue[:0]
}

// Reported returns reported
func (r *reporterImpl) Reported() *uint32 {
	return r.reported
}

// FlushFlag flushes the reported flag
func (r *reporterImpl) FlushFlag() {
	atomic.CompareAndSwapUint32(r.Reported(), 1, 0)
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

func sendHttpReq(messageQueue []plugin.MonitoringDataWrapper) {
	if plugin.DebugEnabled {
		fmt.Printf("MessageQueue:\n %+v \n", messageQueue)
	}
	b, err := json.Marshal(&messageQueue)
	if err != nil {
		fmt.Println("Error in marshalling ", err)
	}

	targetURL := collectorUrl + monitoringDataPath
	if plugin.DebugEnabled {
		fmt.Println("Sending HTTP request to Thundra collector: " + targetURL)
	}

	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(b))
	if err != nil {
		fmt.Println("Error http.NewRequest: ", err)
	}
	req.Close = true
	req.Header.Set("Authorization", "ApiKey "+plugin.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error client.Do(req): ", err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll(resp.Body): ", err)
	}
	if plugin.DebugEnabled {
		fmt.Println("response Status:", resp.Status)
		fmt.Println("response Headers:", resp.Header)
		fmt.Println("response Body:", string(body))
	}
	resp.Body.Close()
}
