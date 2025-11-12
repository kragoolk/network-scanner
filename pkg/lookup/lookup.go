package lookup

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func GetVendor(mac string) string {
	if mac == "Unknown" {
		return "N/A"
	}
	
	url := fmt.Sprintf("https://api.macvendors.com/%s", mac)
	
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return "Lookup Failed"
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != 200 {
		return "Unknown Vendor"
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "Lookup Failed"
	}
	
	vendor := strings.TrimSpace(string(body))
	if len(vendor) > 23 {
		vendor = vendor[:23]
	}
	
	return vendor
}

