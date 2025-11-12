package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
	
	"network-scanner/pkg/lookup"

	"network-scanner/pkg/scanner"
	"network-scanner/pkg/storage"
)

type Server struct {
	network      string
	lastScan     *storage.ScanResult
	scanMutex    sync.RWMutex
	isScanning   bool
}

func NewServer(network string) *Server {
	return &Server{
		network: network,
	}
}

func (s *Server) Start(port int) error {
	http.HandleFunc("/", s.handleRoot)
	http.HandleFunc("/scan", s.handleScan)
	http.HandleFunc("/hosts", s.handleHosts)
	http.HandleFunc("/status", s.handleStatus)
	
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("\n🌐 API Server starting on http://localhost%s\n", addr)
	fmt.Println("Available endpoints:")
	fmt.Println("  GET  /          - API information")
	fmt.Println("  POST /scan      - Trigger new scan")
	fmt.Println("  GET  /hosts     - Get current hosts")
	fmt.Println("  GET  /status    - Get scanner status")
	fmt.Println()
	
	return http.ListenAndServe(addr, nil)
}

func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	info := map[string]interface{}{
		"name":    "Network Scanner API",
		"version": "0.8",
		"network": s.network,
		"endpoints": []string{
			"/scan (POST) - Trigger new scan",
			"/hosts (GET) - Get current hosts",
			"/status (GET) - Get scanner status",
		},
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(info)
}

func (s *Server) handleScan(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	s.scanMutex.Lock()
	if s.isScanning {
		s.scanMutex.Unlock()
		http.Error(w, "Scan already in progress", http.StatusConflict)
		return
	}
	s.isScanning = true
	s.scanMutex.Unlock()
	
	go func() {
		defer func() {
			s.scanMutex.Lock()
			s.isScanning = false
			s.scanMutex.Unlock()
		}()
		
		hosts := scanner.ScanNetwork(s.network)
		
		// Lookup vendors
		for i := range hosts {
			hosts[i].Vendor = lookup.GetVendor(hosts[i].MAC)
		}
		
		result := &storage.ScanResult{
			Timestamp: time.Now(),
			Network:   s.network,
			Hosts:     hosts,
		}
		
		s.scanMutex.Lock()
		s.lastScan = result
		s.scanMutex.Unlock()
	}()
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "Scan started",
		"message": "Scan running in background. Check /status for progress.",
	})
}

func (s *Server) handleHosts(w http.ResponseWriter, r *http.Request) {
	s.scanMutex.RLock()
	defer s.scanMutex.RUnlock()
	
	if s.lastScan == nil {
		http.Error(w, "No scan results available. Trigger a scan first with POST /scan", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.lastScan)
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	s.scanMutex.RLock()
	defer s.scanMutex.RUnlock()
	
	status := map[string]interface{}{
		"scanning": s.isScanning,
		"network":  s.network,
	}
	
	if s.lastScan != nil {
		status["last_scan"] = s.lastScan.Timestamp
		status["host_count"] = len(s.lastScan.Hosts)
	} else {
		status["last_scan"] = nil
		status["host_count"] = 0
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

