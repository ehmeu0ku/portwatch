# portwatch

A lightweight CLI daemon that monitors port activity and alerts on unexpected listeners.

---

## Installation

```bash
go install github.com/yourname/portwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/portwatch.git && cd portwatch && go build -o portwatch .
```

---

## Usage

Start the daemon with a config file defining your expected listeners:

```bash
portwatch --config config.yaml --interval 10s
```

**Example `config.yaml`:**

```yaml
allowed_ports:
  - 22
  - 80
  - 443
alert:
  method: log
  output: /var/log/portwatch.log
```

portwatch will poll active listeners at the specified interval and log (or notify) whenever a port outside your allowlist is detected.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `config.yaml` | Path to config file |
| `--interval` | `30s` | Polling interval |
| `--verbose` | `false` | Enable verbose output |

---

## How It Works

portwatch reads active TCP/UDP listeners from the system at a regular interval and compares them against your defined allowlist. Any unexpected port triggers a configurable alert — log entry, stdout warning, or webhook call.

---

## License

MIT © 2024 yourname