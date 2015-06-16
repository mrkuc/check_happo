package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
	"github.com/heartbeatsjp/check_happo/command"
	"github.com/heartbeatsjp/happo-lib"
)

var GlobalFlags = []cli.Flag{}

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
				Value: happo_agent.DEFAULT_AGENT_PORT,
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
		},
	},

	{
		Name:   "check_happo",
		Usage:  "",
		Action: command.CmdTest,
		Flags:  []cli.Flag{},
	},
}

func CommandNotFound(c *cli.Context, command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command. See '%s --help'.", c.App.Name, command, c.App.Name, c.App.Name)
	os.Exit(2)
}
