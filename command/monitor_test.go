package command

import (
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
