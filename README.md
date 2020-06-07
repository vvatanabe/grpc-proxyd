# grpc-proxyd

## Description

grpc-proxyd is a daemon that allows to easily configure grpc routing with YAML files. It's just a proof of concept.

## Installation

### Go

If you have the Go(go1.14+) installed, you can also install it with go get command.

```sh
$ go get github.com/vvatanabe/grpc-proxyd
```

## Usage

```sh
Usage:
  grpc-proxyd [flags]

Flags:
      --cert_file string   cert file path [GRPC_PROXYD_CERT_FILE]
  -c, --config string      config file path (default "config.yml")
  -h, --help               help for grpc-proxyd
      --key_file string    key file path [GRPC_PROXYD_KEY_FILE]
  -p, --port int           listen port [GRPC_PROXYD_PORT] (default 50051)
  -v, --verbose            debug mode [GRPC_PROXYD_VERBOSE]
```

## Config File

### YAML

```yaml
routes:
  - match: "/echo.EchoService/Echo"
    addr: "localhost:3001"
  - match: "/download.DownloadService/Download"
    addr: "localhost:3002"
  - match: "/upload.UploadService/Upload"
    addr: "localhost:3003"
```

## Acknowledgments

- [mwitkow/grpc-proxy](https://github.com/mwitkow/grpc-proxy)
- [mwitkow/grpc-proxy forks](https://github.com/mwitkow/grpc-proxy/network/members)
- [devsu/grpc-proxy](https://github.com/devsu/grpc-proxy)

## Bugs and Feedback

For bugs, questions and discussions please use the GitHub Issues.

## License

`grpc-proxyd` is released under the Apache 2.0 license.

## Author

[vvatanabe](https://github.com/vvatanabe)