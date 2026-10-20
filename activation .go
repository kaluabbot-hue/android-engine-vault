package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// ClearTerminal wipes the screen clean for the security prompt
func ClearTerminal() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// ExecuteTerminalLockdown runs Step 10: Prompts user for their reference ID
func ExecuteTerminalLockdown() string {
	reader := bufio.NewReader(os.Stdin)

	for {
		ClearTerminal()
		fmt.Println("================================================================")
		fmt.Println("🛑 [SECURITY ALERT] RUNTIME VALIDATION EXPIRED OR FAULT DETECTED")
		fmt.Println("================================================================")
		fmt.Println(" This premium lag-elimination engine requires an active license.")
		fmt.Println(" Please input your Paystack Transaction Reference ID below")
		fmt.Println(" to instantly unlock your device optimization runtime.")
		fmt.Println("================================================================")
		fmt.Print("\nEnter your Paystack Transaction Reference ID to activate: ")

		input, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		referenceID := strings.TrimSpace(input)

		if referenceID == "" {
			fmt.Println("\n❌ [ERROR] Reference ID cannot be empty. Press Enter to retry...")
			reader.ReadString('\n')
			continue
		}

		return referenceID
	}
}
