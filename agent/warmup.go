package agent

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"
)

// Please visit https://github.com/thundra-io/thundra-lambda-warmup to learn more about warmup.

// checkAndHandleWarmupRequest is used to detect thundra-lambda-warmup requests.
// If the incoming request is a warmup request thundra will return nil and stop execution.
func checkAndHandleWarmupRequest(payload json.RawMessage) bool {
	if json.Valid(payload) {
		paylaodStr := string(payload)
		if strings.HasPrefix(paylaodStr, `"#warmup`) {
			paylaodStr, err := strconv.Unquote(paylaodStr)
			if err != nil {
				log.Println("Bad string format while checking warmup")
				return false
			}
			paylaodStr = strings.TrimLeft(paylaodStr, " ")

			// Warmup data has the following format "#warmup wait=200 k1=v1"
			//Therefore we need to parse it to only have arguments in k=v format
			delay := 0
			sp := strings.SplitAfter(paylaodStr, "#warmup")[1]
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

		} else {
			j := make(map[string]interface{})
			err := json.Unmarshal(payload, &j)
			if err != nil {
				log.Println("Bad json format while checking warmup")
				return false
			}

			if len(j) == 0 {
				log.Println("Received warmup request as empty message. Handling with 100 milliseconds delay ...")
				time.Sleep(time.Millisecond * 100)
				return true
			} else if len(j) > 0 {
				body := j["body"]
				if body == nil {
					log.Println("Received warmup request as empty message. Handling with 100 milliseconds delay ...")
					time.Sleep(time.Millisecond * 100)
					return true
				}
			}
		}

	}
	return false
}
