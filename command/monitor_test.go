package command

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	"github.com/heartbeatsjp/happo-agent/halib"

	"github.com/stretchr/testify/assert"
)

const RequestJSONStr = "{\"apikey\":\"\",\"plugin_name\":\"plugin\",\"plugin_option\":\"option\"}"
const ResponseJSONStr = "{\"return_value\":0,\"message\":\"test\"}"

var ResponseData = halib.MonitorResponse{ReturnValue: 0, Message: "test"}

func TestGetMonitorJSON1(t *testing.T) {
	ret, err := getMonitorJSON("plugin", "option")
	jsonStr := string(ret)
	assert.EqualValues(t, RequestJSONStr, jsonStr)
	assert.Nil(t, err)
}

func TestParseMonitorJSON1(t *testing.T) {
	ret, err := parseMonitorJSON(ResponseJSONStr)
	assert.EqualValues(t, ResponseData, ret)
	assert.Nil(t, err)
}

func TestCmdMonitorDirect(t *testing.T) {
	patterns := []struct {
		responseBody string
		responseCode int
		expectOut    string
		expectExit   int
	}{
		{responseBody: `{"return_value":0,"message":"PROCS OK: 100 processes\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "PROCS OK: 100 processes\n",
			expectExit:   halib.MonitorOK},
		{responseBody: `{"return_value":1,"message":"PROCS WARNING: 168 processes\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "PROCS WARNING: 168 processes\n",
			expectExit:   halib.MonitorWarning},
		{responseBody: `{"return_value":2,"message":"PROCS CRITICAL: 170 processes\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "PROCS CRITICAL: 170 processes\n",
			expectExit:   halib.MonitorError},
		{responseBody: `{"return_value":3,"message":"UNKNOWN ERROR OCCURED\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "UNKNOWN ERROR OCCURED\n",
			expectExit:   halib.MonitorUnknown},
		{responseBody: ``,
			responseCode: http.StatusBadRequest,
			expectOut:    `ERROR(happo-agent): happo-agent returns 400`,
			expectExit:   halib.MonitorError},
		{responseBody: `{"return_value":124,"message":""}`,
			responseCode: http.StatusInternalServerError,
			expectOut:    `ERROR(happo-agent): happo-agent returns 500`,
			expectExit:   halib.MonitorError},
		{responseBody: `{"return_value":124,"message":""}`,
			responseCode: http.StatusServiceUnavailable,
			expectOut:    `ERROR(happo-agent): happo-agent returns 503`,
			expectExit:   halib.MonitorError},
		{responseBody: ``,
			responseCode: http.StatusGatewayTimeout,
			expectOut:    "UNKNOWN(happo-agent): happo-agent returns 504",
			expectExit:   halib.MonitorUnknown},
	}
	for _, pattern := range patterns {
		out, code := testCmdMonitorDirect(pattern.responseBody, pattern.responseCode)
		assert.Equal(t, pattern.expectOut, out)
		assert.Equal(t, pattern.expectExit, code)
	}
}

func testCmdMonitorDirect(responseBody string, responseCode int) (string, int) {

	ts := httptest.NewTLSServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, responseBody, responseCode)
			}))
	defer ts.Close()
	re, _ := regexp.Compile("([a-z]+)://([A-Za-z0-9.]+):([0-9]+)(.*)")
	found := re.FindStringSubmatch(ts.URL)
	host := found[2]
	port, _ := strconv.Atoi(found[3])

	out, code := cmdMonitor(false, "", "", make([]string, 0), host, port, 3)

	return out, code
}

func TestCmdMonitorProxy(t *testing.T) {
	patterns := []struct {
		responseBody string
		responseCode int
		expectOut    string
		expectExit   int
	}{
		{responseBody: `{"return_value":0,"message":"PROCS OK: 100 processes\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "PROCS OK: 100 processes\n",
			expectExit:   halib.MonitorOK},
		{responseBody: `{"return_value":1,"message":"PROCS WARNING: 168 processes\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "PROCS WARNING: 168 processes\n",
			expectExit:   halib.MonitorWarning},
		{responseBody: `{"return_value":2,"message":"PROCS CRITICAL: 170 processes\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "PROCS CRITICAL: 170 processes\n",
			expectExit:   halib.MonitorError},
		{responseBody: `{"return_value":3,"message":"UNKNOWN ERROR OCCURED\n"}`,
			responseCode: http.StatusOK,
			expectOut:    "UNKNOWN ERROR OCCURED\n",
			expectExit:   halib.MonitorUnknown},
		{responseBody: ``,
			responseCode: http.StatusBadRequest,
			expectOut:    `ERROR(happo-agent): happo-agent returns 400`,
			expectExit:   halib.MonitorError},
		{responseBody: `{"return_value":124,"message":""}`,
			responseCode: http.StatusInternalServerError,
			expectOut:    `ERROR(happo-agent): happo-agent returns 500`,
			expectExit:   halib.MonitorError},
		{responseBody: `{"return_value":124,"message":""}`,
			responseCode: http.StatusServiceUnavailable,
			expectOut:    `ERROR(happo-agent): happo-agent returns 503`,
			expectExit:   halib.MonitorError},
		{responseBody: ``,
			responseCode: http.StatusGatewayTimeout,
			expectOut:    "UNKNOWN(happo-agent): happo-agent returns 504",
			expectExit:   halib.MonitorUnknown},
	}
	for _, pattern := range patterns {
		out, code := testCmdMonitorDirect(pattern.responseBody, pattern.responseCode)
		assert.Equal(t, pattern.expectOut, out)
		assert.Equal(t, pattern.expectExit, code)
	}
}

func testCmdMonitorProxy(responseBody string, responseCode int) (string, int) {

	ts := httptest.NewTLSServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				http.Error(w, responseBody, responseCode)
			}))
	defer ts.Close()
	re, _ := regexp.Compile("([a-z]+)://([A-Za-z0-9.]+):([0-9]+)(.*)")
	found := re.FindStringSubmatch(ts.URL)
	host := found[2]
	port, _ := strconv.Atoi(found[3])

	proxy := []string{fmt.Sprintf("%s:%v", host, port)}
	out, code := cmdMonitor(false, "", "", proxy, host, port, 3)

	return out, code
}
