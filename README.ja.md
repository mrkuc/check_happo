# `check_happo` - `happo-agent` helper

[![wercker status](https://app.wercker.com/status/29388c3fe43b1bb6cfb75b828a70f104/s/ "wercker status")](https://app.wercker.com/project/byKey/29388c3fe43b1bb6cfb75b828a70f104)

## Description

ネットワーク越しにサーバを監視できる [`happo-agent`](https://github.com/heartbeatsjp/happo-agent)に対してアクセスするためのヘルパープログラムです。


## Usage

### Requires

動作確認を行っているのは、以下のOSです。

- CentOS 6.6 (x86_64)以上
- CentOS 7 (x86_64)以上

恐らくx86_64アーキテクチャの各種Linuxディストーションで動作するはずです。

ビルドする際、go 1.5以上が必要となります。

### コマンドライン

エージェントに対して、監視リクエストを送信し、結果を受信します。

```
/path/to/check_happo monitor -H [ホスト名] -p [プラグイン名] -o [オプション]
```

例えば、`check_procs`を行うには次の通りにします。

```
$ check_happo monitor -H [HOSTNAME] -p check_procs -o '-w 100 -c 200'
2015/09/25 18:18:31 Request: https://HOSTNAME:6777/monitor
PROCS OK: 83 processes
```

踏み台サーバを越えてアクセスするには、`-P`オプションを利用します。
```
/path/to/check_happo monitor -X [プロキシサーバ名] -H [ホスト名] -p [プラグイン名] -o [オプション]
```

プロキシ経由でアクセスした結果は、直接アクセスする場合と変わりありません。

```
$ check_happo monitor -X 192.0.2.1 -H 198.51.100.1 -p check_procs -o '-w 100 -c 200'
2015/09/25 18:18:31 Request: https://192.0.2.1:6777/proxy
PROCS OK: 83 processes
```

詳細は `./check_happo --help` を実行し、確認してください。


## Install

To install, use `go get`:

```bash
$ sudo yum install epel-release
$ sudo yum install nagios-plugins-all
$ go get github.com/heartbeatsjp/check_happo
```

Nagiosのcommandsファイルには以下の通りの設定を行ってください。

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

## License

Copyright 2016 HEARTBEATS Corporation.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
