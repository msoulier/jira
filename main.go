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

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	writer.WriteString("Jira issue number? ")
	writer.Flush()
	line, err := reader.ReadString('\n')
    if err != nil {
        fmt.Fprintf(os.Stderr, "%s\n", err)
    }
	issue := strings.TrimSpace(line)
	url := fmt.Sprintf(JIRAURL_FMT, issue)
	fmt.Printf("Opening %s in your browser\n", url)
	webbrowser.Open(url)
}
