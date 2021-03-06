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

	"github.com/heartbeatsjp/check_happo/util"
	"github.com/heartbeatsjp/happo-agent/halib"
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
		agentPort = halib.DefaultAgentPort
	}

	// Step 2 or later
	proxyHosts = proxyHosts[1:]
	proxyHosts = append(proxyHosts, fmt.Sprintf("%s:%d", host, port))

	proxyRequest := halib.ProxyRequest{
		ProxyHostPort: proxyHosts,
		RequestType:   requestType,
		RequestJSON:   proxyJSONStr,
	}
	jsonData, _ := json.Marshal(proxyRequest)

	return jsonData, agentHost, agentPort, nil
}

// InternalRuntimeError is check_happo internal runtime error
type InternalRuntimeError struct {
	OriginalError error
}

func (e InternalRuntimeError) Error() string {
	return e.OriginalError.Error()
}

// PostToAgent do HTTPS request and returns HTTP response body, HTTP response status code, error
func PostToAgent(host string, port int, method string, jsonData []byte) (string, int, error) {
	uri := fmt.Sprintf("https://%s:%d/%s", host, port, method)
	log := util.Logger()
	log.Debug("Request: ", uri)

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Connection", "close")
	log.Debug("Request struct:\n", util.DumpStruct(req))

	resp, err := _httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}

	log.Debug("Response struct:\n", util.DumpStruct(resp))
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", 0, &InternalRuntimeError{OriginalError: err}
	}

	if resp.StatusCode != http.StatusOK {
		return string(body[:]), resp.StatusCode, fmt.Errorf("happo-agent returns %d", resp.StatusCode)
	}

	return string(body[:]), resp.StatusCode, nil
}

// SetHTTPClientTimeout set timeout of httpClient
func SetHTTPClientTimeout(timeout int) {
	_httpClient.Timeout = time.Duration(timeout) * time.Second
}
