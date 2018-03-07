package thundra

import "sync"

var msgQueue []Message
var mutex = &sync.Mutex{}

func collect(msg Message) {
	defer mutex.Unlock()
	mutex.Lock()
	msgQueue = append(msgQueue, msg)
}

func report() {
	if ShouldSendAsync == "false" {
		sendHttpReq(msgQueue)
	}
}
