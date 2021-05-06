package bb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

const LIQUID_TEST = `{"ret_code":0,"ret_msg":"OK","ext_code":"","ext_info":"","result":[{"id":7438589,"qty":738,"side":"Buy","time":1620172528603,"symbol":"BTCUSD","price":52956.5},{"id":7438590,"qty":14180,"side":"Buy","time":1620172528604,"symbol":"BTCUSD","price":52921.5},{"id":7438599,"qty":30079,"side":"Buy","time":1620172528608,"symbol":"BTCUSD","price":52934.5},{"id":7438600,"qty":100,"side":"Buy","time":1620172528604,"symbol":"BTCUSD","price":52929.5},{"id":7438601,"qty":3572,"side":"Buy","time":1620172528608,"symbol":"BTCUSD","price":52925},{"id":7438606,"qty":2563,"side":"Buy","time":1620172528610,"symbol":"BTCUSD","price":52929.5},{"id":7438607,"qty":3205,"side":"Buy","time":1620172529580,"symbol":"BTCUSD","price":52912},{"id":7438608,"qty":10245,"side":"Buy","time":1620172528610,"symbol":"BTCUSD","price":52935.5},{"id":7438609,"qty":10476,"side":"Buy","time":1620172529580,"symbol":"BTCUSD","price":52911.5},{"id":7438610,"qty":1600,"side":"Buy","time":1620172529580,"symbol":"BTCUSD","price":52910},{"id":7438615,"qty":478,"side":"Buy","time":1620172528611,"symbol":"BTCUSD","price":52951.5},{"id":7438618,"qty":49540,"side":"Buy","time":1620172528616,"symbol":"BTCUSD","price":52935},{"id":7438619,"qty":70,"side":"Buy","time":1620172528618,"symbol":"BTCUSD","price":52957},{"id":7438626,"qty":269,"side":"Buy","time":1620172528618,"symbol":"BTCUSD","price":52925},{"id":7438637,"qty":48680,"side":"Buy","time":1620172528612,"symbol":"BTCUSD","price":52930},{"id":7438638,"qty":14000,"side":"Buy","time":1620172528627,"symbol":"BTCUSD","price":52924},{"id":7438641,"qty":26462,"side":"Buy","time":1620172528612,"symbol":"BTCUSD","price":52929.5},{"id":7438653,"qty":443,"side":"Buy","time":1620172528613,"symbol":"BTCUSD","price":52937},{"id":7438666,"qty":10,"side":"Buy","time":1620172529592,"symbol":"BTCUSD","price":52913}]}`

const LIQUID_TEST2 = `{"ret_code":0,"ret_msg":"OK","ext_code":"","ext_info":"","result":[{"id":7438589,"qty":738,"side":"Buy","time":1620172528603,"symbol":"BTCUSD","price":52956.5}]}`
const LIQUID_TEST3 = `{"ret_code":0,"ret_msg":"OK","ext_code":"","ext_info":"","result":[{"id":7438589,"qty":738,"side":"Buy","time":1620172528603,"symbol":"BTCUSD","price":52956.5}]}`

func TestParseMessage(t *testing.T) {
	var response RestResponse

	err := json.Unmarshal([]byte(LIQUID_TEST), &response)

	if err != nil {
		fmt.Println(err)
		t.Error()
	}

	var liquid []LiquidRec
	err = json.Unmarshal([]byte(response.Result), &liquid)

	if err != nil {
		fmt.Println(err)
		t.Error()
	}

	fmt.Println(liquid)
}

// Test Bybit API(Liquidation )
func TestLiquidRest(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.bybit.com/v2/public/liq-records?symbol=BTCUSD", nil)

	client := new(http.Client)
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(byteArray))
	fmt.Println(LiquidMessageStr(string(byteArray)))
}

func TestLiquidString(t *testing.T) {
	s := LiquidMessageStr(LIQUID_TEST)
	fmt.Println(s)
}

func TestLiquidLoop(t *testing.T) {
	s, time, e := LiquidRequest(0)

	if e != nil {
		fmt.Println(s)
	}

	fmt.Println(time.UTC().String(), time.UTC().UnixNano())
	fmt.Println(LiquidMessageStr(s))
}
