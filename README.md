# go-telnet

This utility is an attempt to replicate the telnet utility's functionality.
Its interface is pretty straightforward, however, note that it lacks most of
the flags available in the original telnet programme.

See: https://linux.die.net/man/1/telnet for reference.

## Usage

```go-telnet [HOST ADDRESS] [PORT] [FLAGS]```

## Supported flags

### timeout (-t)
Sets the timeout that will be used for exiting the programme (default is 30s).

## TODO
1) Add the ability to write output data to a file.
