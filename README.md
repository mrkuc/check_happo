# `check_happo` - Nagios plugin for `happo-agent`

## Description

Call monitoring request helper for [`happo-agent`](https://github.com/heartbeatsjp/happo-agent).

### Requires

- Red Hat Enterprise Linux (RHEL) 6.x, 7.x
- CentOS 6.x, 7.x
- Ubuntu 12.04 or later

### Command

Request to `happo-agent`.

```
/path/to/check_happo monitor -H [HOSTNAME] -p [NAGIOS_PLUGINNAME] -o [PLUGIN_OPTION]
```

For example, execute `check_procs` is below。

```
$ check_happo monitor -H 127.0.0.1 -p check_procs -o '-w 100 -c 200'
2015/09/25 18:18:31 Request: https://127.0.0.1:6777/monitor
PROCS OK: 83 processes
```

If access server via bastion, require `-X` option.

```
/path/to/check_happo monitor -X [BASTION_SERVER] -H [HOSTNAME] -p [NAGIOS_PLUGINNAME] -o [PLUGIN_OPTION]
```

For example.

```
$ check_happo monitor -X 192.0.2.1 -H 198.51.100.1 -p check_procs -o '-w 100 -c 200'
2015/09/25 18:18:31 Request: https://192.0.2.1:6777/proxy
PROCS OK: 83 processes
```

For more information, please execute `./check_happo --help`.


## Install

To install, use `go get`:

```bash
$ sudo yum install epel-release
$ sudo yum install nagios-plugins-all
$ go get github.com/heartbeatsjp/check_happo
```

And define Nagios command:

```
define command {
        command_name    check_by_happo
        command_line    /path/to/check_happo monitor -H $ARG1$ -X $ARG2$ -p $ARG3$ -o $ARG4$
    }
```


## Contribution

1. Fork ([http://github.com/heartbeatsjp/check_happo/fork](http://github.com/heartbeatsjp/check_happo/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request


## Author

[Yuichiro Saito](https://github.com/koemu)
