package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/rangoons/quick-branch/internal/generated"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

// authCmd represents the auth command
var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Store your Linear API key",
	Long: `Store your Linear API key in the config file for future use.
The API key will be hidden while you type or paste it.

Example:
  linear-cli auth`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Print("Enter your Linear API key: ")

		apiKeyBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println() // Print newline after hidden input

		if err != nil {
			return fmt.Errorf("failed to read API key: %w", err)
		}

		apiKey := strings.TrimSpace(string(apiKeyBytes))
		if apiKey == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		// Test the API key before saving
		fmt.Println("Verifying API key...")
		if err := verifyAPIKey(apiKey); err != nil {
			return fmt.Errorf("API key verification failed: %w", err)
		}

		return writeAPIToken(apiKey)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
}

// authorizedTransport adds the Authorization header to all requests
type authorizedTransport struct {
	apiKey string
	base   http.RoundTripper
}

func (t *authorizedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", t.apiKey)
	return t.base.RoundTrip(req)
}

func verifyAPIKey(apiKey string) error {
	httpClient := &http.Client{
		Transport: &authorizedTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.linear.app/graphql", httpClient)

	ctx := context.Background()
	response, err := generated.Me(ctx, graphqlClient)
	if err != nil {
		return err
	}

	fmt.Println("✓ API key verified successfully!")
	fmt.Printf("\nAuthenticated as:\n")
	fmt.Printf("  Name:  %s\n", response.Viewer.Name)
	fmt.Printf("  Email: %s\n\n", response.Viewer.Email)

	return nil
}

func writeAPIToken(apiKey string) error {
	viper.Set("api_key", apiKey)

	// Ensure the config directory exists
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	appConfigDir := configDir + "/quick-branch"
	if err := os.MkdirAll(appConfigDir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := appConfigDir + "/config.yaml"
	viper.SetConfigFile(configPath)

	err = viper.SafeWriteConfig()
	if err != nil {
		// If config doesn't exist, create it
		err = viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}
	fmt.Println("✓ API key saved to", configPath)
	return nil
}
