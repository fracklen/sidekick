package main

import (
  "sidekick"
  "vulcanClient"
  "os"
  "os/signal"
  "syscall"
  "log"
  "flag"
  "net/url"
  "fmt"
  "time"
  "strings"
)

var (
  dockerUrl        = flag.String("docker-url", "unix:///var/run/docker.sock", "Docker socket file/url")
  expectedHttpCode = flag.Int("expected-http-code", 200, "Expected Http Code from health check")
  interval         = flag.Int("interval", 10, "Health check interval")
  //dockerUrl        = flag.String("docker-url", "http://172.16.42.43:4243", "Docker socket file/url")
  containerName    = flag.String("container", "2dc43851e93f", "Container ID/Name")
  virtualHostnames = flag.String("hostname", "www.example.org", "Comma-separated Virtual Hostnames")
  exposedPort      = flag.String("port", "8080", "Port")
  etcdUrl          = flag.String("etcd", "http://172.16.42.43:4001", "Etcd endpoint")
  httpMethod       = flag.String("http-method", "GET", "Http Method for health check")
  healthUrl        = flag.String("health-url", "/", "Health check path (include prefix slash)")
  upstream         = flag.String("upstream", "foobar", "Upstream name")
  location         = flag.String("location", "loc1", "Location name")
  path             = flag.String("path", "/", "Location path")
  verbose          = flag.Bool("verbose", false, "Verbose")

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
  var endpoint, containerId string
  var err interface{}

  for {
    endpoint, containerId, err = sidekick.FindEndpoint(*dockerUrl, *containerName, *exposedPort)
    if err != nil {
      log.Printf("Error finding endpoint: %+v", err)
      wait(1)
    } else {
      break
    }
  }

  log.Printf("Endpoint: %s", endpoint)
  defer vc.Delete(*upstream, containerId)
  go trap(containerId)

  for {
    if ping(endpoint, *healthUrl, *httpMethod, *expectedHttpCode, *verbose) {
      if *verbose {
        log.Printf("OK")
      }
      vc.Set(*upstream, containerId, endpoint, virtualHostnameList, *location, *path)
      wait(*interval)
    } else {
      if *verbose {
        log.Printf("Failed")
      }
      vc.Delete(*upstream, containerId)
      wait(5)
    }

  }
}

func trap(containerId string) {
  c := make(chan os.Signal, 1)
  signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

  for {
    <- c
    vc.Delete(*upstream, containerId)
    os.Exit(0)
  }
}

func ping(endpoint string, healthUrl string, method string, expectedHttpCode int, verbose bool) (result bool) {
  defer func() {
    if r := recover(); r != nil {
      log.Printf("Err pinging: %s", endpoint)
      log.Printf("%+v", r)
      result = false
    }
  }()

  uri  := fmt.Sprintf("http://%s%s", endpoint, healthUrl)
  u, _ := url.Parse(uri)
  if verbose {
    log.Printf("Ping %+v", uri)
  }

  return sidekick.CheckUrl(u, method, expectedHttpCode, verbose)
}

func wait(secs int) {
  time.Sleep(time.Duration(secs) * time.Second)
}
