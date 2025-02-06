package main

import (
	"log"
	"net"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) findPossibleIPAddresses() {
	interfaces, err := net.Interfaces()
	if err != nil {
		a.AlertDialog("Error finding IP addresses", err.Error())
		return
	}

	possibleAddresses := make([]string, 0)

	// Iterate through the interfaces to find the first non-loopback address
	for _, iface := range interfaces {
		// Skip loopback interfaces (127.0.0.1, etc.)
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			log.Fatal(err)
		}

		// Iterate over the addresses and look for an IPv4 address
		for _, addr := range addrs {
			// Check if the address is an IPv4 address
			if ipNet, ok := addr.(*net.IPNet); ok && ipNet.IP.To4() != nil {
				possibleAddresses = append(possibleAddresses, ipNet.IP.String())
			}
		}
	}

	a.sacnConfig.PossibleIpAddresses = possibleAddresses
	if len(possibleAddresses) == 0 {
		a.AlertDialog("No IP addresses found", "This should not happen. Make sure you have a network interface active.")
		return
	}
	a.sacnConfig.IpAddress = possibleAddresses[0]
}

func (a *App) LoadFile() string {
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Load Följe Configuration",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Följe Configurations (*.fconf)",
				Pattern:     "*.fconf",
			},
		}})
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "{}"
	}

	data, err := os.ReadFile(file)
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return "{}"
	}

	return string(data)
}

func (a *App) SaveFile(content string) {
	file, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Load Följe Configuration",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Följe Configurations (*.fconf)",
				Pattern:     "*.fconf",
			},
		},
		DefaultFilename: "conf.fconf",
	})
	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}

	err = os.WriteFile(file, []byte(content), 0644)

	if err != nil {
		runtime.LogError(a.ctx, err.Error())
		return
	}
}
