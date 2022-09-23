//
//	@Description 测试http
//	@return
//  @author hind3ight
//  @createdtime
//  @updatedtime

package internal

import (
	"encoding/json"
	"fmt"
	"gitee.com/wm-data/go-library/surfer"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	baseUrl      = `http://192.168.32.97:3111/goim/push`
	broadcastUrl = `/all?operation=%v`
	roomUrl      = `/room?operation=%v&type=live&room=%v`
	midUrl       = `/mids?operation=%v&mids=%v`

	operation = 1003
	room      = 1000
	mids      = 123
)

func SendMsgByHttp(sendType int, data []byte) {

	url := baseUrl
	switch sendType {
	case 1:
		url += fmt.Sprintf(roomUrl, operation, room)
	case 2:
		url += fmt.Sprintf(midUrl, operation, mids)
	case 3:
		url += fmt.Sprintf(broadcastUrl, operation)
	}

	resp, err := httpRequest2(url, data)
	if err != nil {
		fmt.Printf("发送请求错误,err:%s\n", err)
		return
	}

	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)
	type res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	var resData res
	json.Unmarshal(b, resData)
	//fmt.Println(string(b))
}

func httpRequest(url string, data []byte) (*http.Response, error) {
	header := http.Header{}
	header.Set("Content-Type", "application/json")

	rowData := string(data)

	conf := &surfer.DefaultRequest{
		Url:         url,
		Method:      "POST-J",
		PostData:    rowData,
		DialTimeout: time.Second * 15,
		TryTimes:    2,
		Header:      header,
	}
	return surfer.Download(conf)
}

func httpRequest2(url string, data []byte) (*http.Response, error) {
	client := &http.Client{}
	var msg = strings.NewReader(string(data))
	req, err := http.NewRequest("POST", url, msg)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Apipost client Runtime/+https://www.apipost.cn/")
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}
