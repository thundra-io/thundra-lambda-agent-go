package thundra

import "sync"

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
	mutex.Lock()
	c.msgQueue = append(c.msgQueue, msg)
	mutex.Unlock()
}

func (c *collectorImpl) report() {
	if ShouldSendAsync == "false" {
		sendHttpReq(c.msgQueue)
	}
}

func (c *collectorImpl) clear() {
	c.msgQueue = c.msgQueue[:0]
}
