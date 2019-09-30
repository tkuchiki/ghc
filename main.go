package main

import (
	"bytes"
	"context"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io"
	"log"
	"os"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func comment(noCodeBlock bool, header, body string) string {
	var comment string
	if header != "" {
		header = fmt.Sprintf("# %s\n", header)
	}

	if noCodeBlock {
		comment = fmt.Sprintf("%s%s", header, body)
	} else {
		comment = fmt.Sprintf("%s```\n%s\n```", header, body)
	}

	return comment
}

func main() {
	owner := kingpin.Flag("owner", "GitHub owner").Required().String()
	repo := kingpin.Flag("repo", "GitHub repo").Required().String()
	number := kingpin.Flag("number", "GitHub issue number").Required().Int()
	header := kingpin.Flag("header", "GitHub issue comment header").String()
	noCodeBlock := kingpin.Flag("no-code-block", "no code block").Bool()
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
		body := comment(*noCodeBlock,  *header, buf.String())
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
