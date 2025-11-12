package main

import (
	"flag"
	"fmt"
	"strings"
	
	"network-scanner/pkg/lookup"
	"network-scanner/pkg/scanner"
	"network-scanner/pkg/storage"
)

func main() {
	fmt.Println("=== Network Scanner v0.7 ===")
	
	// Command-line flags
	saveFile := flag.String("save", "scan-latest.json", "File to save scan results")
	compareFile := flag.String("compare", "", "Previous scan file to compare against")
	flag.Parse()
	
	targetNetwork := "192.168.1.0/24"
	
	fmt.Printf("\nScanning network: %s\n", targetNetwork)
	fmt.Println("Scanning hosts and ports - this may take 1-2 minutes...\n")
	
	hosts := scanner.ScanNetwork(targetNetwork)
	
	// Lookup vendors
	for i := range hosts {
		hosts[i].Vendor = lookup.GetVendor(hosts[i].MAC)
	}
	
	// Display results
	fmt.Printf("\n=== Found %d active hosts ===\n", len(hosts))
	fmt.Printf("%-15s %-20s %-25s %-20s\n", "IP Address", "MAC Address", "Vendor", "Open Ports")
	fmt.Println(strings.Repeat("-", 85))
	for _, host := range hosts {
		ports := formatPorts(host.OpenPorts)
		fmt.Printf("%-15s %-20s %-25s %-20s\n", host.IP, host.MAC, host.Vendor, ports)
	}
	
	// Save scan results
	err := storage.SaveScan(targetNetwork, hosts, *saveFile)
	if err != nil {
		fmt.Printf("\nWarning: Failed to save scan: %v\n", err)
	} else {
		fmt.Printf("\n✓ Scan results saved to: %s\n", *saveFile)
	}
	
	// Compare with previous scan if requested
	if *compareFile != "" {
		fmt.Printf("\nComparing with previous scan: %s\n", *compareFile)
		previousScan, err := storage.LoadScan(*compareFile)
		if err != nil {
			fmt.Printf("Error loading previous scan: %v\n", err)
			return
		}
		
		added, removed := storage.CompareScans(previousScan.Hosts, hosts)
		
		if len(added) > 0 {
			fmt.Printf("\n🆕 NEW HOSTS DETECTED (%d):\n", len(added))
			for _, host := range added {
				fmt.Printf("  • %s (%s) - %s\n", host.IP, host.MAC, host.Vendor)
			}
		}
		
		if len(removed) > 0 {
			fmt.Printf("\n❌ HOSTS OFFLINE (%d):\n", len(removed))
			for _, host := range removed {
				fmt.Printf("  • %s (%s) - %s\n", host.IP, host.MAC, host.Vendor)
			}
		}
		
		if len(added) == 0 && len(removed) == 0 {
			fmt.Println("\n✓ No changes detected - network is stable")
		}
	}
}

func formatPorts(ports []int) string {
	if len(ports) == 0 {
		return "None"
	}
	
	var portStrs []string
	for _, p := range ports {
		portStrs = append(portStrs, fmt.Sprintf("%d", p))
	}
	result := strings.Join(portStrs, ",")
	
	if len(result) > 18 {
		result = result[:18] + "..."
	}
	
	return result
}

