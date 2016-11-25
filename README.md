####go-changelog

package `go-changelog` implements a changelog from github in golang.

When cooperating with pull-request model using github, we usually leave review-comments and extra information in the pull-request itself.

However, when the pull-request get merged, those content didn't go to git's log message. In such case, we won't be able to see those valuable information in our changelog. 

This package do the following things:

1. compare two branches: head and base branch

2. visit each merged pull-request and get the useful content back
 
3. generate a CHANGELOG accordingly

###### INSTALL
```
go get github.com/hsinhoyeh/go-changelog
```

##### Example
```
func matchAllChanges(_ changelog.Commit) bool {
        return true
}

func main() {
        const (
                token = "b29cb7fdfb7a065f19e89b72731758e95be9c2b2"
        )
        tc := oauth2.NewClient(
                oauth2.NoContext, oauth2.StaticTokenSource(
                        &oauth2.Token{AccessToken: token},
                ),
        )
        client := github.NewClient(tc)

        cl := changelog.NewChangeLog(client)

        commits, err := cl.Get("hsinhoyeh", "go-changelog", "master", "base")
        if err != nil {
                log.Fatal(err)
        }

        matchAny:= func(_changelog.Commit) bool {
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
``` 