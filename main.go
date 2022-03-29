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
			unparsed := make(chan string, 1)
			parsed := make(chan string, 1)

			go parseRepoUrl(unparsed, parsed)
			err = getRepoUrl(unparsed)

			if err == nil {
				err = browser.OpenURL(<-parsed)
			}
			return
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func getRepoUrl(c chan string) (err error) {
	url, err := GetGitConfig("remote.origin.url")
	c <- url
	return
}

func parseRepoUrl(in chan string, out chan string) {
	unparsedUrl := <-in
	re := regexp.MustCompile("[^@]+@([^:]+):(.+)\\.git")
	parts := re.FindStringSubmatch(unparsedUrl)
	if parts != nil {
		out <- "https://" + parts[1] + "/" + parts[2]
	} else {
		re = regexp.MustCompile("(.*)\\.git")
		out <- re.FindStringSubmatch(unparsedUrl)[1]
	}
}
