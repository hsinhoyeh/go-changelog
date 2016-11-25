package changelog

import (
	"strconv"
	"strings"

	"github.com/google/go-github/github"
)

// Commit holds git commit message and the corresponding pull request link
// and issue link
type Commit struct {
	SHA                string `json:"sha"`
	ShortMessage       string `json:"shortMessage"`
	CommitMessage      string `json:"commitMessage"`
	PullRequestContent string `json:"pullRequestContent"`
	PrURL              string `json:"prURL"`
	IssueURL           string `json:"issueURL"`
}

// Commits is a collection of commit
type Commits []Commit

// FindAll applies a commitmatcher to find all matched commits
func (cs Commits) FindAll(cm CommitMatcher) Commits {
	var found Commits

	for _, commit := range cs {
		if cm(commit) {
			found = append(found, commit)
		}
	}
	return found
}

type CommitMatcher func(Commit) bool

type ChangeLog struct {
	client *github.Client
}

func NewChangeLog(client *github.Client) *ChangeLog {
	return &ChangeLog{
		client: client,
	}
}

// Get calls api: https://developer.github.com/v3/repos/commits/#compare-two-commits
func (c *ChangeLog) Get(owner, repo, base, compare string) (Commits, error) {
	comp, _, err := c.client.Repositories.CompareCommits(owner, repo, base, compare)
	if err != nil {
		return nil, err
	}

	var commits Commits
	for _, remoteCommit := range comp.Commits {
		commit := Commit{}
		if remoteCommit.SHA != nil {
			commit.SHA = *remoteCommit.SHA
		}
		if remoteCommit.Commit.Message != nil {
			commit.CommitMessage = *remoteCommit.Commit.Message
			firstEOL := strings.Index(commit.CommitMessage, "\n\n")
			if firstEOL == -1 {
				commit.ShortMessage = commit.CommitMessage
			} else {
				commit.ShortMessage = commit.CommitMessage[0:firstEOL]
			}
		}

		// parse pull request number from (#xxx)
		commitMessage := *remoteCommit.Commit.Message
		leftParenthesis := strings.Index(commitMessage, "(")
		rightParenthesis := strings.Index(commitMessage, ")")
		if leftParenthesis != -1 && rightParenthesis != -1 {
			prNum, _ := strconv.Atoi(commitMessage[leftParenthesis+2 : /*skip ( and #*/ rightParenthesis])

			// get pull reqeust content
			// TODO: use batch api
			pr, _, err := c.client.PullRequests.Get(owner, repo, prNum)
			if err != nil {
				return nil, err
			}
			if pr.HTMLURL != nil {
				commit.PrURL = *pr.HTMLURL
			}
			if pr.Body != nil {
				commit.PullRequestContent = *pr.Body
			}
			if pr.IssueURL != nil {
				commit.IssueURL = *pr.IssueURL
			}
		}
		commits = append(commits, commit)
	}
	return commits, nil
}
