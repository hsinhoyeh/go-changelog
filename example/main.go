package main

import (
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/hsinhoyeh/go-changelog"
	"golang.org/x/oauth2"
)

func main() {
	const (
		token = "<SECRET>"
	)

	tc := oauth2.NewClient(
		oauth2.NoContext, oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: token},
		),
	)
	client := github.NewClient(tc)

	cl := changelog.NewChangeLog(client)

	commits, err := cl.Get("hsinhoyeh", "go-changelog", "base", "master")
	if err != nil {
		log.Fatal(err)
	}

	matchAny := func(_ changelog.Commit) bool {
		return true
	}

	list := map[string]changelog.Commits{
		"Features": commits.FindAll(matchAny),
	}

	f, err := os.OpenFile("CHANGELOG.md", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		log.Fatal(err)
	}
	changelog.Generate(f, list)
	f.Close()
}
