package main

import (
  "sidekick"
  "vulcanClient"

  "log"
  "flag"
  "os"
  "os/signal"
  "net/url"
  "fmt"
  "time"
  "strings"
)

var (
  //dockerUrl = flag.String("docker-url", "unix:///var/run/docker.sock", "Docker socket file/url")
  expectedHttpCode = flag.Int("expected-http-code", 200, "Expected Http Code from health check")
  dockerUrl        = flag.String("docker-url", "http://172.16.42.43:4243", "Docker socket file/url")
  containerId      = flag.String("container", "2dc43851e93f", "Container ID")
  virtualHostnames = flag.String("hostname", "www.example.org", "Comma-separated Virtual Hostnames")
  exposedPort      = flag.String("port", "8080", "Port")
  etcdUrl          = flag.String("etcd", "http://172.16.42.43:4001", "Etcd endpoint")
  httpMethod       = flag.String("http-method", "GET", "Http Method for health check")
  healthUrl        = flag.String("health-url", "/", "Health check path (include prefix slash)")
  upstream         = flag.String("upstream", "foobar", "Upstream name")
  location         = flag.String("location", "loc1", "Location name")
  path             = flag.String("path", "/", "Location path")

  vc                  = &vulcanClient.VulcanClient{}
  virtualHostnameList = make([]string, 0)
)

func init() {
  flag.Parse()

  for _, vh := range strings.Split(*virtualHostnames, ",") {
    virtualHostnameList = append(virtualHostnameList, vh)
  }

  if len(virtualHostnameList) == 0 || virtualHostnameList[0] == "www.example.org" {
    log.Fatal("No hostname given")
  }

  vc = vulcanClient.New(*etcdUrl)
}

func main() {
  go trap()

  endpoint, err := sidekick.FindEndpoint(*dockerUrl, *containerId, *exposedPort)
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Endpoint: %s", endpoint)

  for {
    if ping(endpoint, *healthUrl, *httpMethod, *expectedHttpCode) {
      vc.Set(*upstream, *containerId, endpoint, virtualHostnameList, *location, *path)
    } else {
      vc.Delete(*upstream, *containerId, virtualHostnameList, *location)
    }
    time.Sleep(time.Duration(1 * time.Second))
  }
}

func ping(endpoint string, healthUrl string, method string, expectedHttpCode int) bool {
  defer func() bool {
    if r := recover(); r != nil {
      log.Printf("Err pinging: %s", endpoint)
      log.Printf("%+v", r)
      return false
    }
    return true
  }()

  uri  := fmt.Sprintf("http://%s%s", endpoint, healthUrl)
  u, _ := url.Parse(uri)

  return sidekick.CheckUrl(u, method, expectedHttpCode)
}

func trap() {
  c := make(chan os.Signal, 1)
  signal.Notify(c, os.Interrupt)

  for {
    <- c
    vc.Delete(*upstream, *containerId, virtualHostnameList, *location)
    os.Exit(0)
  }
}
