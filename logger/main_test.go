package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

// Test Bybit API(LS ratio)
func TestLSRest(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.bybit.com/v2/public/account-ratio?symbol=BTCUSD&period=5min", nil)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
}
