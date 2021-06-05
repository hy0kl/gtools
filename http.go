package gtools

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path/filepath"
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

// 支持post原始多文件上传,同时携带表单数据
func MultipartClient(reqUrl string, queryString map[string]interface{}, reqHeaders map[string]string, files map[string]string, timeoutConf HttpTimeout) (httpStatusCode int, respHeader http.Header, originByte []byte, err error) {
	client := httClientWithTimeout(timeoutConf)

	// 创建一个缓冲区对象,后面的要上传的body都存在这个缓冲区里
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	if len(files) <= 0 {
		log.Println("no file for upload")
	}
	// name: 上传表单中字段名; localABS: 待上传文件路径
	for fname, filename := range files {
		// 创建第一个需要上传的文件,filepath.Base获取文件的名称
		var fileWriter io.Writer
		fileWriter, err = bodyWriter.CreateFormFile(fname, filepath.Base(filename))
		if err != nil {
			log.Println("the uploaded file does not exist. fname:", fname, ", filename:", filename, ", err:", err)
			return
		}

		// 打开文件
		var fd *os.File
		fd, err = os.Open(filename)
		if err != nil {
			log.Println("Can NOT open file. fname:", fname, ", filename:", filename, ", err:", err)
			return
		}

		// 把第文件流写入到缓冲区里去
		_, err = io.Copy(fileWriter, fd)
		_ = fd.Close()
		if err != nil {
			log.Println("Can NOT copy stream. fname:", fname, ", filename:", filename, ", err:", err)
			return
		}
	}

	// 写入附加字段必须在_,_=io.Copy(fileWriter,fd)后面
	// 写入常规k,v参数
	for k, v := range queryString {
		_ = bodyWriter.WriteField(k, Stringify(v))
	}

	// 获取请求Content-Type类型,后面有用
	contentType := bodyWriter.FormDataContentType()
	_ = bodyWriter.Close()

	// 创建一个post请求
	req, err := http.NewRequest("POST", reqUrl, nil)
	if err != nil {
		log.Println("http.NewRequest has wrong. reqUrl:", reqUrl, ", queryString:", queryString, ", reqHeaders:", reqHeaders, ", files:", files)
		return
	}

	for k, v := range reqHeaders {
		req.Header.Set(k, v)
	}
	req.Header.Set("Content-Type", contentType)
	// 转换类型
	req.Body = ioutil.NopCloser(bodyBuf)
	// 发送数据
	resp, err := client.Do(req)
	if err != nil {
		log.Println("client.Do has wrong. reqUrl:", reqUrl, ", queryString:", queryString, ", reqHeaders:", reqHeaders, ", files:", files, ", err:", err)
		return
	}

	//读取请求返回的数据
	originByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("ReadAll has wrong. reqUrl:", reqUrl, ", queryString:", queryString, ", reqHeaders:", reqHeaders, ", files:", files, ", err:", err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	httpStatusCode = resp.StatusCode
	respHeader = resp.Header

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
