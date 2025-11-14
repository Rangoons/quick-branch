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
		// Prompt for API key
		fmt.Print("Enter your Linear API key: ")

		// Read the API key securely (hidden input)
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

		// Save the API key
		return writeAPIToken(apiKey)
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// authCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// authCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
	// Create HTTP client with authorization
	httpClient := &http.Client{
		Transport: &authorizedTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	// Create GraphQL client
	graphqlClient := graphql.NewClient("https://api.linear.app/graphql", httpClient)

	// Make the request
	ctx := context.Background()
	response, err := generated.Me(ctx, graphqlClient)
	if err != nil {
		return err
	}

	// Display verification results
	fmt.Println("✓ API key verified successfully!")
	fmt.Printf("\nAuthenticated as:\n")
	fmt.Printf("  Name:  %s\n", response.Viewer.Name)
	fmt.Printf("  Email: %s\n\n", response.Viewer.Email)

	return nil
}

func writeAPIToken(apiKey string) error {
	viper.Set("api_key", apiKey)

	// Try to write the config file
	err := viper.SafeWriteConfig()
	if err != nil {
		// If config doesn't exist, create it
		err = viper.WriteConfig()
		if err != nil {
			return fmt.Errorf("failed to write config: %w", err)
		}
	}
	configPath := viper.ConfigFileUsed()
	if configPath != "" {
		fmt.Println("✓ API key saved to", configPath)
	} else {
		fmt.Println("✓ API key saved to config")
	}
	return nil
}
