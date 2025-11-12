package models

type Host struct {
	IP        string
	MAC       string
	Vendor    string
	OpenPorts []int
}

