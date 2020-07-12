package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func comment(noCodeBlock, noDetails bool, header, body string) string {
	var details string
	var comment string

	if header != "" {
		header = fmt.Sprintf("# %s\n", header)
	}

	if noCodeBlock {
		comment = body
	} else {
		comment = fmt.Sprintf("```\n%s\n```", body)
	}

	if noDetails {
		details = comment
	} else {
		details = fmt.Sprintf("\n<details>\n\n%s\n\n</details>", comment)
	}

	return fmt.Sprintf(`%s%s`, header, details)
}

func main() {
	owner := kingpin.Flag("owner", "GitHub owner").Required().String()
	repo := kingpin.Flag("repo", "GitHub repo").Required().String()
	number := kingpin.Flag("number", "GitHub issue number").Required().Int()
	header := kingpin.Flag("header", "GitHub issue comment header").String()
	noCodeBlock := kingpin.Flag("no-code-block", "no code block").Bool()
	noDetails := kingpin.Flag("no-details", "no details tag").Bool()
	kingpin.Parse()

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		log.Fatal(fmt.Errorf("GITHUB_TOKEN env is required"))
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	ctx := context.Background()
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	for {
		buf := &bytes.Buffer{}
		_, err := io.CopyN(buf, os.Stdin, 65535)

		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		body := comment(*noCodeBlock, *noDetails, *header, buf.String())
		comment := &github.IssueComment{
			Body: &body,
		}

		if _, _, cerr := client.Issues.CreateComment(ctx, *owner, *repo, *number, comment); cerr != nil {
			log.Fatal(cerr)
		}

		if err == io.EOF {
			break
		}
	}
}
