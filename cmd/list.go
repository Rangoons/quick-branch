/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Khan/genqlient/graphql"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/rangoons/quick-branch/internal/generated"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Linear issues based on your saved filters",
	RunE: func(cmd *cobra.Command, args []string) error {
		// teamName := viper.GetString("list.team_name")
		// assigneeFilter := viper.GetString("list.assignee_filter")
		// fmt.Printf("Fetching %s issues for team \"%s\"...\n\n", assigneeFilter, teamName)

		resp, err := fetchIssues()
		if err != nil {
			return err
		}

		if len(resp.Issues.Nodes) == 0 {
			fmt.Println("No issues found.")
			return nil
		}
		var (
			header = lipgloss.Color("#957FB8")
			border = lipgloss.Color("#54546D")
			text   = lipgloss.Color("#DCD7BA")

			headerStyle = lipgloss.NewStyle().Foreground(header).Bold(true).Align(lipgloss.Center)
			cellStyle   = lipgloss.NewStyle().Padding(0, 1)
			oddRowStyle = cellStyle.Foreground(text)
		)
		t := table.New().Border(lipgloss.NormalBorder()).BorderStyle(lipgloss.NewStyle().Foreground(border)).StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return headerStyle
			default:
				return oddRowStyle
			}
		}).Headers("PRIORITY", "TICKET", "TITLE", "STATE")

		for _, issue := range resp.Issues.Nodes {
			t.Row(strconv.FormatFloat(issue.Priority, 'f', 0, 64), issue.Identifier, issue.Title, issue.State.Name)
			// fmt.Printf("(%g) [%s] %s  [%s]\n", issue.Priority, issue.Identifier, issue.Title, issue.State.Name)
		}
		fmt.Println(t)
		return nil
	},
}

var listSetupCmd = &cobra.Command{
	Use:   "setup",
	Short: "Configure the filters used by the list command",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSetupWizard()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.AddCommand(listSetupCmd)
}

func newGraphQLClient() graphql.Client {
	httpClient := &http.Client{
		Transport: &authorizedTransport{
			apiKey: viper.GetString("api_key"),
			base:   http.DefaultTransport,
		},
	}
	return graphql.NewClient("https://api.linear.app/graphql", httpClient)
}

func runSetupWizard() error {
	if viper.GetString("api_key") == "" {
		return fmt.Errorf("no API key found. Please run 'quick-branch auth' first")
	}

	ctx := context.Background()
	client := newGraphQLClient()

	// Step 1: fetch teams and pick team + assignee filter
	teamsResp, err := generated.ViewerTeams(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to fetch teams: %w", err)
	}
	teams := teamsResp.Viewer.Teams.Nodes
	if len(teams) == 0 {
		return fmt.Errorf("no teams found for your account")
	}

	teamOpts := make([]huh.Option[string], len(teams))
	for i, t := range teams {
		teamOpts[i] = huh.NewOption(t.Name, t.Id)
	}

	var selectedTeamID, assigneeFilter string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select your team").
				Options(teamOpts...).
				Value(&selectedTeamID),
			huh.NewSelect[string]().
				Title("Show which issues").
				Options(
					huh.NewOption("Assigned to me", "me"),
					huh.NewOption("Unassigned", "unassigned"),
					huh.NewOption("All", "all"),
				).
				Value(&assigneeFilter),
		),
	).Run()
	if err != nil {
		return err
	}

	// Step 2: fetch states for the chosen team and pick which to include
	statesResp, err := generated.TeamStatesById(ctx, client, selectedTeamID)
	if err != nil {
		return fmt.Errorf("failed to fetch states: %w", err)
	}
	allStates := statesResp.Team.States.Nodes

	stateOpts := make([]huh.Option[string], len(allStates))
	for i, s := range allStates {
		stateOpts[i] = huh.NewOption(s.Name, s.Id)
	}

	var selectedStateIDs []string
	err = huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select states to include").
				Options(stateOpts...).
				Value(&selectedStateIDs),
		),
	).Run()
	if err != nil {
		return err
	}

	// Look up team name and selected state names for the confirmation message
	var teamName string
	for _, t := range teams {
		if t.Id == selectedTeamID {
			teamName = t.Name
			break
		}
	}
	selectedStateNames := make([]string, 0, len(selectedStateIDs))
	stateIDSet := make(map[string]bool, len(selectedStateIDs))
	for _, id := range selectedStateIDs {
		stateIDSet[id] = true
	}
	for _, s := range allStates {
		if stateIDSet[s.Id] {
			selectedStateNames = append(selectedStateNames, s.Name)
		}
	}

	viper.Set("list.team_id", selectedTeamID)
	viper.Set("list.team_name", teamName)
	viper.Set("list.assignee_filter", assigneeFilter)
	viper.Set("list.state_ids", selectedStateIDs)

	if err := saveConfig(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("\nSaved! `list` will show %s issues on team \"%s\" in states: %s\n",
		assigneeFilter, teamName, strings.Join(selectedStateNames, ", "))
	return nil
}

func fetchIssues() (*generated.FilteredIssuesResponse, error) {
	if viper.GetString("api_key") == "" {
		return nil, fmt.Errorf("no API key found. Please run 'quick-branch auth' first")
	}
	teamID := viper.GetString("list.team_id")
	if teamID == "" {
		return nil, fmt.Errorf("list not configured. Please run 'quick-branch list setup' first")
	}

	filter := buildIssueFilter(
		teamID,
		viper.GetStringSlice("list.state_ids"),
		viper.GetString("list.assignee_filter"),
	)

	ctx := context.Background()
	return generated.FilteredIssues(ctx, newGraphQLClient(), filter)
}

func buildIssueFilter(teamID string, stateIDs []string, assigneeFilter string) *generated.IssueFilter {
	filter := &generated.IssueFilter{
		Team: &generated.TeamFilter{
			Id: &generated.IDComparator{Eq: &teamID},
		},
		State: &generated.WorkflowStateFilter{
			Id: &generated.IDComparator{In: stateIDs},
		},
	}

	switch assigneeFilter {
	case "me":
		t := true
		filter.Assignee = &generated.NullableUserFilter{
			IsMe: &generated.BooleanComparator{Eq: &t},
		}
	case "unassigned":
		t := true
		filter.Assignee = &generated.NullableUserFilter{Null: &t}
	}

	return filter
}

func saveConfig() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}
	appConfigDir := configDir + "/quick-branch"
	if err := os.MkdirAll(appConfigDir, 0o700); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}
	viper.SetConfigFile(appConfigDir + "/config.yaml")
	return viper.WriteConfig()
}
