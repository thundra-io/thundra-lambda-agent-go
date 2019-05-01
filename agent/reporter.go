package agent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/thundra-io/thundra-lambda-agent-go/config"

	"github.com/thundra-io/thundra-lambda-agent-go/constants"
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
	client       *http.Client
	reported     *uint32
}

var shouldSendAsync string
var collectorURL string
var mutex = &sync.Mutex{}

func init() {
	shouldSendAsync = os.Getenv(constants.ThundraLambdaPublishCloudwatchEnable)
	if url := os.Getenv(constants.ThundraLambdaReportRestBaseURL); url != "" {
		collectorURL = url
	} else {
		collectorURL = constants.DefaultCollectorURL
	}
}

func newReporter() *reporterImpl {
	return &reporterImpl{
		client:   createHTTPClient(),
		reported: new(uint32),
	}
}

// Collect collects the data from plugins. If async is on, it sends the data immediately.
func (r *reporterImpl) Collect(messages []plugin.MonitoringDataWrapper) {
	defer mutex.Unlock()
	mutex.Lock()
	if shouldSendAsync == "true" && !config.ReportCloudwatchCompositeDataEnabled {
		sendAsync(messages)
		return
	}
	r.messageQueue = append(r.messageQueue, messages...)
}

// Report sends the data to collector
func (r *reporterImpl) Report() {
	atomic.CompareAndSwapUint32(r.reported, 0, 1)
	if shouldSendAsync == "false" || shouldSendAsync == "" {
		r.sendHTTPReq()
	} else if config.ReportCloudwatchCompositeDataEnabled {
		r.sendAsyncComposite()
	}
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
	switch v := msg.(type) {
	case []plugin.MonitoringDataWrapper:
		for i := range v {
			b, err := json.Marshal(v[i])
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println("Sending ASYNC request to Thundra collector")
			fmt.Println(string(b))
		}
	default:
		b, err := json.Marshal(&msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Sending ASYNC request to Thundra collector")
		fmt.Println(string(b))
	}
}

func (r *reporterImpl) sendAsyncComposite() {
	batchSize := config.ReportCloudwatchCompositeBatchSize
	for i := 0; i < len(r.messageQueue); i += batchSize {
		end := i + batchSize
		if end > len(r.messageQueue) {
			end = len(r.messageQueue)
		}
		baseData := plugin.PrepareBaseData()
		compositeData := plugin.WrapMonitoringData(plugin.PrepareCompositeData(baseData, r.messageQueue[i:end]), "Composite")
		sendAsync(compositeData)
	}
}

func (r *reporterImpl) sendHTTPReq() {
	if config.DebugEnabled {
		fmt.Printf("MessageQueue:\n %+v \n", r.messageQueue)
	}
	targetURL := collectorURL + constants.MonitoringDataPath
	if config.ReportRestCompositeDataEnabled {
		targetURL = collectorURL + constants.CompositeMonitoringDataPath
	}

	if config.DebugEnabled {
		fmt.Println("Sending HTTP request to Thundra collector: " + targetURL)
	}

	batchSize := config.ReportRestCompositeBatchSize
	var wg sync.WaitGroup
	ch := make(chan string)
	for i := 0; i < len(r.messageQueue); i += batchSize {

		end := i + batchSize

		if end > len(r.messageQueue) {
			end = len(r.messageQueue)
		}
		if config.ReportRestCompositeDataEnabled {
			baseData := plugin.PrepareBaseData()
			compositeData := plugin.WrapMonitoringData(plugin.PrepareCompositeData(baseData, r.messageQueue[i:end]), "Composite")

			b, err := json.Marshal(compositeData)
			if err != nil {
				fmt.Println("Error in marshalling ", err)
				return
			}
			wg.Add(1)
			go r.sendBatch(targetURL, b, &wg, ch)
		} else {
			b, err := json.Marshal(r.messageQueue[i:end])
			if err != nil {
				fmt.Println("Error in marshalling ", err)
				return
			}
			wg.Add(1)
			go r.sendBatch(targetURL, b, &wg, ch)
		}
	}
	go printDebugLogs(ch)

	wg.Wait()

	close(ch)
}

func printDebugLogs(ch <-chan string) {
	for {
		select {
		case str := <-ch:
			fmt.Print(str)
		}
	}
}

func (r *reporterImpl) sendBatch(targetURL string, messages []byte, wg *sync.WaitGroup, ch chan<- string) {
	defer wg.Done()
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(messages))
	if err != nil {
		ch <- fmt.Sprintln("Error http.NewRequest: ", err)
		return
	}
	req.Close = true
	req.Header.Set("Authorization", "ApiKey "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		ch <- fmt.Sprintln("Error client.Do(req): ", err)
		return
	}
	if config.DebugEnabled {
		ch <- fmt.Sprintln("response Status:", resp.Status)
		ch <- fmt.Sprintln("response Headers:", resp.Header)
	}
	if resp.Body == nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ch <- fmt.Sprintln("ioutil.ReadAll(resp.Body): ", err)
	} else if config.DebugEnabled {
		ch <- fmt.Sprintln("response Body:", string(body))
	}

	resp.Body.Close()
}

func createHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: config.TrustAllCertificates,
		},
	}
	return &http.Client{Transport: tr}
}
