# Hermes

Hetzner Failover IP Monitor - Automatically updates failover IP routing via Hetzner Robot API.

## Build

```bash
go build -o hermes .
```

## Configuration

Create `/home/nixos/robot.json`:

```json
{
  "user": "#your_user",
  "password": "your_password",
  "failover_ip": "203.0.113.1"
}
```

Or use environment variables:

```bash
export HETZNER_USER="#your_user"
export HETZNER_PASS="your_password"
export FAILOVER_IP="203.0.113.1"
```

## Usage

```bash
# Normal run
./hermes

# Dry run (test without API calls)
./hermes --dry-run
```

## Keepalived Integration

Add to `keepalived.conf`:

```bash
vrrp_instance VI_1 {
    notify_master "/usr/local/bin/hermes"
}
```

## Development

```bash
devbox shell
go build
```



