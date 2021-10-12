# Ha-cert-manager-for-traefik
The app can be used to backup and replication the acme file across multiple traefik nodes

[![Go Report Card](https://goreportcard.com/badge/github.com/jakubjastrabik/treafikCertKVstore)](https://goreportcard.com/report/github.com/jakubjastrabik/treafikCertKVstore)
[![goreleaser](https://github.com/jakubjastrabik/treafikCertKVstore/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/jakubjastrabik/treafikCertKVstore/actions/workflows/goreleaser.yml)
[![Docker](https://github.com/jakubjastrabik/treafikCertKVstore/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/jakubjastrabik/treafikCertKVstore/actions/workflows/docker-publish.yml)

<p align="center">
  <img src="[http://some_place.com/image.png](https://github.com/jakubjastrabik/treafikCertKVstore/tree/master/docu/images/basic-topo.svg)" />
</p>


## Possible configuration tags

| Flag Name             | Description                                           | Default Value          |
|--                     |--                                                     |--                      |
| members               | comma separated list of members                       |                        |
| httpPort              | Port to be used for connection                        | 7900                   |
| httpAddress           | Address to be use for connection                      | 0.0.0.0                |
| traefikCertLocalStore | path with file name where are stored certificates     | /etc/traefik/acme.json |
| consulKey             | Consul key for storage certificates                   | traefik/acme.json      |
| path                  | Log file path with name                               | /var/log/hacert.log    |
| logLevel              | Possible level of debugging, DEBUG, WARN, INFO, ERROR | DEBUG                  |
| appName               | Aplication tag for logging                            | traefikCertKVStore     |
| backupCount           | Count of rotated backup version                       | 3                      |
| waitAfterStart        | Waiting to start to do tasks after started in seconds | 5                      |

## Prometheus metrics

| metricsPath           |	URL path for surfacing collected metrics              | /metrics	             |
| --                    | --                                                    | --                     |

## Grafana
Grafana ID: 15100 https://grafana.com/grafana/dashboards/15100

## Usage

### From Binary

``` Bash
./traefikCertKVStore -members="192.168.1.10,192.168.1.11"
```
### From container

``` bash
docker pull ghcr.io/jakubjastrabik/treafikcertkvstore:latest
docker run -p 7900:7900 --name treafikcertkvstore --env-file=.env ghcr.io/jakubjastrabik/treafikcertkvstore 
```

### Build from source

* example for linux amd64
  
``` Bash
git clone https://github.com/jakubjastrabik/treafikCertKVstore.git
go get
GOOS=linux GOARCH=amd64 go build -o traefikCertKVStore 
```
