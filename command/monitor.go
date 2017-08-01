package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/heartbeatsjp/check_happo/comm"
	"github.com/heartbeatsjp/happo-agent/lib"
)

func CmdMonitor(c *cli.Context) {
	var agent_host string
	var agent_port int
	var request_type string
	var monitor_jsonStr []byte
	var jsonStr []byte
	var timeout int

	monitor_jsonStr, err := getMonitorJSON(c.String("plugin_name"), c.String("plugin_option"))
	if err != nil {
		log.Print(err)
		os.Exit(lib.MonitorUnknown)
	}

	if len(c.StringSlice("proxy")) >= 1 {
		request_type = "proxy"
		jsonStr, agent_host, agent_port, err = comm.GetProxyJSON(c.StringSlice("proxy"), c.String("host"), c.Int("port"), "monitor", monitor_jsonStr)
		if err != nil {
			log.Print(err)
			os.Exit(lib.MonitorUnknown)
		}
	} else {
		request_type = "monitor"
		agent_host = c.String("host")
		agent_port = c.Int("port")
		jsonStr = monitor_jsonStr
	}

	timeout = c.Int("timeout")
	if timeout > 0 {
		comm.SetHTTPClientTimeout(timeout)
	}

	res, err := comm.PostToAgent(agent_host, agent_port, request_type, jsonStr)
	if err != nil {
		log.Print(err)
		os.Exit(lib.MonitorError)
	}
	monitor_result, err := parseMonitorJSON(res)
	if err != nil {
		log.Print(err)
		os.Exit(lib.MonitorError)
	}

	fmt.Print(monitor_result.Message)
	os.Exit(monitor_result.ReturnValue)
}

func getMonitorJSON(plugin_name string, plugin_option string) ([]byte, error) {
	monitor_request := lib.MonitorRequest{
		APIKey:       "",
		PluginName:   plugin_name,
		PluginOption: plugin_option,
	}
	data, err := json.Marshal(monitor_request)

	return data, err
}

func parseMonitorJSON(jsonStr string) (lib.MonitorResponse, error) {
	var m lib.MonitorResponse
	err := json.Unmarshal([]byte(jsonStr), &m)
	return m, err
}
