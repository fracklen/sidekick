package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"sidekick"
	"strings"
	"syscall"
	"time"
	"vulcanClient"
)

var (
	dockerURL        = flag.String("docker-url", "unix:///var/run/docker.sock", "Docker socket file/url")
	expectedHTTPCode = flag.Int("expected-http-code", 200, "Expected Http Code from health check")
	interval         = flag.Int("interval", 10, "Health check interval")
	containerName    = flag.String("container", "2dc43851e93f", "Container ID/Name")
	virtualHostnames = flag.String("hostname", "www.example.org", "Comma-separated Virtual Hostnames")
	exposedPort      = flag.String("port", "8080", "Port")
	etcdURL          = flag.String("etcd", "http://172.16.42.43:4001", "Etcd endpoint")
	httpMethod       = flag.String("http-method", "GET", "Http Method for health check")
	healthURL        = flag.String("health-url", "/", "Health check path (include prefix slash)")
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

	// Create new Vulcan client with an TTL of 3 times the health check interval
	vc = vulcanClient.New(*etcdURL, uint64(3 * (*interval)))
}

func main() {
	var endpoint, containerID string
	var err interface{}

	for {
		endpoint, containerID, err = sidekick.FindEndpoint(*dockerURL, *containerName, *exposedPort)
		if err != nil {
			log.Printf("Error finding endpoint: %+v", err)
			wait(1)
		} else {
			break
		}
	}

	log.Printf("Endpoint: %s", endpoint)
	defer vc.Delete(*upstream, containerID)
	go trap(containerID)

	for {
		if ping(endpoint, *healthURL, *httpMethod, *expectedHTTPCode, *verbose) {
			if *verbose {
				log.Printf("OK")
			}
			vc.Set(*upstream, containerID, endpoint, virtualHostnameList, *location, *path)
			wait(*interval)
		} else {
			if *verbose {
				log.Printf("Failed")
			}
			vc.Delete(*upstream, containerID)
			wait(5)
		}

	}
}

func trap(containerID string) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-c
		vc.Delete(*upstream, containerID)
		os.Exit(0)
	}()
}

func ping(endpoint string, healthURL string, method string, expectedHTTPCode int, verbose bool) (result bool) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Err pinging: %s", endpoint)
			log.Printf("%+v", r)
			result = false
		}
	}()

	uri := fmt.Sprintf("http://%s%s", endpoint, healthURL)
	u, _ := url.Parse(uri)
	if verbose {
		log.Printf("Ping %+v", uri)
	}

	return sidekick.CheckURL(u, method, expectedHTTPCode, verbose)
}

func wait(secs int) {
	time.Sleep(time.Duration(secs) * time.Second)
}
