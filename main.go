package main

import (
	"log"
	"os"
	"regexp"

	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "git-repo",
		Usage: "Open the current repo in browser",
		Action: func(context *cli.Context) (err error) {
			url, err := getRepoUrl()
			if err != nil {
				return
			}
			err = browser.OpenURL(url)
			return
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getRepoUrl() (url string, err error) {
	unparsedUrl, err := GetGitConfig("remote.origin.url")
	if err != nil {
		return
	}

	re := regexp.MustCompile("[^@]+@([^:]+):(.+)\\.git")
	parts := re.FindStringSubmatch(unparsedUrl)
	if parts != nil {
		url = "https://" + parts[1] + "/" + parts[2]
	} else {
		re = regexp.MustCompile("(.*)\\.git")
		url = re.FindStringSubmatch(unparsedUrl)[1]
	}
	return
}
