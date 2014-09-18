package sidekick

import (
	"log"
	"net"
	"net/http"
	"net/url"
	"time"
)

var timeout = time.Duration(2 * time.Second)

func dialTimeout(network, addr string) (net.Conn, error) {
	return net.DialTimeout(network, addr, timeout)
}

// CheckURL - Checks that the given URL responds with expectedStatus
func CheckURL(u *url.URL, method string, expectedStatus int, verbose bool) bool {
	httpTransport := &http.Transport{
		Dial: dialTimeout,
	}

	client := http.Client{Transport: httpTransport}

	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		if verbose {
			log.Printf("Err checking:%+v", err)
		}
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		if verbose {
			log.Printf("Err checking:%+v", err)
		}
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == expectedStatus {
		return true
	}
	if verbose {
		log.Printf("Err checking:%+v", resp.StatusCode)
	}
	return false
}
