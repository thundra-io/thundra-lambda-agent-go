package plugins

import (
	"encoding/json"
	"fmt"
	"bytes"
	"net/http"
)

func sendReport(msg Message) {
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(b))

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(b))
	req.Header.Set("Authorization", "ApiKey "+msg.ApiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
}