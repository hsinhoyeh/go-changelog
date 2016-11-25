package changelog

import (
	"fmt"
	"io"
)

func Generate(w io.Writer, m map[string]Commits) {
	writeString(w, "##CHANGELOG\r\n")
	for head, commits := range m {
		if len(commits) == 0 {
			continue
		}
		writeString(w, fmt.Sprintf("####%s\r\n", head))
		for _, commit := range commits {
			writeString(w, fmt.Sprintf("[issue](%s)|[pull-request](%s) - %s\r\n\r\n",
				commit.IssueURL,
				commit.PrURL,
				commit.ShortMessage,
			))
		}
	}
}

func writeString(w io.Writer, str string) (int, error) {
	return w.Write([]byte(str))
}
