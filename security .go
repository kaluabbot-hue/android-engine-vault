package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
)

const DeviceSalt = "ALPHA_ENGINE_VAULT_ROUTER_SECURE_2026_MX"

func GetHardwareFingerprint() string {
	var hardwarePayload strings.Builder

	if cpuData, err := ioutil.ReadFile("/proc/cpuinfo"); err == nil {
		lines := strings.Split(string(cpuData), "\n")
		for _, line := range lines {
			if strings.Contains(line, "Serial") || strings.Contains(line, "Hardware") || strings.Contains(line, "Revision") {
				hardwarePayload.WriteString(line)
			}
		}
	}

	if cidData, err := ioutil.ReadFile("/sys/block/mmcblk0/device/cid"); err == nil {
		hardwarePayload.Write(cidData)
	} else {
		if ufsData, err := ioutil.ReadFile("/sys/class/block/sda/device/model"); err == nil {
			hardwarePayload.Write(ufsData)
		}
	}

	interfaces, err := net.Interfaces()
	if err == nil {
		for _, iface := range interfaces {
			if iface.HardwareAddr != nil && iface.Name != "lo" {
				hardwarePayload.WriteString(iface.HardwareAddr.String())
			}
		}
	}

	hasher := sha256.New()
	hasher.Write([]byte(hardwarePayload.String() + DeviceSalt))
	
	return hex.EncodeToString(hasher.Sum(nil))
}

func VerifyDeviceAccess(authorizedFingerprint string) bool {
	currentDeviceFingerprint := GetHardwareFingerprint()
	return currentDeviceFingerprint == authorizedFingerprint
}

func main() {
	myFingerprint := GetHardwareFingerprint()
	
	fmt.Println("==================================================")
	fmt.Println("         ENGINE VAULT DEVICE IDENTIFIER           ")
	fmt.Println("==================================================")
	fmt.Printf("Your Unique Device Fingerprint:\n%s\n\n", myFingerprint)

	var AuthorizedSystemKey = myFingerprint 

	if !VerifyDeviceAccess(AuthorizedSystemKey) {
		fmt.Println("[CRITICAL SECURITY BREACH]")
		fmt.Println("Error: Hardware Binding Violation.")
		fmt.Println("This binary is locked to an authorized device profile.")
		fmt.Println("Unauthorised sharing detected. Execution terminated.")
		os.Exit(1)
	}

	fmt.Println("[SUCCESS] Hardware verified. Launching engine environment...")
}
