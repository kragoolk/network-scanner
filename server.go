package main

import (
	"flag"
	"fmt"
	"os"
	
	"network-scanner/pkg/api"
)

func main() {
	network := flag.String("network", "192.168.1.0/24", "Network to scan (CIDR notation)")
	port := flag.Int("port", 8080, "API server port")
	flag.Parse()
	
	fmt.Println("=== Network Scanner API v0.8 ===")
	
	server := api.NewServer(*network)
	err := server.Start(*port)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

