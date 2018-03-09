package thundra

import "sync"
type Collector interface{
	collect(msg Message)
	report()
	clear()
}

type collectorImpl struct{
	msgQueue []Message
}

var mutex = &sync.Mutex{}

func (c *collectorImpl)collect(msg Message) {
	defer mutex.Unlock()
	mutex.Lock()
	c.msgQueue = append(c.msgQueue, msg)
}

func (c *collectorImpl)report() {
	if ShouldSendAsync == "false" {
		sendHttpReq(c.msgQueue)
	}
}

func (c *collectorImpl)clear(){
	c.msgQueue = c.msgQueue[:0]
}