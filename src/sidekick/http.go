package sidekick

import (
 "net/http"
 "net/url"
 "net"
 "time"
)

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
    return net.DialTimeout(network, addr, timeout)
}

func CheckUrl(u *url.URL, method string, expectedStatus int) bool {
  httpTransport := &http.Transport{
    Dial: dialTimeout,
  }

  u.Path = ""
  client := http.Client{Transport: httpTransport}

  req, err := http.NewRequest(method, u.String(), nil)
  if err != nil {
    return false
  }

  resp, err := client.Do(req)
  if err != nil {
    return false
  }
  defer resp.Body.Close()

  if resp.StatusCode == expectedStatus {
    return true
  }
  return false
}
