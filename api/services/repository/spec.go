package repository

import "time"

type DbRepository struct {
	ID                                 int       `json:"id"`
	UserID                             int       `json:"user_id"`
	Size                               int       `json:"size"`
	Issues                             int       `json:"issues"`
	Stars                              int       `json:"stars"`
	Forks                              int       `json:"forks"`
	Watchers                           int       `json:"watchers"`
	Subscribers                        int       `json:"subscribers"`
	Network                            int       `json:"network"`
	IsFork                             bool      `json:"is_fork"`
	IsPrivate                          bool      `json:"is_private"`
	IsTemplate                         bool      `json:"is_template"`
	IsArchived                         bool      `json:"is_archived"`
	IsDisabled                         bool      `json:"is_disabled"`
	HasIssues                          bool      `json:"has_issues"`
	HasProjects                        bool      `json:"has_projects"`
	HasDownloads                       bool      `json:"has_downloads"`
	HasWiki                            bool      `json:"has_wiki"`
	HasPages                           bool      `json:"has_pages"`
	HasDiscussions                     bool      `json:"has_discussions"`
	CreatedAt                          time.Time `json:"created_at"`
	UpdatedAt                          time.Time `json:"updated_at"`
	PushedAt                           time.Time `json:"pushed_at"`
	DefaultBranch                      string    `json:"default_branch"`
	NodeID                             string    `json:"node_id"`
	GitURL                             string    `json:"git_url"`
	SSHURL                             string    `json:"ssh_url"`
	CloneURL                           string    `json:"clone_url"`
	SvnURL                             string    `json:"svn_url"`
	MirrorURL                          string    `json:"mirror_url"`
	Name                               string    `json:"name"`
	FullName                           string    `json:"full_name"`
	Description                        string    `json:"description"`
	Language                           string    `json:"language"`
	License                            string    `json:"license"`
	URL                                string    `json:"url"`
	Homepage                           string    `json:"homepage"`
	Topics                             []string  `json:"topics"`
	OSSFScorecardTotalScore            float64   `json:"ossf_scorecard_total_score"`
	OSSFScorecardDependencyUpdateScore float64   `json:"ossf_scorecard_dependency_update_score"`
	OSSFScorecardFuzzingScore          float64   `json:"ossf_scorecard_fuzzing_score"`
	OSSFScorecardMaintainedScore       float64   `json:"ossf_scorecard_maintained_score"`
	OSSFScorecardUpdatedAt             time.Time `json:"ossf_scorecard_updated_at"`
	OpenIssuesCount                    int       `json:"open_issues_count"`
	ClosedIssuesCount                  int       `json:"closed_issues_count"`
	IssuesVelocityCount                float64   `json:"issues_velocity_count"`
	OpenPRsCount                       int       `json:"open_prs_count"`
	ClosedPRsCount                     int       `json:"closed_prs_count"`
	MergedPRsCount                     int       `json:"merged_prs_count"`
	DraftPRsCount                      int       `json:"draft_prs_count"`
	SpamPRsCount                       int       `json:"spam_prs_count"`
	PRVelocityCount                    float64   `json:"pr_velocity_count"`
	ForkVelocity                       float64   `json:"fork_velocity"`
	PRActiveCount                      int       `json:"pr_active_count"`
	ActivityRatio                      float64   `json:"activity_ratio"`
	ContributorConfidence              float64   `json:"contributor_confidence"`
	Health                             float64   `json:"health"`
	LastPushedAt                       time.Time `json:"last_pushed_at"`
	LastMainPushedAt                   time.Time `json:"last_main_pushed_at"`
}

// DbContributorInfo represents the structure of a single contributor
type DbContributorInfo struct {
	ID                 int       `json:"id"`
	Login              string    `json:"login"`
	AvatarURL          string    `json:"avatar_url"`
	Company            string    `json:"company"`
	Location           string    `json:"location"`
	OSCR               float64   `json:"oscr"`
	Repos              []string  `json:"repos"`
	Tags               []string  `json:"tags"`
	Commits            int       `json:"commits"`
	PRsCreated         int       `json:"prs_created"`
	PRsReviewed        int       `json:"prs_reviewed"`
	IssuesCreated      int       `json:"issues_created"`
	CommitComments     int       `json:"commit_comments"`
	IssueComments      int       `json:"issue_comments"`
	PRReviewComments   int       `json:"pr_review_comments"`
	Comments           int       `json:"comments"`
	TotalContributions int       `json:"total_contributions"`
	LastContributed    time.Time `json:"last_contributed"`
	DevstatsUpdatedAt  time.Time `json:"devstats_updated_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// ContributorsResponse represents the structure of the contributors endpoint response
type ContributorsResponse struct {
	Data []DbContributorInfo `json:"data"`
	Meta struct {
		Page            int  `json:"page"`
		Limit           int  `json:"limit"`
		ItemCount       int  `json:"itemCount"`
		PageCount       int  `json:"pageCount"`
		HasPreviousPage bool `json:"hasPreviousPage"`
		HasNextPage     bool `json:"hasNextPage"`
	} `json:"meta"`
}
