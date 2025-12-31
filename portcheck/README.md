# portcheck

A fast, concurrent TCP port scanner written in Go.

## Features

- Concurrent scanning using worker pools (CPU cores × 10 workers)
- Configurable port ranges and lists
- 3-second connection timeout

## Installation

```bash
go build -o portcheck
```

## Usage

```bash
# Scan all ports (1-65535) on a host
./portcheck <host>

# Scan specific ports
./portcheck <host> <ports>
```

### Port Specification

Ports can be specified as:
- Single ports: `80`
- Ranges: `1-1000`
- Comma-separated combinations: `22,80,443,8000-9000`

### Examples

```bash
# Scan all ports on localhost
./portcheck localhost

# Scan common web ports
./portcheck example.com 80,443

# Scan a range of ports
./portcheck 192.168.1.1 1-1024

# Scan mixed ports and ranges
./portcheck example.com 22,80,443,8000-8100
```

## Output

Open ports are printed to stdout:

```
SUCCESS: 192.168.1.1:22
SUCCESS: 192.168.1.1:80
```

## License

MIT
