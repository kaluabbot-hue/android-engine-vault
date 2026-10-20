package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// PaystackResponse structures the data returned from Paystack's validation servers
type PaystackResponse struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Status    string `json:"status"`
		Reference string `json:"reference"`
		Amount    int    `json:"amount"` 
		Currency  string `json:"currency"`
	} `json:"data"`
}

const (
	PaystackBaseURL    = "https://paystack.co"
	MonthlyAmountKobo  = 50000  // ₦500
	LifetimeAmountKobo = 150000 // ₦1,500
)

// VerifyPaystackPayment securely pings the server to validate the reference ID
func VerifyPaystackPayment(reference string) (bool, string) {
	// SECURE FIX: Grabs the key dynamically from the system memory during runtime
	secretKey := os.Getenv("PAYSTACK_SECRET_KEY")
	if secretKey == "" {
		return false, "MISSING_CONFIGURATION"
	}

	url := fmt.Sprintf("%s%s", PaystackBaseURL, strings.TrimSpace(reference))
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, "REQUEST_ERR"
	}

	req.Header.Set("Authorization", "Bearer "+secretKey)
	req.Header.Set("Cache-Control", "no-cache")

	client := &http.Client{Timeout: 12 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, "NETWORK_TIMEOUT"
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "READ_ERR"
	}

	var res PaystackResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return false, "PARSE_ERR"
	}

	if res.Status && res.Data.Status == "success" && res.Data.Currency == "NGN" {
		if res.Data.Amount >= LifetimeAmountKobo {
			return true, "LIFETIME"
		} else if res.Data.Amount >= MonthlyAmountKobo {
			return true, "MONTHLY"
		}
	}

	return false, "INVALID_TRANSACTION"
}
