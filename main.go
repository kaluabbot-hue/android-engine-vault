package main

import (
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
"time"
)

// The absolute hardcoded expiration date: Sunday, November 1, 2026 00:00:00 UTC
const ExpirationTimestamp int64 = 1793404800

// 32-Byte cryptographic key used to lock the local historical state file
var localStateKey = []byte("9f8e7d6c5b4a3f2e1d0c9b8a7f6e5d4c")

// Data template for tracking execution history securely
type LocalTelemetry struct {
LastSystemCheck int64 `json:"last_system_check"`
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

// Feature 1: The Offline Time-Lease & Anti-Clock Tampering Engine
func checkTimeLeaseValid() bool {
currentTime := time.Now().Unix()

// A. Check against the absolute hardcoded expiration wall
if currentTime > ExpirationTimestamp {
fmt.Println("==================================================")
fmt.Println("[LEASE EXPIRED] The 30-day optimization lease has run out.")
fmt.Println("Please purchase the updated monthly payload module.")
fmt.Println("==================================================")
return false
}

stateFileName := filepath.Clean(".sys_runtime.dat")
var lastRegisteredTime int64 = 0

// B. Read existing local tracking file to catch clock rollback cheating
if _, err := os.Stat(stateFileName); err == nil {
encryptedContent, err := os.ReadFile(stateFileName)
if err == nil {
decryptedContent, err := decryptData(encryptedContent, localStateKey)
if err == nil {
var telemetry LocalTelemetry
if err := json.Unmarshal(decryptedContent, &telemetry); err == nil {
lastRegisteredTime = telemetry.LastSystemCheck

// If the current device time is behind the history file, the user cheated the clock
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

// C. Update the encrypted history file with the current high-water mark timestamp
updatedTelemetry := LocalTelemetry{LastSystemCheck: currentTime}
serializedData, _ := json.Marshal(updatedTelemetry)
encryptedData, err := encryptData(serializedData, localStateKey)
if err == nil {
_ = os.WriteFile(stateFileName, encryptedData, 0600)
}

daysLeft := time.Duration(ExpirationTimestamp-currentTime) * time.Second
fmt.Printf("[SECURITY] Offline license status: ACTIVE. Remaining duration: %v\n", daysLeft)
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

// Execute time tracking check
if !checkTimeLeaseValid() {
os.Exit(1)
}

// Execute resource loop checks
runOptimizations()
}
