package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	
	"network-scanner/pkg/models"
)

type ScanResult struct {
	Timestamp time.Time     `json:"timestamp"`
	Network   string        `json:"network"`
	Hosts     []models.Host `json:"hosts"`
}

func SaveScan(network string, hosts []models.Host, filename string) error {
	result := ScanResult{
		Timestamp: time.Now(),
		Network:   network,
		Hosts:     hosts,
	}
	
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	
	return nil
}

func LoadScan(filename string) (*ScanResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	
	var result ScanResult
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	
	return &result, nil
}

func CompareScans(oldHosts, newHosts []models.Host) (added, removed []models.Host) {
	oldIPs := make(map[string]models.Host)
	newIPs := make(map[string]models.Host)
	
	for _, h := range oldHosts {
		oldIPs[h.IP] = h
	}
	
	for _, h := range newHosts {
		newIPs[h.IP] = h
	}
	
	// Find new hosts
	for ip, host := range newIPs {
		if _, exists := oldIPs[ip]; !exists {
			added = append(added, host)
		}
	}
	
	// Find removed hosts
	for ip, host := range oldIPs {
		if _, exists := newIPs[ip]; !exists {
			removed = append(removed, host)
		}
	}
	
	return added, removed
}

