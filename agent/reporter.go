package agent

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"sync/atomic"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/config"

	"github.com/thundra-io/thundra-lambda-agent-go/v2/constants"
	"github.com/thundra-io/thundra-lambda-agent-go/v2/plugin"
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

var collectorURL string
var mutex = &sync.Mutex{}

func init() {
	if url := os.Getenv(constants.ThundraLambdaReportRestBaseURL); url != "" {
		collectorURL = url
	} else {
		collectorURL = config.CollectorUrl
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
	if config.ReportCloudwatchEnabled && !config.ReportCloudwatchCompositeDataEnabled {
		sendAsync(messages)
		return
	}
	r.messageQueue = append(r.messageQueue, messages...)
}

// Report sends the data to collector
func (r *reporterImpl) Report() {
	atomic.CompareAndSwapUint32(r.reported, 0, 1)
	if !config.ReportCloudwatchEnabled {
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

func sendAsync(data []plugin.MonitoringDataWrapper) {
	for i := range data {
		b, err := json.Marshal(data[i])
		if err != nil {
			log.Println(err)
			return
		}
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
		compositeData := plugin.PrepareCompositeData(baseData, r.messageQueue[i:end])
		wrappedCompositeData := plugin.WrapMonitoringData(compositeData, "Composite")
		sendAsync([]plugin.MonitoringDataWrapper{wrappedCompositeData})
	}
}

func (r *reporterImpl) sendHTTPReq() {
	if config.DebugEnabled {
		log.Printf("MessageQueue:\n %+v \n", r.messageQueue)
	}
	targetURL := collectorURL + constants.MonitoringDataPath
	if config.ReportRestCompositeDataEnabled {
		targetURL = collectorURL + constants.CompositeMonitoringDataPath
	}

	if config.DebugEnabled {
		log.Println("Sending HTTP request to Thundra collector: " + targetURL)
	}

	batchSize := config.ReportRestCompositeBatchSize
	var wg sync.WaitGroup
	for i := 0; i < len(r.messageQueue); i += batchSize {

		end := i + batchSize

		if end > len(r.messageQueue) {
			end = len(r.messageQueue)
		}
		if config.ReportRestCompositeDataEnabled {
			baseData := plugin.PrepareBaseData()
			compositeData := plugin.PrepareCompositeData(baseData, r.messageQueue[i:end])
			wrappedCompositeData := plugin.WrapMonitoringData(compositeData, "Composite")

			b, err := json.Marshal(wrappedCompositeData)
			if err != nil {
				log.Println("Error in marshalling ", err)
				return
			}
			wg.Add(1)
			go r.sendBatch(targetURL, b, &wg)
		} else {
			b, err := json.Marshal(r.messageQueue[i:end])
			if err != nil {
				log.Println("Error in marshalling ", err)
				return
			}
			wg.Add(1)
			go r.sendBatch(targetURL, b, &wg)
		}
	}
	wg.Wait()
}

func (r *reporterImpl) sendBatch(targetURL string, messages []byte, wg *sync.WaitGroup) {
	defer wg.Done()
	req, err := http.NewRequest("POST", targetURL, bytes.NewBuffer(messages))
	if err != nil {
		log.Println("Error http.NewRequest:", err)
		return
	}
	req.Close = true
	req.Header.Set("Authorization", "ApiKey "+config.APIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		log.Println("Error client.Do(req):", err)
		return
	}
	if config.DebugEnabled {
		log.Println("response Status:", resp.Status)
		log.Println("response Headers:", resp.Header)
	}
	if resp.Body == nil {
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ioutil.ReadAll(resp.Body): ", err)
	} else if config.DebugEnabled {
		log.Println("response Body:", string(body))
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
