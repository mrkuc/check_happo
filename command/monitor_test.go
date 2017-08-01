package command

import (
	"testing"

	"github.com/heartbeatsjp/happo-agent/lib"

	"github.com/stretchr/testify/assert"
)

const REQUEST_JSON_STR = "{\"apikey\":\"\",\"plugin_name\":\"plugin\",\"plugin_option\":\"option\"}"
const RESPONSE_JSON_STR = "{\"return_value\":0,\"message\":\"test\"}"

var RESPONSE_DATA = lib.MonitorResponse{ReturnValue: 0, Message: "test"}

func TestGetMonitorJSON1(t *testing.T) {
	ret, err := getMonitorJSON("plugin", "option")
	json_str := string(ret)
	assert.EqualValues(t, json_str, REQUEST_JSON_STR)
	assert.Nil(t, err)
}

func TestParseMonitorJSON1(t *testing.T) {
	ret, err := parseMonitorJSON(RESPONSE_JSON_STR)
	assert.EqualValues(t, ret, RESPONSE_DATA)
	assert.Nil(t, err)
}
