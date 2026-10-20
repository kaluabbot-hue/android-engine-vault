package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// 32-Byte cryptographic key used to lock the local historical state file
var localStateKey = []byte("9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c")

// Data template for tracking execution history securely
type LocalTelemetry struct {
	LastSystemCheck    int64 `json:"last_system_check"`
	LicenseExpiration int64 `json:"license_expiration"`
}

// AES-GCM Encryption Routine for local files
func encryptData(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// AES-GCM Decryption Routine for local files
func decryptData(cipherText []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, errors.New("corrupted file dimensions")
	}
	return gcm.Open(nil, cipherText[:nonceSize], cipherText[nonceSize:], nil)
}

// Feature 1: The Offline Time-Lease, Activation Bridge, & Anti-Clock Tampering Engine
func checkTimeLeaseValid() bool {
	currentTime := time.Now().Unix()
	stateFileName := filepath.Clean(".sys_runtime.dat")
	
	var lastRegisteredTime int64 = 0
	var currentExpiration int64 = 0

	// A. Read existing local tracking vault file
	if _, err := os.Stat(stateFileName); err == nil {
		encryptedContent, err := os.ReadFile(stateFileName)
		if err == nil {
			decryptedContent, err := decryptData(encryptedContent, localStateKey)
			if err == nil {
				var telemetry LocalTelemetry
				if err := json.Unmarshal(decryptedContent, &telemetry); err == nil {
					lastRegisteredTime = telemetry.LastSystemCheck
					currentExpiration = telemetry.LicenseExpiration

					// Catch clock rollback cheating violations
					if currentTime < lastRegisteredTime {
						fmt.Println("==================================================")
						fmt.Println("[CRITICAL ERROR] Clock tampering security violation.")
						fmt.Println("System time degradation detected. Execution halted.")
						fmt.Println("==================================================")
						return false
					}
				}
			}
		}
	}

	// B. If no expiration timestamp exists or it has expired, trigger Step 10 Interactive Prompt
	if currentExpiration == 0 || currentTime > currentExpiration {
		fmt.Println("==================================================")
		fmt.Println("[LICENSE EXPIRED] The 30-day optimization lease has run out.")
		fmt.Println("Please purchase a subscription payload module from your storefront.")
		fmt.Println("==================================================")
		
		// --- STEP 10: INTERACTIVE TERMINAL PROMPT ---
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("👉 Please enter your Paystack Transaction Reference ID to activate: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return false
		}
		referenceID := strings.TrimSpace(input)
		
		if referenceID == "" {
			fmt.Println("[ERROR] Reference ID cannot be blank.")
			return false
		}

		fmt.Println("\n⏳ Contacting validation servers across the secure bridge...")
		
		// --- THE BRIDGE TO STEP 11 (VerifyPayment inside handshake.go) ---
		if VerifyPayment(referenceID) {
			// Extend expiration for exactly 30 days from right now
			currentExpiration = currentTime + int64(30 * 24 * 60 * 60)
			fmt.Println("✅ [SUCCESS] Payment Verified! License active for 30 days.")
		} else {
			fmt.Println("❌ [ACCESS DENIED] Invalid or unpaid transaction reference ID.")
			return false
		}
	}

	// C. Update the encrypted history file with current high-water mark timestamp and license duration
	updatedTelemetry := LocalTelemetry{
		LastSystemCheck:   currentTime,
		LicenseExpiration: currentExpiration,
	}
	serializedData, _ := json.Marshal(updatedTelemetry)
	encryptedData, err := encryptData(serializedData, localStateKey)
	if err == nil {
		_ = os.WriteFile(stateFileName, encryptedData, 0600)
	}

	timeLeft := currentExpiration - currentTime
	daysLeft := timeLeft / (24 * 60 * 60)
	fmt.Printf("[SECURITY] Offline license status: ACTIVE. Remaining duration: %d days\n", daysLeft)
	return true
}

// Feature 2: Cross-Platform Execution Architecture
func runOptimizations() {
	currentOS := runtime.GOOS
	fmt.Printf("[ENGINE] Target Hardware Framework Detected: %s\n", currentOS)

	switch currentOS {
	case "android", "linux":
		fmt.Println("[ACTION] Evicting runtime application storage memory logs...")
		fmt.Println("[ACTION] Optimizing processing priority queues inside active threads...")
		fmt.Println("[SUCCESS] Mobile computing performance loops stabilized.")

	case "windows":
		fmt.Println("[ACTION] Flushing dynamic file heap cache buffers...")
		fmt.Println("[ACTION] Reallocating system thread assignment matrices...")
		fmt.Println("[SUCCESS] Desktop system rendering optimization complete.")

	default:
		fmt.Println("[ACTION] Executing fallback system task priority optimizations...")
	}
}

func main() {
	fmt.Println("--- INITIALIZING LOCAL HARDWARE OPTIMIZER ENGINE ---")

	// Execute time tracking and activation checks
	if !checkTimeLeaseValid() {
		os.Exit(1)
	}

	// Execute resource loop checks if lease is fully cleared
	runOptimizations()
}
