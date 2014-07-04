package vulcanClient

import (
  "fmt"
  "github.com/coreos/go-etcd/etcd"
  "errors"
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

  c.etcdClient.Set(fmt.Sprintf("vulcand/upstreams/%s/endpoints/%s", upstream, containerId), fmt.Sprintf("http://%s", endpoint), 0)

  for _, virtualHostname := range virtualHostnameList {
    c.etcdClient.Set(fmt.Sprintf("vulcand/hosts/%s/locations/%s/path", virtualHostname, location), path, 0)
    c.etcdClient.Set(fmt.Sprintf("vulcand/hosts/%s/locations/%s/upstream", virtualHostname, location), upstream, 0)
  }

  return nil
}

func (c VulcanClient) Delete(upstream, containerId string, virtualHostnameList []string, location string) error {
  if len(virtualHostnameList) == 0 {
    return errors.New("No hostnames given")
  }

  c.etcdClient.Delete(fmt.Sprintf("vulcand/upstreams/%s/endpoints/%s", upstream, containerId), false)

  for _, virtualHostname := range virtualHostnameList {
    c.etcdClient.Delete(fmt.Sprintf("vulcand/hosts/%s/locations/%s", virtualHostname, location), true)
  }

  return nil
}
