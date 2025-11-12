package scanner

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"sync"
	"time"
	
	"network-scanner/pkg/models"
)

var CommonPorts = []int{
	21, 22, 23, 25, 80, 443, 445, 3389, 8080, 8443,
}

func ScanNetwork(cidr string) []models.Host {
	ip, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}
	
	var activeHosts []models.Host
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); incIP(ip) {
		wg.Add(1)
		go func(ipAddr string) {
			defer wg.Done()
			if PingHost(ipAddr) {
				mac := GetMACAddress(ipAddr)
				openPorts := ScanPorts(ipAddr, CommonPorts)
				
				mu.Lock()
				activeHosts = append(activeHosts, models.Host{
					IP:        ipAddr,
					MAC:       mac,
					OpenPorts: openPorts,
				})
				mu.Unlock()
			}
		}(ip.String())
	}
	
	wg.Wait()
	return activeHosts
}

func PingHost(ip string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, "ping", "-c", "1", "-W", "1", ip)
	return cmd.Run() == nil
}

func GetMACAddress(ip string) string {
	exec.Command("ping", "-c", "1", "-W", "1", ip).Run()
	
	out, err := exec.Command("ip", "neigh", "show", ip).Output()
	if err != nil {
		return "Unknown"
	}
	
	fields := strings.Fields(string(out))
	for i, field := range fields {
		if field == "lladdr" && i+1 < len(fields) {
			return fields[i+1]
		}
	}
	
	return "Unknown"
}

func ScanPorts(ip string, ports []int) []int {
	var openPorts []int
	var mu sync.Mutex
	var wg sync.WaitGroup
	
	for _, port := range ports {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", ip, p)
			conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
			if err == nil {
				conn.Close()
				mu.Lock()
				openPorts = append(openPorts, p)
				mu.Unlock()
			}
		}(port)
	}
	
	wg.Wait()
	return openPorts
}

func incIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

