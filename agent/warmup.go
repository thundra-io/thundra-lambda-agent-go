package agent

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Please visit https://github.com/thundra-io/thundra-lambda-warmup to learn more about warmup.

// checkAndHandleWarmupRequest is used to detect thundra-lambda-warmup requests.
// If the incoming request is a warmup request thundra will return nil and stop execution.
func checkAndHandleWarmupRequest(event interface{}, eventType reflect.Type) bool {
	// Check if the event is an empty struct
	if isZeroEvent(event, eventType) {
		log.Println("Received warmup request as empty message. Handling with 100 milliseconds delay ...")
		time.Sleep(time.Millisecond * 100)
		return true
	}

	if eventType.Kind() == reflect.String {
		// Check whether it is a warmup request
		s := fmt.Sprint(event)
		if strings.HasPrefix(s, "#warmup") {
			delay := 100

			// Warmup data has the following format "#warmup wait=200 k1=v1"
			//Therefore we need to parse it to only have arguments in k=v format
			sp := strings.SplitAfter(s, "#warmup")[1]
			args := strings.Fields(sp)
			// Iterate over all warmup arguments
			for _, a := range args {
				argParts := strings.Split(a, "=")
				// Check whether argument is in key=value format
				if len(argParts) == 2 {
					k := argParts[0]
					v := argParts[1]
					// Check whether argument is "wait" argument
					// which specifies extra wait time before returning from request
					if k == "wait" {
						w, err := strconv.Atoi(v)
						if err != nil {
							log.Println(err)
						} else {
							delay += w
						}
					}
				}
			}
			log.Println("Received warmup request as warmup message. Handling with ", delay, " milliseconds delay ...")
			time.Sleep(time.Millisecond * time.Duration(delay))
			return true
		}
	}

	return false
}

// isZeroEvent compares whether event is equals to the zero value for event's type.
func isZeroEvent(event interface{}, eventType reflect.Type) bool {
	zeroEvent := reflect.Zero(eventType)

	// Use sprint to compare each individual field values of event with zero event
	ev := fmt.Sprint(event)
	ze := fmt.Sprint(zeroEvent)
	if ev == ze {
		return true
	}
	return false
}
