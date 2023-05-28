# run-o11y-run

```sh
_____ _____ _____     _____ ___   ___   __ __     _____ _____ _____
| __  |  |  |   | |___|     |_  | |_  | |  |  |___| __  |  |  |   | |
|    -|  |  | | | |___|  |  |_| |_ _| |_|_   _|___|    -|  |  | | | |
|__|__|_____|_|___|   |_____|_____|_____| |_|     |__|__|_____|_|___|

```

A single-binary 🌯 wrapper around `docker compose` with embedded configurations to effortlessly run your local observability stack.

The underlying observability stack is built on [Grafana](https://grafana.com/) products and [OpenTelemetry](https://opentelemetry.io/). It includes the following services:

* [Grafana](https://grafana.com/oss/grafana/)
* [Grafana Loki](https://grafana.com/oss/loki/)
* [Grafana Tempo](https://grafana.com/oss/tempo/)
* [MinIO](https://min.io/)
* [OpenTelemetry Collector](https://opentelemetry.io/docs/collector/)
* [Prometheus](https://grafana.com/oss/prometheus/)


## Prerequisites

`run-o11y-run` depends on the latest version of [Docker Desktop](https://www.docker.com/products/docker-desktop/), which includes the `docker compose` command.

## Quick Start

### Install

#### brew

Install [brew](https://brew.sh/) and then run:

```sh
brew install krzko/tap/run-o11y-run
```

#### Download Binary

Download the latest version from the [Releases](https://github.com/krzko/run-o11y-run/releases) page.

### Run

* Run `run-o11y-run`
* Configure your service to push telemetry data to one of the following endpoints. This can be done by setting the `OTEL_EXPORTER_OTLP_ENDPOINT` environment variable or updating your config file. Make sure to set the traffic to **insecure**:
    * OTLP (grpc): http://localhost:4317
    * OTLP (http): http://localhost:4318
    * Jaeger: http://localhost:14268
    * Zipkin: http://localhost:9411
* Logs are procedd via two means:
  * Tailed from `/var/log/*.log` and `./*.log` on your local machine.
  * A Syslog RFC 3164 header format, `syslog` receiver operates on `localhost:8094`
* To exit gracefully, press `CTRL+C`.

## Commands

`run-o11y-run` is a powerful command-line tool that provides seamless management of your local observability stack. It offers three simple commands: `start`, `stop`, and `clean`.

```sh
$ run-o11y-run start

✨ Starting...
[+] Running 56/39
 ✔ tempo 5 layers [⣿⣿⣿⣿⣿]      0B/0B      Pulled                                                                                                           142.9s
 ⠦ minio 9 layers [⣿⣿⣿⣿⣿⣄⣿⣿⣿] 48.23MB/96.92MB Pulling                                                                                                      170.6s
 ⠦ otel-collector 6 layers [⣿⣿⣿⣿⣿⣷] 45.09MB/48.15MB Pulling                                                                                                170.6s
 ⠦ grafana 12 layers [⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿]      0B/0B      Pulling                                                                                                170.6s
 ⠦ prometheus 15 layers [⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿]      0B/0B      Pulling                                                                                          170.6s
 ✔ loki 9 layers [⣿⣿⣿⣿⣿⣿⣿⣿⣿]      0B/0B      Pulled                                                                                                         81.8s
```

### Start Command

The `start` command allows you to launch run-o11y-run containers and start your observability stack effortlessly. You can customise the behaviour using various flags.

One such flag is the `--registry` flag, which enables you to specify a Docker Registry from which to pull the required images. By default, run-o11y-run uses Docker Hub as the registry. Here's an example of using the `--registry flag`:

```sh
run-o11y-run start --registry <registry-url>
```

Replace `<registry-url>` with the URL of your desired Docker Registry.

To further enhance your setup, you can also utilise the `--external-network` flag, which enables integration of your own docker-compose configurations with run-o11y-run. This allows you to combine the services of run-o11y-run with your existing infrastructure seamlessly.

To start run-o11y-run in `detached` mode, use the `--detach` flag. This will start the containers in the background.

For more details on using the `--external-network` flag, refer to the [External Network](docs/external-network.md) guide.

### Stop Command

The `stop` command is used to gracefully stop the run-o11y-run containers. It ensures a clean shutdown of your observability stack. Here's an example of using the `stop` command:

```sh
run-o11y-run stop
```

### Clean Command

The `clean` command is used to stop and remove run-o11y-run containers, files, and networks. It helps you clean up your environment after using run-o11y-run. Here's an example of using the `clean` command:

```sh
run-o11y-run clean
```

## Local Service Links

* [Grafana Loki](http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22P8E80F9AEF21F6940%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22loki%22,%22uid%22:%22P8E80F9AEF21F6940%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D)
* [Grafana Tempo](http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22tempo%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22tempo%22,%22uid%22:%22tempo%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D)
* [Grafana Prometheus](http://localhost:3000/explore?orgId=1&left=%7B%22datasource%22:%22prometheus%22,%22queries%22:%5B%7B%22refId%22:%22A%22,%22datasource%22:%7B%22type%22:%22prometheus%22,%22uid%22:%22prometheus%22%7D%7D%5D,%22range%22:%7B%22from%22:%22now-1h%22,%22to%22:%22now%22%7D%7D)
* [Prometheus Direct](http://localhost:9090/)

## Documentation

* [LogQL](https://grafana.com/docs/loki/latest/logql/)
* [PromQL](https://prometheus.io/docs/prometheus/latest/querying/basics/)
* [TraceQL](https://grafana.com/docs/tempo/latest/traceql/)

## Troubleshooting

`run-o11y-run` is built on top of Docker, and if you encounter any issues or things don't seem to be working as expected, please use the standard Docker debugging techniques.

Make sure you run `run-o11y-run -clean` to clean up the configurations before attempting any troubleshooting steps.

In case you need further assistance, refer to the [Docker documentation](https://docs.docker.com/) and [Docker troubleshooting guide](https://docs.docker.com/engine/troubleshooting/).
