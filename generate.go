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
			var buffer string
			if commit.IssueURL == "" {
				buffer += "[issue]"
			} else {
				buffer += fmt.Sprintf("[issue](%s)", commit.IssueURL)
			}
			buffer += "|"
			if commit.PrURL == "" {
				buffer += "[pull-request]"
			} else {
				buffer += fmt.Sprintf("[pull-request](%s)", commit.PrURL)
			}
			buffer += fmt.Sprintf(" - %s\r\n\r\n", commit.ShortMessage)

			writeString(w, buffer)
		}
	}
}

func writeString(w io.Writer, str string) (int, error) {
	return w.Write([]byte(str))
}
