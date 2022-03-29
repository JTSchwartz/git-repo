package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/pkg/browser"
	"github.com/urfave/cli/v2"
)

const configPath = "/.git/config"

func main() {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	app := &cli.App{
		Name:  "git-repo",
		Usage: "Open the current repo in browser",
		Action: func(context *cli.Context) (err error) {
			path := pwd + configPath

			url, err := getRepoUrl(path)
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

func getRepoUrl(path string) (url string, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	var line string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "remote \"origin\"") {
			break
		}
	}
	for scanner.Scan() {
		line = scanner.Text()
		if strings.Contains(line, "url = ") {
			break
		}
	}

	re := regexp.MustCompile("\\s+url = [^@]+@([^:]+):(.+)\\.git")

	// re := regexp.MustCompile("[^@]+@([^:]+):(.+)\\.git")
	// re := regexp.MustCompile("(.*)\\.git")

	parts := re.FindStringSubmatch(line)
	url = "https://" + parts[1] + "/" + parts[2] + parts[7]

	return
}
