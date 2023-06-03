# Patch compose

The `patch` command in run-o11y-run allows you to seamlessly enhance your Docker Compose file by injecting the `o11y` network configuration. This command assumes that you have already started a bridged `o11y` network using the `run-o11y-run start --external-network` command. For detailed information on setting up the **external network**, please refer to the [documentation](./external-network.md).

## Usage

To patch any Docker Compose file and add the `o11y` network, simply execute the following command:

```sh
run-o11y-run patch -f ${PATH_TO_DOCKER_COMPOSE_YAML}
```

This command will inject the necessary network configuration into your Docker Compose file, ensuring that all services within your customer-owned environment are connected to the `o11y` network.

By default, the injected environment variables will point to the following exporters:

```yaml
OTEL_EXPORTER_OTLP_ENDPOINT: otel-collector:4317
OTEL_EXPORTER_ZIPKIN_ENDPOINT: tempo:9411
```
With the patched Docker Compose file, you can seamlessly integrate observability capabilities and leverage the power of the o11y stack to monitor and trace your services effectively.
