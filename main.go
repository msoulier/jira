package main

import (
	"github.com/toqueteos/webbrowser"
	"os"
	"bufio"
	"fmt"
	"strings"
	"flag"
	"path"

	"github.com/op/go-logging"
    "github.com/go-ini/ini"
)

var (
	config string
	log          *logging.Logger
	debug        = false
	args []string
    jira_url_fmt = ""
	sample bool = false
)

func init() {
	default_config := "Override path to config file\n" +
		"($XDG_CONFIG_HOME/jira/jira.ini, $HOME/.jira.ini)"
	flag.StringVar(&config, "config", "", default_config)
	flag.BoolVar(&debug, "debug", false, "Debug logging")
	flag.BoolVar(&sample, "sample", false, "Display sample config")
	flag.Parse()
	args = flag.Args()

	if sample {
		sample_config()
		os.Exit(0)
	}

	if (config == "") {
		// If it's blank, then we go to defaults, in order
		// $XDG_CONFIG_HOME/jira/jira.ini
		// $HOME/.jira.ini
		// An array of possible config paths.
		config_paths := make([]string, 0)
		xdg_config_home := os.Getenv("XDG_CONFIG_HOME")
		if xdg_config_home != "" {
			config_paths = append(config_paths,
				path.Join(xdg_config_home, "jira", "jira.ini"))
		}
		home := os.Getenv("HOME")
		if home != "" {
			config_paths = append(config_paths,
				path.Join(home, ".jira.ini"))
		}
		// Now set the first one that exists
		for _, path := range config_paths {
			if _, err := os.Stat(path); os.IsNotExist(err) {
				// Does not exist
				continue
			} else {
				// It exists. Use it.
				config = path
				break
			}
		}
	}
	if config == "" {
		fmt.Fprintf(os.Stderr, "ERROR: Cannot determine path to config file\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Configure the logger
	format := logging.MustStringFormatter(
		`%{time:2006-01-02 15:04:05.000-0700} %{level} [%{shortfile}] %{message}`,
	)
	stderrBackend := logging.NewLogBackend(os.Stderr, "", 0)
	stderrFormatter := logging.NewBackendFormatter(stderrBackend, format)
	stderrBackendLevelled := logging.AddModuleLevel(stderrFormatter)
	logging.SetBackend(stderrBackendLevelled)
	if debug {
		stderrBackendLevelled.SetLevel(logging.DEBUG, "jira")
	} else {
		stderrBackendLevelled.SetLevel(logging.INFO, "jira")
	}
	log = logging.MustGetLogger("jira")

	log.Debugf("config path is %s", config)
}

func sample_config() {
	conf := `
[main]
jira_url_format = https://atlassian.com/browse/%s
`
	fmt.Printf("%s\n", conf)
}

func load_config() (*ini.File, error) {
    data, err := os.ReadFile(config)
    if err != nil {
        return nil, err
    }
    conffile, err := ini.Load(data)
    if err != nil {
        return nil, err
    }
    main, err := conffile.GetSection("main")
    if err != nil {
        return nil, err
    }
    if main.HasKey("jira_url_format") {
        log.Debug("main section has a jira_url_format key")
        urlfmt_key := main.Key("jira_url_format")
        urlfmt := urlfmt_key.Value()
        log.Debugf("jira_url_format from config is '%s'", urlfmt)
        jira_url_fmt = urlfmt
    }
    return conffile, nil
}

// This was written primary to run as "jira bug-1234" and it would
// open the jira url for that bug for me.
func main() {
    _, err := load_config()
    if err != nil {
        log.Errorf("config load error: %s", err)
        os.Exit(1)
    }
    // If invoked with an issue number, we do not need to prompt for one.
    issue := ""
    if len(args) > 0 {
        issue = args[0]
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
	url := fmt.Sprintf(jira_url_fmt, issue)
	fmt.Printf("Opening %s in your browser\n", url)
	webbrowser.Open(url)
}
