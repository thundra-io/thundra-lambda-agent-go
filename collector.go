package thundra

var msgQueue []Message

func collect(msg Message) {
	//TODO mutex
	msgQueue = append(msgQueue, msg)
}

func report() {
	if ShouldSendAsync == "false" {
		sendHttpReq(msgQueue)
	}
}
