package gtools

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"runtime"
	"strings"
	"time"
)

const (
	HttpMethodGet  string = "GET"
	HttpMethodPOST string = "POST"
)

type HttpTimeout struct {
	DialTimeout           int
	DialKeepAlive         int
	TLSHandshakeTimeout   int
	ResponseHeaderTimeout int
	ExpectContinueTimeout int
	Timeout               int
}

func DefaultHttpTimeout() HttpTimeout {
	return SetHttpTimeout(60, 75, 75, 90, 60, 90)
}

func SetHttpTimeout(dialTimeout, dialKeepAlive, tlsHandshakeTimeout, responseHeaderTimeout, expectContinueTimeout, timeout int) HttpTimeout {
	var tt HttpTimeout

	tt.DialTimeout = dialTimeout
	tt.DialKeepAlive = dialKeepAlive
	tt.TLSHandshakeTimeout = tlsHandshakeTimeout
	tt.ResponseHeaderTimeout = responseHeaderTimeout
	tt.ExpectContinueTimeout = expectContinueTimeout
	tt.Timeout = timeout

	return tt
}

func httClientWithTimeout(timeoutConf HttpTimeout) (client *http.Client) {
	var netTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   time.Second * time.Duration(timeoutConf.DialTimeout),
			KeepAlive: time.Second * time.Duration(timeoutConf.DialKeepAlive),
		}).DialContext,
		TLSHandshakeTimeout:   time.Second * time.Duration(timeoutConf.TLSHandshakeTimeout),
		ResponseHeaderTimeout: time.Second * time.Duration(timeoutConf.ResponseHeaderTimeout),
		ExpectContinueTimeout: time.Second * time.Duration(timeoutConf.ExpectContinueTimeout),
	}
	client = &http.Client{
		Timeout:   time.Second * time.Duration(timeoutConf.Timeout),
		Transport: netTransport,
	}
	return
}

// SimpleHttpClient
// Deprecated: please use more advanced go-resty instead
// see: https://github.com/go-resty/resty
func SimpleHttpClient(reqMethod string, reqUrl string, reqHeaders map[string]string, reqBody string, timeoutConf HttpTimeout) (statusCode int, respHeader http.Header, respBody []byte, err error) {
	req, err := http.NewRequest(reqMethod, reqUrl, strings.NewReader(reqBody))
	if err != nil {
		log.Printf("[SimpleHttpClient] http.NewRequest fail, reqUrl: %s", reqUrl)
		return
	}

	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}

	client := httClientWithTimeout(timeoutConf)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("[SimpleHttpClient] do request fail, reqUrl:", reqUrl, ", err:", err)
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	respBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("[SimpleHttpClient] read request fail, reqUrl:", reqUrl, ", err:", err)
		return
	}

	respHeader = resp.Header
	statusCode = resp.StatusCode

	return
}

func BasicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func DefaultFormReqHeaders() map[string]string {
	return map[string]string{
		"Connection": "keep-alive",
		//"Content-Type": "application/x-www-form-urlencoded",
		"Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
		"User-Agent":   runtime.Version(),
	}
}

func DefaultJsonReqHeaders() map[string]string {
	return map[string]string{
		"Connection":   "keep-alive",
		"Content-Type": "application/json",
		"User-Agent":   runtime.Version(),
	}
}

func ParseCookie(cookie, key string) (value string, err error) {
	cookieBox := strings.Split(cookie, "; ")
	if len(cookieBox) == 0 {
		err = fmt.Errorf(`empty origin data, cooke: %s`, cookie)
		return
	}

	for _, item := range cookieBox {
		itemBox := strings.Split(item, "=")
		if len(itemBox) != 2 {
			log.Printf("[ParseCookie] get unexpected data: %s", item)
			continue
		} else {
			if key == itemBox[0] {
				value = itemBox[1]
				break
			}
		}
	}

	return
}
