package vulcanClient

import (
	"errors"
	"fmt"
	"github.com/coreos/go-etcd/etcd"
	"log"
)

// VulcanClient - A etcd wrapper
type VulcanClient struct {
	etcdClient *etcd.Client
	ttl        uint64
}

// New - returns a new VulcanClient
func New(etcdURL string, ttl uint64) *VulcanClient {
	return &VulcanClient{etcd.NewClient([]string{etcdURL}), ttl}
}

// Set - creates a new endpoint for container in etcd
func (c VulcanClient) Set(upstream, containerID, endpoint string, virtualHostnameList []string, location, path string) error {
	if len(virtualHostnameList) == 0 {
		return errors.New("No hostnames given")
	}

	uppath := fmt.Sprintf("vulcand/upstreams/%s/endpoints/%s", upstream, containerID)
	value := fmt.Sprintf("http://%s", endpoint)

	prev, err := c.etcdClient.Get(uppath, false, false)

	if err != nil {
		log.Printf("Error Updating %+v", err)
	} else {
		if (*prev.Node).Value != value {
			log.Printf("Updating %+v -> %+v", uppath, value)
		}
	}

	_, err = c.etcdClient.Set(uppath, value, c.ttl)
	if err != nil {
		log.Printf("Error updating %+v -> %+v: %+v", uppath, value, err)
	}

	for _, virtualHostname := range virtualHostnameList {
		c.etcdClient.Set(fmt.Sprintf("vulcand/hosts/%s/locations/%s/path", virtualHostname, location), path, 0)
		c.etcdClient.Set(fmt.Sprintf("vulcand/hosts/%s/locations/%s/upstream", virtualHostname, location), upstream, 0)
	}

	return nil
}

// Delete - deletes the endpoint for given container
func (c VulcanClient) Delete(upstream, containerID string) error {
	c.etcdClient.Delete(fmt.Sprintf("vulcand/upstreams/%s/endpoints/%s", upstream, containerID), false)

	return nil
}
