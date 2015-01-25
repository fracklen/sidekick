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

## License

The MIT License (MIT)

Copyright (c) 2014 Lokalebasen.dk 

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
