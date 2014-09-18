package vulcanClient

import (
  "fmt"
  "github.com/coreos/go-etcd/etcd"
  "errors"
  "log"
)

type VulcanClient struct {
  etcdClient *etcd.Client
}

func New(etcdUrl string) *VulcanClient {
  return &VulcanClient{etcd.NewClient([]string{etcdUrl})}
}

func (c VulcanClient) Set(upstream, containerId, endpoint string, virtualHostnameList []string, location, path string) error {
  if len(virtualHostnameList) == 0 {
    return errors.New("No hostnames given")
  }

  uppath := fmt.Sprintf("vulcand/upstreams/%s/endpoints/%s", upstream, containerId)
  value := fmt.Sprintf("http://%s", endpoint)

  prev, err := c.etcdClient.Get(uppath, false, false)

  if err != nil {
    log.Printf("Error Updating %+v", err)
  } else {
    if (*prev.Node).Value != value {
      log.Printf("Updating %+v -> %+v", uppath, value)
    }
  }

  _, err = c.etcdClient.Set(uppath, value, 0)
  if err != nil {
   log.Printf("Error updating %+v -> %+v: %+v", uppath, value, err)
  }

  for _, virtualHostname := range virtualHostnameList {
    c.etcdClient.Set(fmt.Sprintf("vulcand/hosts/%s/locations/%s/path", virtualHostname, location), path, 0)
    c.etcdClient.Set(fmt.Sprintf("vulcand/hosts/%s/locations/%s/upstream", virtualHostname, location), upstream, 0)
  }

  return nil
}

func (c VulcanClient) Delete(upstream, containerId string) error {
  c.etcdClient.Delete(fmt.Sprintf("vulcand/upstreams/%s/endpoints/%s", upstream, containerId), false)

  return nil
}
