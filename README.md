# parityexporter
Parity / Etherscan prometheus exporter

### Env
* `EXPORTER_PORT`
* `ETHERSCAN_URL`
* `ETHERSCAN_KEY`
* `PARITY_URL`
* `ENDPOINT` - default "/"

### Build

```shell script
$ docker build -t exporter .
```

### Run

```shell script
$ docker run -e EXPORTER_PORT=<EXPORTER_PORT> \
-e ETHERSCAN_URL=<ETHERSCAN_URL> \
-e ETHERSCAN_KEY=<ETHERSCAN_KEY> \
-e PARITY_URL=<PARITY_URL> \
-e ENDPOINT=<ENDPOINT> \
-p <EXPORTER_PORT>:<EXPORTER_PORT> exporter
```
