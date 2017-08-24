package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/heartbeatsjp/check_happo/command"
	"github.com/heartbeatsjp/happo-agent/halib"
)

// GlobalFlags are global level options
var GlobalFlags = []cli.Flag{}

// Commands is list of subcommand
var Commands = []cli.Command{

	{
		Name:   "monitor",
		Usage:  "",
		Action: command.CmdMonitor,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "host, H",
				Usage: "hostname or IP address",
			},
			cli.IntFlag{
				Name:  "port, P",
				Value: halib.DefaultAgentPort,
				Usage: "Port number",
			},
			cli.StringSliceFlag{
				Name:  "proxy, X",
				Value: &cli.StringSlice{},
				Usage: "Proxy hostname[:port] (You can multiple define.)",
			},
			cli.StringFlag{
				Name:  "plugin_name, p",
				Usage: "Plugin Name",
			},
			cli.StringFlag{
				Name:  "plugin_option, o",
				Usage: "Plugin Option",
			},
			cli.StringFlag{
				Name:  "timeout, t",
				Usage: "Connect Timeout",
			},
			cli.BoolFlag{
				Name:  "verbose, v",
				Usage: "verbose output",
			},
		},
	},

	{
		Name:   "check_happo",
		Usage:  "",
		Action: command.CmdTest,
		Flags:  []cli.Flag{},
	},
}

// CommandNotFound implements action when subcommand not found
func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
