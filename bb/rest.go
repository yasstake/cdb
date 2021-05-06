package bb

import (
	"cdb/trans"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type RestResponse struct {
	Code    json.Number     `json:"ret_code"`
	Message string          `json:"ret_msg"`
	ExtCode string          `json:"ext_code"`
	Result  json.RawMessage `json:"result"`
	Time    json.Number     `json:"time_now"`
}

func RestRequest(url string) (result string, time_now time.Time, err error) {
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Println(err)
		return result, time_now, err
	}

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return result, time_now, err
	}

	defer resp.Body.Close()
	byteArray, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
		return result, time_now, err
	}

	var response RestResponse
	err = json.Unmarshal(byteArray, &response)
	if err != nil {
		log.Println(err)
		return result, time_now, err
	}

	code, err := response.Code.Int64()

	if err != nil || code != 0 {
		return result, time_now, err
	}

	t, err := response.Time.Float64()
	if err != nil {
		return result, time_now, err
	}
	time_now = trans.DateTime(int64(float64(t) * float64(time.Second)))

	return string(response.Result), time_now, err
}
