package commands

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/bluefunda/abaper-cli/internal/client"
	"github.com/bluefunda/abaper-cli/internal/config"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with ABAPer using device authorization flow",
	Long: `Initiates an OAuth2 device authorization flow.
A browser window will open for you to authenticate.
Credentials are stored locally in ~/.abaper/tokens.yaml.`,
	RunE: runLogin,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Clear stored credentials",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ClearTokens(); err != nil {
			return fmt.Errorf("logout: %w", err)
		}
		fmt.Println("Logged out successfully.")
		return nil
	},
}

func runLogin(cmd *cobra.Command, args []string) error {
	cfg := config.Load()
	realm := cfg.Realm

	// Step 1: Request device code
	fmt.Println("Requesting device authorization...")
	deviceResp, err := client.RequestDeviceCode(realm)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Step 2: Open browser
	verifyURL := deviceResp.VerificationURIComplete
	if verifyURL == "" {
		verifyURL = deviceResp.VerificationURI
	}

	fmt.Printf("\nOpen this URL in your browser to log in:\n  %s\n\n", verifyURL)
	fmt.Printf("Your code: %s\n\n", deviceResp.UserCode)

	_ = openBrowser(verifyURL)

	// Step 3: Poll for token
	fmt.Println("Waiting for authorization...")
	tokenResp, err := client.PollForToken(realm, deviceResp.DeviceCode, deviceResp.Interval)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Step 4: Save tokens
	tokens := &config.Tokens{
		AccessToken:  tokenResp.AccessToken,
		RefreshToken: tokenResp.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).UnixMilli(),
	}

	if err := config.SaveTokens(tokens); err != nil {
		return fmt.Errorf("save credentials: %w", err)
	}

	fmt.Println("Successfully logged in!")
	return nil
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}
