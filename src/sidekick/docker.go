package sidekick

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
)

var (
	errNotFound = errors.New("Not Found")
)

// FindEndpoint - Finds the endpoint of the docker container
func FindEndpoint(dockerURL string, containerName string, exposedPort string) (string, string, error) {
	info, err := findApplicationContainer(dockerURL, containerName)
	if err != nil {
		return "", "", err
	}

	internalPort := fmt.Sprintf("%s/tcp", exposedPort)

	portBindings := info.NetworkSettings.Ports[internalPort]

	for _, host := range portBindings {
		endpoint := fmt.Sprintf("%s:%s", host.HostIP, host.HostPort)
		return endpoint, info.ID, nil
	}

	return "", "", errNotFound
}

func findApplicationContainer(dockerURL string, containerName string) (*ContainerInfo, error) {
	u, err := url.Parse(dockerURL)

	if err != nil {
		log.Fatal("Error parsing docker-url")
	}

	client := newHTTPClient(u)

	uri := fmt.Sprintf("%s/v1.12/containers/%s/json", u.String(), containerName)

	info := &ContainerInfo{}
	data, err := doRequest(client, "GET", uri, nil)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

func doRequest(client *http.Client, method string, path string, body []byte) ([]byte, error) {
	b := bytes.NewBuffer(body)

	req, err := http.NewRequest(method, path, b)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, errNotFound
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("%s: %s", resp.Status, data)
	}

	return data, nil
}

func newHTTPClient(u *url.URL) *http.Client {
	httpTransport := &http.Transport{}

	if u.Scheme == "unix" {
		socketPath := u.Path
		unixDial := func(proto string, addr string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		}

		httpTransport.Dial = unixDial

		// Override the main URL object so the HTTP lib won't complain
		u.Scheme = "http"
		u.Host = "unix.sock"
	}

	u.Path = ""
	return &http.Client{Transport: httpTransport}
}
