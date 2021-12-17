# time-format

A tiny convenience tool to convert time formats.

```
Example:
    time-format -u -- $(date)

Flags:
   --version       Displays the program version string.
-h --help          Displays help with available flag, subcommand, and positional value parameters.
   --rfc3339-utc   Format time as RFC 3339 UTC (default)
-r --rfc3339       Format time as RFC 3339
-u --unix          Format time as unix time
-m --unix-milli    Format time as unix milli
-f --unix-float    Format time as unix float
```

`time-format` currently parses:

- RFC3339Nano
- UnixDate
- Unix timestamp
- Unix millisecond timestamp
- Unix float timestamp
