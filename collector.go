package thundra

import "sync"

type collector struct{
	msgQueue []Message
}

var mutex = &sync.Mutex{}

func (c *collector)collect(msg Message) {
	defer mutex.Unlock()
	mutex.Lock()
	c.msgQueue = append(c.msgQueue, msg)
}

func (c *collector)report() {
	if ShouldSendAsync == "false" {
		sendHttpReq(c.msgQueue)
	}
}

func (c *collector)clear(){
	//TODO not nil
	c.msgQueue = nil
}