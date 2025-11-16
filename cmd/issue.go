package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/Khan/genqlient/graphql"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/rangoons/quick-branch/internal/generated"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	url         bool
	branch      bool
	checkout    bool
	description bool
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue [issue number]",
	Short: "Fetches linear issue",
	Long: `issue is a CLI tool for quickly working with linear

	You can provide an issue number as arguments & copy the issue URL or branch name & create a new branch with that branch name
	`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueID := args[0]
		issue, err := fetchIssue(issueID)
		if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		if description {
			// Style the title (bold)
			titleStyle := lipgloss.NewStyle().Bold(true)

			// Style the state with its color from Linear
			stateStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(issue.State.Color))

			fmt.Printf("%s: %s\n\n",
				titleStyle.Render(issue.Title),
				stateStyle.Render(issue.State.Name))

			// Render the markdown description prettily
			renderer, err := glamour.NewTermRenderer(
				glamour.WithAutoStyle(),
			)
			if err == nil {
				out, err := renderer.Render(*issue.Description)
				if err == nil {
					fmt.Print(out)
				} else {
					// Fallback to plain text if rendering fails
					fmt.Println(issue.Description)
				}
			} else {
				// Fallback to plain text if renderer creation fails
				fmt.Println(issue.Description)
			}
		}
		if url {
			err := clipboard.WriteAll(issue.Url)
			if err == nil {
				fmt.Println("Copied issue url to clipboard")
			}
		} else if branch {
			err := clipboard.WriteAll(issue.BranchName)
			if err == nil {
				fmt.Println("Copied branch name to clipboard")
			}
		}
		if checkout {
			if err := checkoutBranch(issue.BranchName); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(issueCmd)
	issueCmd.Flags().BoolVarP(&url, "url", "u", false, "Copies the issue URL to your clipboard")
	issueCmd.Flags().BoolVarP(&branch, "branch", "b", false, "Copies the branch name to your clipboard")
	issueCmd.Flags().BoolVarP(&checkout, "checkout", "c", false, "Creates a new branch in the cwd using the branch name from linear")
	issueCmd.Flags().BoolVarP(&description, "verbose", "v", false, "Prints the issue description")
}

func fetchIssue(issueID string) (*generated.IssueIssue, error) {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return nil, fmt.Errorf("no API key found. Please run 'quick-branch auth' first")
	}
	httpClient := &http.Client{
		Transport: &authorizedTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.linear.app/graphql", httpClient)

	ctx := context.Background()
	response, err := generated.Issue(ctx, graphqlClient, issueID)
	if err != nil {
		return nil, err
	}

	return &response.Issue, nil
}

func checkoutBranch(branchName string) error {
	err := exec.Command("git", "switch", "-c", branchName).Run()
	if err != nil {
		return err
	}
	fmt.Printf("Success! Now working on %v\n", branchName)
	return nil
}
