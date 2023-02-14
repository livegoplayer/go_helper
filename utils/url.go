package utils

import (
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func EncodeUrlWithoutSort(v url.Values, keys []string, withOutEncode bool) string {
	if v == nil {
		return ""
	}
	var buf string
	for _, k := range keys {
		vs, ok := v[k]
		if !ok {
			continue
		}
		keyEscaped := url.QueryEscape(k)
		if withOutEncode {
			keyEscaped = k
		}
		for _, vl := range vs {
			if len(buf) > 0 {
				buf += "&"
			}
			buf += keyEscaped
			buf += "="
			if withOutEncode {
				buf += vl
			} else {
				buf += url.QueryEscape(vl)
			}
		}
	}
	res := buf
	return res
}

func HttpGet(url string) (map[string]interface{}, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rep := JsonDecodeToMap(string(body))
	return rep, nil
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		conn.SetDeadline(time.Now().Add(rwTimeout))
		return conn, nil
	}
}

func HttpPost(url string, params interface{}, seconds int, header ...interface{}) (map[string]interface{}, error) {
	body := JsonEncode(params)

	req, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	if len(header) > 0 {
		headers := header[0].(map[string]string)
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	connectTimeout := time.Duration(seconds) * time.Second
	readWriteTimeout := time.Duration(seconds) * time.Second
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			Dial:            TimeoutDialer(connectTimeout, readWriteTimeout),
		},
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	content, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	rep := JsonDecodeToMap(string(content))
	return rep, nil
}
