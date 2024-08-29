package histogram

import "time"

type PrHistogramData struct {
	Bucket                    time.Time `json:"bucket"`
	PrCount                   int       `json:"prs_count"`
	AcceptedPrs               int       `json:"accepted_prs"`
	OpenPrs                   int       `json:"open_prs"`
	ClosedPrs                 int       `json:"closed_prs"`
	DraftPrs                  int       `json:"draft_prs"`
	ActivePrs                 int       `json:"active_prs"`
	SpamPrs                   int       `json:"spam_prs"`
	PRVelocity                int       `json:"pr_velocity"`
	CollaboratorAssociatedPrs int       `json:"collaborator_associated_prs"`
	ContributorAssociatedPrs  int       `json:"contributor_associated_prs"`
	MemberAssociatedPrs       int       `json:"member_associated_prs"`
	NonAssociatedPrs          int       `json:"non_associated_prs"`
	OwnerAssociatedPrs        int       `json:"owner_associated_prs"`
	CommentsOnPrs             int       `json:"comments_on_prs"`
	ReviewCommentsOnPrs       int       `json:"review_comments_on_prs"`
}
