package command

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/heartbeatsjp/check_happo/comm"
	"github.com/heartbeatsjp/check_happo/util"
	"github.com/heartbeatsjp/happo-agent/halib"
)

const (
	localErrorMessageFormat    = "ERROR(check_happo): %s"
	localUnknownMessageFormat  = "UNKNOWN(check_happo): %s"
	remoteErrorMessageFormat   = "ERROR(happo-agent): %s"
	remoteUnknownMessageFormat = "UNKNOWN(happo-agent): %s"
)

// CmdMonitor implements `monitor` subcommand
func CmdMonitor(c *cli.Context) {

	out, code := cmdMonitor(c.Bool("verbose"),
		c.String("plugin_name"),
		c.String("plugin_option"),
		c.StringSlice("proxy"),
		c.String("host"),
		c.Int("port"),
		c.Int("timeout"),
	)
	fmt.Print(out)
	os.Exit(code)
}

func cmdMonitor(verbose bool, pluginName, pluginOption string, proxy []string, host string, port int, timeout int) (string, int) {
	var agentHost string
	var agentPort int
	var requestType string
	var monitorJSONStr []byte
	var jsonStr []byte

	if verbose {
		util.LoggerLevelDebug()
	}

	monitorJSONStr, err := getMonitorJSON(pluginName, pluginOption)
	if err != nil {
		return fmt.Sprintf(localErrorMessageFormat, err.Error()), halib.MonitorError
	}

	if len(proxy) >= 1 {
		requestType = "proxy"
		jsonStr, agentHost, agentPort, err = comm.GetProxyJSON(proxy, host, port, "monitor", monitorJSONStr)
		if err != nil {
			return fmt.Sprintf(localErrorMessageFormat, err.Error()), halib.MonitorError
		}
	} else {
		requestType = "monitor"
		agentHost = host
		agentPort = port
		jsonStr = monitorJSONStr
	}

	if timeout > 0 {
		comm.SetHTTPClientTimeout(timeout)
	}

	responseBody, responseStatusCode, err := comm.PostToAgent(agentHost, agentPort, requestType, jsonStr)

	if err != nil {
		if internalError, ok := err.(comm.InternalRuntimeError); ok {
			return fmt.Sprintf(localErrorMessageFormat, internalError.Error()), halib.MonitorError
		}
		// connection failed or responseStatusCode != 200
		if responseStatusCode == http.StatusGatewayTimeout {
			return fmt.Sprintf(remoteUnknownMessageFormat, err.Error()), halib.MonitorUnknown
		}
		msg := fmt.Sprintf("%s %s", err.Error(), responseBody)
		return fmt.Sprintf(remoteErrorMessageFormat, msg), halib.MonitorError
	}

	monitorResult, err := parseMonitorJSON(responseBody)
	if err != nil {
		return fmt.Sprintf(localErrorMessageFormat, err.Error()), halib.MonitorError
	}

	if responseStatusCode == http.StatusOK {
		return monitorResult.Message, monitorResult.ReturnValue
	}
	return fmt.Sprintf(remoteErrorMessageFormat, monitorResult.Message), halib.MonitorError
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
