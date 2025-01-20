package main

import (
	"github.com/toqueteos/webbrowser"
	"os"
	"bufio"
	"fmt"
	"strings"
)

const (
	JIRAURL_FMT string = "https://mitel.atlassian.net/browse/%s"
)

// This was written primary to run as "jira bug-1234" and it would
// open the jira url for that bug for me.
func main() {
    // If invoked with an issue number, we do not need to prompt for one.
    issue := ""
    if len(os.Args) > 1 {
        issue = os.Args[1]
    } else {
        reader := bufio.NewReader(os.Stdin)
        writer := bufio.NewWriter(os.Stdout)
        writer.WriteString("Jira issue number? ")
        writer.Flush()
        line, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintf(os.Stderr, "%s\n", err)
        }
        issue = line
    }
	issue = strings.TrimSpace(issue)
	url := fmt.Sprintf(JIRAURL_FMT, issue)
	fmt.Printf("Opening %s in your browser\n", url)
	webbrowser.Open(url)
}
