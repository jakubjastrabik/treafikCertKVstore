# Ha-cert-manager-for-traefik
The app can be used to backup and replication the acme file across multiple traefik nodes

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

## Usage

``` Bash
./traefikCertKVStore -members="192.168.1.10,192.168.1.11"
```

### Build from source

* example for linux amd64
  
``` Bash
git clone https://github.com/jakubjastrabik/treafikCertKVstore.git
go get
go build -o traefikCertKVStore GOOS=linux GOARCH=amd64
```
