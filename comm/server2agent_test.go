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
const PORT = lib.DEFAULT_AGENT_PORT
const METHOD = "TEST"
const JSON = "{\"test\":100}"

func TestGetProxyJSON1(t *testing.T) {
	PROXY_HOSTS := []string{"192.168.0.1"}
	var json_data lib.ProxyRequest
	PROXY_REQUEST := lib.ProxyRequest{
		ProxyHostPort: []string{fmt.Sprintf("%s:%d", HOST, PORT)},
		RequestType:   METHOD,
		RequestJSON:   ([]byte)(JSON),
	}

	json_str, agent_host, agent_port, err := GetProxyJSON(PROXY_HOSTS, HOST, PORT, "TEST", ([]byte(JSON)))
	assert.Nil(t, err)
	assert.EqualValues(t, agent_host, "192.168.0.1")
	assert.EqualValues(t, agent_port, lib.DEFAULT_AGENT_PORT)

	json.Unmarshal(json_str, &json_data)
	assert.EqualValues(t, json_data, PROXY_REQUEST)
}

func TestGetProxyJSON2(t *testing.T) {
	PROXY_HOSTS := []string{"192.168.0.1", "172.16.0.1"}
	var json_data lib.ProxyRequest
	PROXY_REQUEST := lib.ProxyRequest{
		ProxyHostPort: []string{"172.16.0.1", fmt.Sprintf("%s:%d", HOST, PORT)},
		RequestType:   METHOD,
		RequestJSON:   ([]byte)(JSON),
	}

	json_str, agent_host, agent_port, err := GetProxyJSON(PROXY_HOSTS, HOST, PORT, "TEST", ([]byte(JSON)))
	assert.Nil(t, err)
	assert.EqualValues(t, agent_host, "192.168.0.1")
	assert.EqualValues(t, agent_port, lib.DEFAULT_AGENT_PORT)

	json.Unmarshal(json_str, &json_data)
	assert.EqualValues(t, json_data, PROXY_REQUEST)
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

	json_str, err := PostToAgent(host, port, METHOD, ([]byte(JSON)))
	assert.NotNil(t, json_str)
	assert.Nil(t, err)
}

func TestPostToAgent2(t *testing.T) {
	json_str, err := PostToAgent("localhost", 12345, METHOD, ([]byte(JSON)))
	assert.EqualValues(t, json_str, "")
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

	json_str, err := PostToAgent(host, port, METHOD, ([]byte(JSON)))
	assert.EqualValues(t, json_str, "{}\n")
	assert.NotNil(t, err)
}
