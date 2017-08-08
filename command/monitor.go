package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/heartbeatsjp/check_happo/comm"
	"github.com/heartbeatsjp/happo-agent/halib"
)

// CmdMonitor implements `monitor` subcommand
func CmdMonitor(c *cli.Context) {
	var agentHost string
	var agentPort int
	var requestType string
	var monitorJSONStr []byte
	var jsonStr []byte
	var timeout int

	monitorJSONStr, err := getMonitorJSON(c.String("plugin_name"), c.String("plugin_option"))
	if err != nil {
		log.Print(err)
		os.Exit(halib.MonitorUnknown)
	}

	if len(c.StringSlice("proxy")) >= 1 {
		requestType = "proxy"
		jsonStr, agentHost, agentPort, err = comm.GetProxyJSON(c.StringSlice("proxy"), c.String("host"), c.Int("port"), "monitor", monitorJSONStr)
		if err != nil {
			log.Print(err)
			os.Exit(halib.MonitorUnknown)
		}
	} else {
		requestType = "monitor"
		agentHost = c.String("host")
		agentPort = c.Int("port")
		jsonStr = monitorJSONStr
	}

	timeout = c.Int("timeout")
	if timeout > 0 {
		comm.SetHTTPClientTimeout(timeout)
	}

	res, err := comm.PostToAgent(agentHost, agentPort, requestType, jsonStr)
	if err != nil {
		log.Print(err)
		os.Exit(halib.MonitorError)
	}
	monitorResult, err := parseMonitorJSON(res)
	if err != nil {
		log.Print(err)
		os.Exit(halib.MonitorError)
	}

	fmt.Print(monitorResult.Message)
	os.Exit(monitorResult.ReturnValue)
}

func getMonitorJSON(pluginName string, pluginOption string) ([]byte, error) {
	monitorRequest := halib.MonitorRequest{
		APIKey:       "",
		PluginName:   pluginName,
		PluginOption: pluginOption,
	}
	data, err := json.Marshal(monitorRequest)

	return data, err
}

func parseMonitorJSON(jsonStr string) (halib.MonitorResponse, error) {
	var m halib.MonitorResponse
	err := json.Unmarshal([]byte(jsonStr), &m)
	return m, err
}
