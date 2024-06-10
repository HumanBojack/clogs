# clogs
Contextual logs, a tool used to display multiple log files.

# Build
```bash
go build
```

# Usage
```bash
clogs <log file 1> <datetime start> <datetime end> <log file 2> <datetime start> <datetime end> ...
```

# Example
In this example, we're using logs from [loghub](https://github.com/logpai/loghub) to display logs from apache, linux, and openssh logs.
```bash
logs/apache.log 1 20 logs/linux.log 0 15 logs/openssh.log 0 15
```