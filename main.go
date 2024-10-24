package main

import "time"

const apiPort = 8080
const DiscoveryPort = 32227
const DefaultAlpacaApiPort = 11111
const ListenIP = "127.0.0.1"
const Location = "Earth"

func main() {
	// Load initial switch values
	MhpSetInit()
	discovery := NewDiscoverySever(DiscoveryPort, apiPort)
	api := NewApiServer(apiPort)
	go discovery.Start()
	defer discovery.Close()
	go api.Start()
	for {
		time.Sleep(10 * time.Second)
	}
}
