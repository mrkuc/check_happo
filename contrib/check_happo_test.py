# coding: utf-8

"""
# Usage

```
python check_happo_test.py | sort -n
```

- may need about 1 minutes
- sometime fails. timing problem?
    - exit status is correct, but message is not same

# Requirements

- python 2.7
- run on CentOS 7 + happo-agent
    - happo-agent listen 0.0.0.0
    - 127.0.1.1 can reach to happo-agent
- 192.168.0.1 leads to timeout
"""

import subprocess
import shlex
import re
from multiprocessing import Process


def run_test(index, command_line, expect_out, expect_code):
    p = subprocess.Popen(shlex.split(command_line),
                         stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    out, err = p.communicate()
    if isinstance(expect_out, re._pattern_type):
        if re.search(expect_out, out.strip()) and p.returncode == expect_code:
            return u"%d OK: %s" % (index, command_line)
        return u"%d NG: %s\t(expect exit=%d)\t'%s'\t(actual exit=%d)\t'%s'" % (
            index,
            command_line,
            expect_code,
            expect_out.pattern,
            p.returncode,
            out)
    else:
        if expect_out == out.strip() and p.returncode == expect_code:
            return u"%d OK: %s" % (index, command_line)
        return u"%d NG: %s\t(expect exit=%d)\t'%s'\t(actual exit=%d)\t'%s'" % (
            index,
            command_line,
            expect_code,
            expect_out,
            p.returncode,
            out)


def process(*args, **kwargs):
    print run_test(*args, **kwargs)


def main():
    patterns = [
        ('monitor -H 127.0.0.1 -p check_procs -o ""',
         re.compile('PROCS OK: .* processes'), 0),
        ('monitor -H 127.0.0.1 -p check_procs -o "-w 0"',
         re.compile('PROCS WARNING: .* processes'), 1),
        ('monitor -H 127.0.0.1 -p check_procs -o "-c 0"',
         re.compile('PROCS CRITICAL: .* processes'), 2),
        ('monitor -H 127.0.0.1 -p check_procs -o "-c 0 -a"',
         re.compile('Usage:'), 3),
        ('monitor -H 127.0.0.1 -p notfound -o ""',
         'stdout=, stderr=/bin/sh: /usr/local/bin/notfound: No such file or directory', 127),
        ('monitor -X 127.0.0.1 -H 127.0.0.1 -p check_procs -o ""',
         re.compile('PROCS OK: .* processes'), 0),
        ('monitor -X 127.0.0.1 -H 127.0.0.1 -p check_procs -o "-w 0"',
         re.compile('PROCS WARNING: .* processes'), 1),
        ('monitor -X 127.0.0.1 -H 127.0.0.1 -p check_procs -o "-c 0"',
         re.compile('PROCS CRITICAL: .* processes'), 2),
        ('monitor -X 127.0.0.1 -H 127.0.0.1 -p notfound -o ""',
         'stdout=, stderr=/bin/sh: /usr/local/bin/notfound: No such file or directory', 127),
        ('monitor -H 127.0.1.1 -p check_procs -o ""',
         'ERROR(happo-agent): happo-agent returns 403 Access Denied', 2),
        ('monitor -X 127.0.0.1 -H 127.0.1.1 -p check_procs -o ""',
         'ERROR(happo-agent): happo-agent returns 403 Access Denied', 2),
        ('monitor -X 127.0.1.1 -H 127.0.0.1 -p check_procs -o ""',
         'ERROR(happo-agent): happo-agent returns 403 Access Denied', 2),
        ('monitor -H 192.168.0.1 -p check_procs -o ""',
         'ERROR(happo-agent): Post https://192.168.0.1:6777/monitor: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)', 2),
        ('monitor -X 127.0.0.1 -H 192.168.0.1 -p check_procs -o ""',
         'ERROR(happo-agent): Post https://127.0.0.1:6777/proxy: net/http: request canceled (Client.Timeout exceeded while awaiting headers)', 2),
        ('monitor -X 192.168.0.1 -H 127.0.0.1 -p check_procs -o ""',
         'ERROR(happo-agent): Post https://192.168.0.1:6777/proxy: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)', 2),
        ('monitor -H 192.168.0.1 -p check_procs -o "" -t 1',
         'ERROR(happo-agent): Post https://192.168.0.1:6777/monitor: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)', 2),
        ('monitor -X 127.0.0.1 -H 192.168.0.1 -p check_procs -o "" -t 1',
         'ERROR(happo-agent): Post https://127.0.0.1:6777/proxy: net/http: request canceled (Client.Timeout exceeded while awaiting headers)', 2),
        ('monitor -X 192.168.0.1 -H 127.0.0.1 -p check_procs -o "" -t 1',
         'ERROR(happo-agent): Post https://192.168.0.1:6777/proxy: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)', 2),
        ('monitor -H 192.168.0.1 -p check_procs -o "" -t 300',
         re.compile('ERROR\(happo-agent\): Post https://192.168.0.1:6777/monitor: dial tcp 192.168.0.1:6777: getsockopt:.*'), 2),
        ('monitor -X 127.0.0.1 -H 192.168.0.1 -p check_procs -o "" -t 300',
         'UNKNOWN(happo-agent): happo-agent returns 504', 3),
        ('monitor -X 192.168.0.1 -H 127.0.0.1 -p check_procs -o "" -t 300',
         re.compile('ERROR\(happo-agent\): Post https://192.168.0.1:6777/proxy: dial tcp 192.168.0.1:6777: getsockopt:.*'), 2),
    ]

    processes = []
    for index, pattern in enumerate(patterns):
        processes.append(Process(target=process, args=(
            index + 1, "./check_happo monitor %s" % (pattern[0]), pattern[1], pattern[2])))
    for p in processes:
        p.start()
    for p in processes:
        p.join()


if __name__ == "__main__":
    main()
