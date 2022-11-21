package cmd

import (
	"ahkpm/src/utils"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

//go:embed search-long.md
var searchLong string

var searchCmd = &cobra.Command{
	Use:   "search <searchTerm>...",
	Short: "Searches GitHub for packages matching the specified query",
	Long:  searchLong,
	Run: func(cmd *cobra.Command, args []string) {
		rawQuery := strings.Join(args, " ") + " topic:ahkpm-package"
		fullUrl := "https://api.github.com/search/repositories?q=" + url.QueryEscape(rawQuery)
		resp, err := http.Get(fullUrl)
		if err != nil {
			utils.Exit("Error searching for packages. Unable to reach GitHub.")
		}
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			utils.Exit("Error reading response body")
		}

		var searchResponse SearchResponse
		err = json.Unmarshal(body, &searchResponse)
		if err != nil {
			utils.Exit("Error parsing response body")
		}

		fmt.Println(GetSearchResultsTable(searchResponse.Items))
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}

func GetSearchResultsTable(items []SearchResponseItem) string {
	lengthOfLongestName := 0
	lengthOfLongestDescription := 0
	for _, item := range items {
		if len(item.FullName) > lengthOfLongestName {
			lengthOfLongestName = len(item.FullName)
		}
		if len(item.Description) > lengthOfLongestDescription {
			lengthOfLongestDescription = len(item.Description)
		}
	}

	ghPrefix := "github.com/"

	maxNameLength := lengthOfLongestName + len(ghPrefix)
	nameHeader := utils.RightPad("Name", " ", maxNameLength)
	descriptionHeader := utils.RightPad("Description", " ", lengthOfLongestDescription)
	nameUnderline := strings.Repeat("-", maxNameLength)
	descriptionUnderline := strings.Repeat("-", lengthOfLongestDescription)

	var table strings.Builder
	table.WriteString(nameHeader + "\t" + descriptionHeader + "\n")
	table.WriteString(nameUnderline + "\t" + descriptionUnderline + "\n")
	for _, item := range items {
		table.WriteString(fmt.Sprintf(
			"%s\t%s\n",
			utils.RightPad(ghPrefix+item.FullName, " ", maxNameLength),
			utils.RightPad(item.Description, " ", lengthOfLongestDescription),
		))
	}
	return table.String()
}

type SearchResponse struct {
	TotalCount        int `json:"total_count"`
	IncompleteResults bool
	Items             []SearchResponseItem
}

type SearchResponseItem struct {
	Id              int
	Name            string
	FullName        string `json:"full_name"`
	Owner           SearchResponseItemOwner
	Private         bool
	HtmlUrl         string `json:"html_url"`
	Description     string
	Fork            bool
	Url             string
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
	PushedAt        string `json:"pushed_at"`
	Homepage        string
	Size            int
	StargazersCount int `json:"stargazers_count"`
	WatchersCount   int `json:"watchers_count"`
	Language        string
	ForksCount      int    `json:"forks_count"`
	OpenIssuesCount int    `json:"open_issues_count"`
	MasterBranch    string `json:"master_branch"`
	DefaultBranch   string `json:"default_branch"`
	Score           float64
}

type SearchResponseItemOwner struct {
	Login             string
	Id                int
	NodeId            string `json:"node_id"`
	AvatarUrl         string `json:"avatar_url"`
	GravatarId        string `json:"gravatar_id"`
	Url               string
	HtmlUrl           string `json:"html_url"`
	FollowersUrl      string `json:"followers_url"`
	FollowingUrl      string `json:"following_url"`
	GistsUrl          string `json:"gists_url"`
	StarredUrl        string `json:"starred_url"`
	SubscriptionsUrl  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposUrl          string `json:"repos_url"`
	EventsUrl         string `json:"events_url"`
	ReceivedEventsUrl string `json:"received_events_url"`
	Type              string
	SiteAdmin         bool `json:"site_admin"`
}
