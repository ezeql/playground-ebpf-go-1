# playground-ebpf-go-1

Learning playground for zero-code observability with eBPF and Go.

A toy Go service is traced by [OpenTelemetry eBPF Instrumentation (OBI)](https://opentelemetry.io/docs/zero-code/obi/) without any SDK. Traces and RED metrics go to [grafana/otel-lgtm](https://github.com/grafana/docker-otel-lgtm).

## Requirements

- A Linux kernel with BTF (`/sys/kernel/btf/vmlinux`). Tested on OrbStack (kernel 7.2, arm64).
- Docker with Compose.

## Run

```bash
make run    # app + MongoDB + load generator + OBI + LGTM, then follows OBI's span log
make down
```

- App: http://localhost:8080/hello
- Grafana: http://localhost:3300 (Explore → Tempo for traces, Prometheus for metrics)

## What's inside

| Path | What |
|---|---|
| `app/` | `/hello` → `/work` → MongoDB (insert + find), with random latency and 500s |
| `compose.yaml` | app, mongo, load generator, OBI, LGTM, and the labs |
| `trace.bt` | bpftrace script used by `make bpftrace` |
| `uprobe/` | A hand-written eBPF program ([cilium/ebpf](https://github.com/cilium/ebpf)) that counts `ServeHTTP` calls in the app |

## Labs

- `make bpftrace`: prints every `openat` syscall in the VM.
- `make uprobe`: prints how many HTTP requests the app has served.

## Notes

- OBI, bpftrace and the uprobe run as privileged containers. On OrbStack, a privileged container with host PID can reach your macOS home, because OrbStack mounts it into the VM. Only run images you trust.
- OBI shares the app's PID namespace, so recreating only `app` makes OBI exit. Use `make run`, which recreates everything.
- Ports are bound to `127.0.0.1`.
