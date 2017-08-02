package comm

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
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
var tlsConfig = &tls.Config{
	InsecureSkipVerify: true,
	MinVersion:         tls.VersionTLS12,
}
var tr = &http.Transport{
	TLSClientConfig: tlsConfig,
}
var _httpClient = &http.Client{Transport: tr, Timeout: 60 * time.Second} //default timeout 60 sec

// GetProxyJSON returns Marshal-ed ProxyRequest in []byte
func GetProxyJSON(proxyHosts []string, host string, port int, requestType string, proxyJSONStr []byte) ([]byte, string, int, error) {
	var agentHost string
	var agentPort int
	var err error

	// Step 1
	agentHostport := strings.Split(proxyHosts[0], ":")
	agentHost = agentHostport[0]
	if len(agentHostport) == 2 {
		agentPort, err = strconv.Atoi(agentHostport[1])
		if err != nil {
			return nil, "", 0, err
		}
	} else {
		agentPort = lib.DefaultAgentPort
	}

	// Step 2 or later
	proxyHosts = proxyHosts[1:]
	proxyHosts = append(proxyHosts, fmt.Sprintf("%s:%d", host, port))

	proxyRequest := lib.ProxyRequest{
		ProxyHostPort: proxyHosts,
		RequestType:   requestType,
		RequestJSON:   proxyJSONStr,
	}
	jsonData, _ := json.Marshal(proxyRequest)

	return jsonData, agentHost, agentPort, nil
}

// PostToAgent do HTTPS request and returns result
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
		return string(body[:]), fmt.Errorf("HTTP status error: %d", resp.StatusCode)
	}
	return string(body[:]), nil
}

// SetHTTPClientTimeout set timeout of httpClient
func SetHTTPClientTimeout(timeout int) {
	_httpClient.Timeout = time.Duration(timeout) * time.Second
}
