# Patch compose

Patching compose assumes [external-network] for `run-o11y-run` is active. Please refer to [docs](./external-network.md) for details.

## Usage

To patch any docker compose file with a `o11y` network, please run:
```sh
run-o11y-run  patch-compose  -f ${PATH_TO_DOCKER_COMPOSE_YAML}
```

iw will inject `o11y`network configuration to compose and inject network to ALL services within customer owned docker-compose.

By default injected env var will point to http exporters:
```yaml
OTEL_EXPORTER_OTLP_ENDPOINT: otel-collector:4317
```