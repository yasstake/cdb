package bb

import (
	"fmt"
	"testing"
)

const TICKER_MESSAGE = `[{"symbol":"BTCUSD","bid_price":"57615.5","ask_price":"57616","last_price":"57616.00","last_tick_direction":"ZeroPlusTick","prev_price_24h":"57269.50","price_24h_pcnt":"0.00605","high_price_24h":"57729.00","low_price_24h":"55311.00","prev_price_1h":"57182.00","price_1h_pcnt":"0.007589","mark_price":"57602.40","index_price":"57596.67","open_interest":2081684903,"open_value":"14616.80","total_turnover":"77137353.10","turnover_24h":"90800.77","total_volume":1530791399383,"volume_24h":5126153288,"funding_rate":"0.0001","predicted_funding_rate":"0.0001","next_funding_time":"2021-05-08T00:00:00Z","countdown_hour":8,"delivery_fee_rate":"0","predicted_delivery_price":"0.00","delivery_time":""}]`

func TestTickerRequest(t *testing.T) {
	body, time, err := TickerRequest()

	fmt.Println(body)
	fmt.Println(time)
	fmt.Println(err)
}

func TestTickerMessage(t *testing.T) {
	ticker, err := TickerMessage(TICKER_MESSAGE)

	fmt.Println(ticker)
	fmt.Println(err)
}
