package main

import (
  "sidekick"
  "log"
)

func main() {
  endpoint, err := sidekick.FindEndpoint()
  if err != nil {
    log.Fatal(err)
  }
  log.Printf("Endpoint: %s", endpoint)
}

