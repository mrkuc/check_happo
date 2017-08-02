package comm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/heartbeatsjp/happo-agent/lib"
)

// --- Global Variables
// See http://golang.org/pkg/net/http/#Client
var tls_config = &tls.Config{
	InsecureSkipVerify: true,
	MinVersion:         tls.VersionTLS12,
}
var tr = &http.Transport{
	TLSClientConfig: tls_config,
}
var _httpClient = &http.Client{Transport: tr, Timeout: 60 * time.Second} //default timeout 60 sec

func GetProxyJSON(proxy_hosts []string, host string, port int, request_type string, proxy_jsonStr []byte) ([]byte, string, int, error) {
	var agent_host string
	var agent_port int
	var err error

	// Step 1
	agent_hostport := strings.Split(proxy_hosts[0], ":")
	agent_host = agent_hostport[0]
	if len(agent_hostport) == 2 {
		agent_port, err = strconv.Atoi(agent_hostport[1])
		if err != nil {
			return nil, "", 0, err
		}
	} else {
		agent_port = lib.DefaultAgentPort
	}

	// Step 2 or later
	proxy_hosts = proxy_hosts[1:]
	proxy_hosts = append(proxy_hosts, fmt.Sprintf("%s:%d", host, port))

	proxy_request := lib.ProxyRequest{
		ProxyHostPort: proxy_hosts,
		RequestType:   request_type,
		RequestJSON:   proxy_jsonStr,
	}
	jsonData, _ := json.Marshal(proxy_request)

	return jsonData, agent_host, agent_port, nil
}

func PostToAgent(host string, port int, method string, jsonData []byte) (string, error) {
	uri := fmt.Sprintf("https://%s:%d/%s", host, port, method)
	log := Logger()
	log.Debug("Request: ", uri)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := _httpClient.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return string(body[:]), errors.New(fmt.Sprintf("HTTP status error: %d", resp.StatusCode))
	}
	return string(body[:]), nil
}

func SetHTTPClientTimeout(timeout int) {
	_httpClient.Timeout = time.Duration(timeout) * time.Second
}
