package common

import (
	"fmt"
	"net/url"
	"os"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetCurrentPath() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return strings.Replace(dir, "\\", "/", -1), nil
}


func BuildUrlParam(sUrl string, params map[string]interface{}) string {
	var buf strings.Builder
	buf.WriteString(sUrl)
	buf.WriteByte('?')

	for k, v := range params {
		buf.WriteString(url.QueryEscape(k))
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(fmt.Sprint(v)))
		buf.WriteByte('&')
	}

	return buf.String()
}
