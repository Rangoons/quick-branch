/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Khan/genqlient/graphql"
	"github.com/rangoons/quick-branch/internal/generated"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	status       bool
	checkoutFlag bool
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <issueID>",
	Short: "start will assign you to the issue you pass in",
	Long:  ``,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		issueID := args[0]
		err := assignMe(issueID)
		if err != nil {
			fmt.Println(err)
		}
		if status {
			err := updateIssueStatus(issueID)
			if err != nil {
				fmt.Println(err)
			}
		}
		if checkoutFlag {
			issue, err := fetchIssue(issueID)
			if err != nil {
				fmt.Printf("Error fetching issue: %v\n", err)
				return
			}
			if err := checkoutBranch(issue.BranchName); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolVarP(&status, "status", "s", false, "Updates the status of the issue to 'In Dev'")
	startCmd.Flags().BoolVarP(&checkoutFlag, "checkout", "c", false, "Creates a new branch in the cwd using the branch name from linear")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// startCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func assignMe(issueID string) error {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return fmt.Errorf("no API key found. Please run 'quick-branch auth' first")
	}
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
	input := generated.IssueUpdateInput{AssigneeId: &response.Viewer.Id}
	mutation, err := generated.IssueUpdate(ctx, graphqlClient, issueID, input)
	if err != nil {
		return err
	}
	fmt.Printf("Success! Assigned %v to %v\n", response.Viewer.Name, mutation.IssueUpdate.Issue.Title)
	return nil
}

func updateIssueStatus(issueID string) error {
	apiKey := viper.GetString("api_key")
	if apiKey == "" {
		return fmt.Errorf("no API key found. Please run 'quick-branch auth' first")
	}
	httpClient := &http.Client{
		Transport: &authorizedTransport{
			apiKey: apiKey,
			base:   http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.linear.app/graphql", httpClient)

	ctx := context.Background()
	response, err := generated.TeamStates(ctx, graphqlClient, issueID)
	if err != nil {
		return err
	}
	var inDev string
	for _, s := range response.Issue.Team.States.Nodes {
		if s.Name == "In Dev" {
			inDev = s.Id
		}
	}

	input := generated.IssueUpdateInput{StateId: &inDev}
	mutation, err := generated.IssueUpdate(ctx, graphqlClient, issueID, input)
	if err != nil {
		return err
	}
	fmt.Printf("Success! Updated %v to %v\n", mutation.IssueUpdate.Issue.Title, mutation.IssueUpdate.Issue.State.Name)
	return nil
}
