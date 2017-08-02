package comm

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/heartbeatsjp/happo-agent/lib"
)

const HOST = "10.0.0.1"
const PORT = lib.DefaultAgentPort
const METHOD = "TEST"
const JSON = "{\"test\":100}"

func TestGetProxyJSON1(t *testing.T) {
	ProxyHosts := []string{"192.168.0.1"}
	var jsonData lib.ProxyRequest
	ProxyRequest := lib.ProxyRequest{
		ProxyHostPort: []string{fmt.Sprintf("%s:%d", HOST, PORT)},
		RequestType:   METHOD,
		RequestJSON:   ([]byte)(JSON),
	}

	jsonStr, agentHost, agentPort, err := GetProxyJSON(ProxyHosts, HOST, PORT, "TEST", ([]byte(JSON)))
	assert.Nil(t, err)
	assert.EqualValues(t, "192.168.0.1", agentHost)
	assert.EqualValues(t, lib.DefaultAgentPort, agentPort)

	json.Unmarshal(jsonStr, &jsonData)
	assert.EqualValues(t, ProxyRequest, jsonData)
}

func TestGetProxyJSON2(t *testing.T) {
	ProxyHosts := []string{"192.168.0.1", "172.16.0.1"}
	var jsonData lib.ProxyRequest
	ProxyRequest := lib.ProxyRequest{
		ProxyHostPort: []string{"172.16.0.1", fmt.Sprintf("%s:%d", HOST, PORT)},
		RequestType:   METHOD,
		RequestJSON:   ([]byte)(JSON),
	}

	jsonStr, agentHost, agentPort, err := GetProxyJSON(ProxyHosts, HOST, PORT, "TEST", ([]byte(JSON)))
	assert.Nil(t, err)
	assert.EqualValues(t, "192.168.0.1", agentHost)
	assert.EqualValues(t, lib.DefaultAgentPort, agentPort)

	json.Unmarshal(jsonStr, &jsonData)
	assert.EqualValues(t, ProxyRequest, jsonData)
}

func TestPostToAgent1(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintln(w, "{}")
			}))
	defer ts.Close()
	re, _ := regexp.Compile("([a-z]+)://([A-Za-z0-9.]+):([0-9]+)(.*)")
	found := re.FindStringSubmatch(ts.URL)
	host := found[2]
	port, _ := strconv.Atoi(found[3])

	jsonStr, err := PostToAgent(host, port, METHOD, ([]byte(JSON)))
	assert.NotNil(t, jsonStr)
	assert.Nil(t, err)
}

func TestPostToAgent2(t *testing.T) {
	jsonStr, err := PostToAgent("localhost", 12345, METHOD, ([]byte(JSON)))
	assert.EqualValues(t, "", jsonStr)
	assert.NotNil(t, err)
}

func TestPostToAgent3(t *testing.T) {
	ts := httptest.NewTLSServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, "{}", http.StatusNotAcceptable)
			}))
	defer ts.Close()
	re, _ := regexp.Compile("([a-z]+)://([A-Za-z0-9.]+):([0-9]+)(.*)")
	found := re.FindStringSubmatch(ts.URL)
	host := found[2]
	port, _ := strconv.Atoi(found[3])

	jsonStr, err := PostToAgent(host, port, METHOD, ([]byte(JSON)))
	assert.EqualValues(t, "{}\n", jsonStr)
	assert.NotNil(t, err)
}
