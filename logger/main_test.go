package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

const TRADE_MESSAGE = `{"topic":"trade.BTCUSD","data":[{"trade_time_ms":1619392252147,"timestamp":"2021-04-25T23:10:52.000Z","symbol":"BTCUSD","side":"Sell","size":11,"price":48590,"tick_direction":"MinusTick","trade_id":"54055af6-d73d-5eb3-a8e7-13c87\
9f74ae6","cross_seq":6166836535},{"trade_time_ms":1619392252148,"timestamp":"2021-04-25T23:10:52.000Z","symbol":"B\
TCUSD","side":"Sell","size":4,"price":48590,"tick_direction":"ZeroMinusTick","trade_id":"2c769500-feb4-52bb-9c9b-9\
3d0a0b43c29","cross_seq":6166836535},{"trade_time_ms":1619392252149,"timestamp":"2021-04-25T23:10:52.000Z","symbol\
":"BTCUSD","side":"Sell","size":3,"price":48590,"tick_direction":"ZeroMinusTick","trade_id":"b9f2158a-65fd-52a7-b9\
52-473143b5bf2f","cross_seq":6166836535},{"trade_time_ms":1619392252149,"timestamp":"2021-04-25T23:10:52.000Z","sy\
mbol":"BTCUSD","side":"Sell","size":3,"price":48590,"tick_direction":"ZeroMinusTick","trade_id":"77759ace-cd3e-59e\
4-99a9-c7f934fb68b0","cross_seq":6166836535},{"trade_time_ms":1619392252150,"timestamp":"2021-04-25T23:10:52.000Z"\
,"symbol":"BTCUSD","side":"Sell","size":1,"price":48590,"tick_direction":"ZeroMinusTick","trade_id":"6bc67978-a12d\
-5b8c-a72d-cdc9287b4603","cross_seq":6166836535},{"trade_time_ms":1619392252151,"timestamp":"2021-04-25T23:10:52.0\
00Z","symbol":"BTCUSD","side":"Sell","size":1,"price":48590,"tick_direction":"ZeroMinusTick","trade_id":"fbfa637f-\
83ac-5326-a292-f1035e930a24","cross_seq":6166836535},{"trade_time_ms":1619392252155,"timestamp":"2021-04-25T23:10:\
52.000Z","symbol":"BTCUSD","side":"Sell","size":1,"price":48590,"tick_direction":"ZeroMinusTick","trade_id":"b7a2e\
ece-3f52-5761-849a-2f58e43229f3","cross_seq":6166836535},{"trade_time_ms":1619392252158,"timestamp":"2021-04-25T23\
:10:52.000Z","symbol":"BTCUSD","side":"Sell","size":1,"price":48590,"tick_direction":"ZeroMinusTick","trade_id":"4\
79a6750-653d-56dc-8cbb-45cdbcc4307c","cross_seq":6166836535}]}`

func TestUnMarshall(t *testing.T) {
	var decoded Message

	var json_message = ORDER_BOOK_SNAP
	err := json.Unmarshal([]byte(json_message), &decoded)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---ORDERBOOK---", json_message)
	fmt.Println(decoded.Topic)
	fmt.Println(decoded.Time)
	fmt.Println(decoded.Sequence)
	fmt.Println(decoded.Type)

	json_message = ORDER_BOOK_DELTA
	err = json.Unmarshal([]byte(json_message), &decoded)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---SNAP---", json_message)
	fmt.Println(decoded)

	json_message = TRADE_MESSAGE
	err = json.Unmarshal([]byte(json_message), &decoded)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("---SNAP---", json_message)
	fmt.Println(decoded)

}
