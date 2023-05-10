# run-o11y-run

A single-binary 🌯 wrapper around `docker compose` with embedded configurations to effortlessly run your local observability stack

## Prerequisites

`run-o11y-run` depends on the latest version of `docker`, which includes the `docker compose` command.

## Quick Start

* Download the archive from [releases](https://github.service.anz/kowalskk/run-o11y-run/releases).
* Extract the archive
* Run `run-o11y-run`
    * **NOTE:** If you encounter the error "run-o11y-run" can't be opened because Apple cannot check it for malicious software:
    * Visit https://support.apple.com/en-au/guide/mac-help/mh40616/mac for guidance.
    * **TODO:** The app will be notarised in future updates.
* Configure your service to push data to one of the following endpoints. This can be done by setting the `OTEL_EXPORTER_OTLP_ENDPOINT` environment variable or updating your config file. Make sure to set the traffic to **insecure**:
    * OTLP (grpc): http://localhost:4317
    * OTLP (http): http://localhost:4318
    * Jaeger: http://localhost:14268
    * Zipkin: http://localhost:9411
* To exit gracefully, press `CTRL+C`.

## Commands

```sh
# basic
run-o11y-run

# clean
run-o11y-run -clean
```

## Links

* [Grafana Tempo](http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22tempo%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22tempo%22,%22uid%22:%22tempo%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D)
* [Grafana Prometheus](http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22prometheus%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22prometheus%22,%22uid%22:%22prometheus%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D)
* [Prometheus Direct](http://localhost:9090/)

## Services

The underlying observability stack is built on Grafana products. It includes the following components:

* Grafana
* Grafana Tempo
* OpenTelemetry Collector
* Prometheus

## Documentation

* [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/)
* [TraceQL](https://grafana.com/docs/tempo/latest/traceql/)

## Troubleshooting

`run-o11y-run` is built on top of Docker, and if you encounter any issues or things don't seem to be working as expected, please use the standard Docker debugging techniques.

Make sure you run `run-o11y-run -clean` to clean up the configurations before attempting any troubleshooting steps.

In case you need further assistance, refer to the [Docker documentation](https://docs.docker.com/) and [Docker troubleshooting guide](https://docs.docker.com/engine/troubleshooting/).
