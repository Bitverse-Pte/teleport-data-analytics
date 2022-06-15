package tools

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

func HttpGet(url string, header map[string]string) ([]byte, error) {
	return sendHttpRequest("GET", url, nil, header)
}

func sendHttpRequest(method, url string, body io.Reader, header map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return []byte{}, err
	}
	//设置请求头
	setHeaders(req, &header)
	cli := http.Client{
		Timeout: 45 * time.Second,
	}
	resp, err := cli.Do(req)
	if err != nil {
		return []byte{}, err
	}
	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return out, err
	}
	if resp.StatusCode > 299 || resp.StatusCode < 200 {
		return out, errors.New("http response status:" + strconv.Itoa(resp.StatusCode) + " ,resp:" + string(out))
	}
	return out, nil
}

func setHeaders(req *http.Request, header *map[string]string) {
	if header == nil {
		req.Header.Add("Content-Type", "application/json")
		return
	}
	xh := *header
	if xh["Content-Type"] == "" {
		req.Header.Add("Content-Type", "application/json")
	}
	for key, value := range xh {
		req.Header.Add(key, value)
	}
}

func Post(url string, body io.Reader, header map[string]string) ([]byte, error) {
	return sendHttpRequest("POST", url, body, header)
}

// post请求，并且json unmarshal
func PostAndUnmarshal(url string, body io.Reader, header map[string]string, respPointer interface{}) error {
	byt, err := Post(url, body, header)
	if err != nil {
		return err
	}
	return json.Unmarshal(byt, respPointer)
}
