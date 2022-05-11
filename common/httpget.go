package common



import (
	"bytes"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

/** headers 常用 Content-type
	"text/html"
    "application/json"
    "application/xml"
    "text/plain"
    "application/x-www-form-urlencoded"
    "multipart/form-data"
*/
func HttpDo(queryUrl, method string, params, headers map[string]string, data []byte, timeout int64, options map[string]interface{}) (body []byte, err error) {
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}
	//可选参数
	if options != nil {
		transport := &http.Transport{}
		//代理请求
		if proxyUrl, ok := options["proxy"]; ok && proxyUrl != "" {
			parsedProxyUrl, tmpErr := url.Parse(proxyUrl.(string))
			if tmpErr != nil {
				err = tmpErr
				return
			}
			transport.Proxy = http.ProxyURL(parsedProxyUrl)
		}
		if tlsConfig, ok := options["tls"]; ok {
			transport.TLSClientConfig = tlsConfig.(*tls.Config)
		}
		client.Transport = transport
	}
	req, err := http.NewRequest(method, queryUrl, bytes.NewBuffer(data))
	if err != nil {
		return
	}
	//add params
	q := req.URL.Query()
	if params != nil {
		for key, val := range params {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	resp, err := client.Do(req) //  默认的resp ,err :=  http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	return
}


func HttpGet(queryUrl string, params map[string]string, timeout int64) (body []byte, err error) {
	return HttpDo(queryUrl, "GET", params, map[string]string{}, []byte{}, timeout, map[string]interface{}{})
}
