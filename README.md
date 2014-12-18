# Sidekick

A small Go-application, which sits next to your own application and acts as its sidekick.

## What is a sidekick?

In this sense, a sidekick will poke to its "hero" with a specific `interval`.
It will poke a `health-url`, where it expect the "hero" to return an `expected-http-code`.
If every went OK, the sidekick will add the "hero" as an upstream in [Vulcand](https://github.com/mailgun/vulcand) via [etcd](https://github.com/coreos/go-etcd).

The sidekick will check each `interval` seconds, and pull it out of Vulcan, if it gets an unexpected return code.

## Flags

The following flags are available:

| Name                 | Type   | Description                              | Default                         |
| -------------------- | ------ | ---------------------------------------- | ------------------------------- |
| `docker-url`         | String | Docker socket file/url                   | `"unix:///var/run/docker.sock"` |
| `expected-http-code` | Int    | Expected HTTP Code from health check     | `200`                           |
| `interval`           | Int    | Health check interval                    | `10`                            |
| `container`          | String | Container ID/Name                        | `"2dc43851e93f"`                |
| `hostname`           | String | Comma-separated Virtual Hostnames        | `"www.example.org"`             |
| `port`               | String | Port                                     | `"8080"`                        |
| `etcd`               | String | etcd endpoint                            | `"http://172.16.42.43:4001"`    |
| `http-method`        | String | HTTP Method for health check             | `"GET"`                         |
| `health-url`         | String | Health check path (include prefix slash) | `"/"`                           |
| `upstream`           | String | Upstream name                            | `"foobar"`                      |
| `location`           | String | Location name                            | `"loc1"`                        |
| `path`               | String | Location path                            | `"/"`                           |
| `verbose`            | Bool   | Verbose                                  | `false`                         |
