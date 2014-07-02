package main

import (
  "sidekick"
  "log"
  "flag"
  "os"
  "net/url"
  "fmt"
  "time"
)

var (
  //dockerUrl = flag.String("docker-url", "unix:///var/run/docker.sock", "Docker socket file/url")
  dockerUrl = flag.String("docker-url", "http://172.16.42.43:4243", "Docker socket file/url")
  containerId = flag.String("container", "a380ed47f37f", "Container ID")
  hostname = os.Getenv("HOSTNAME")
  exposedPort = flag.String("port", "2345", "Port")
  etcdUrl = flag.String("etcd", "http://0.0.0.0:4001", "Etcd endpoint")
  expectedHttpCode = flag.Int("expected-http-code", 200, "Expected Http Code from health check")
  httpMethod = flag.String("http-method", "GET", "Http Method for health check")
  healthUrl = flag.String("health-url", "/", "Health check path (include prefix slash)")
)

func main() {
  endpoint, err := sidekick.FindEndpoint(*dockerUrl, *containerId, *exposedPort)
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Endpoint: %s", endpoint)

  for {
    if ping(endpoint, *healthUrl, *httpMethod, *expectedHttpCode) {
      fmt.Println("Alive")
    } else {
      fmt.Println("Dead")
    }
    time.Sleep(time.Duration(1 * time.Second))
  }
}

func ping(endpoint string, healthUrl string, method string, expectedHttpCode int) bool{
  uri  := fmt.Sprintf("http://%s%s", endpoint, healthUrl)
  u, _ := url.Parse(uri)
  return sidekick.CheckUrl(u, method, expectedHttpCode)
}



