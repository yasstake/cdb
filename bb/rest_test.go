package bb

import (
	"fmt"
	"testing"
)

func TestRestRequest(t *testing.T) {
	url := "https://api.bybit.com/v2/public/liq-records?symbol=BTCUSD"

	result, time, err := RestRequest(url)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(time.UTC().String())
	fmt.Println(result)
}
