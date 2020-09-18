package main

import (
	"context"
	"log"
	"strings"

	"github.com/google/go-github/github"
	"github.com/swinton/go-probot/probot"
)

func main() {
	probot.HandleEvent("issues", func(ctx *probot.Context) error {
		// Because we're listening for "issues" we know the payload is a *github.IssuesEvent
		event := ctx.Payload.(*github.IssuesEvent)
		log.Printf("🌈 Got issues %+v\n", event)

		return nil
	})

	probot.HandleEvent("pull_request", func(ctx *probot.Context) error {
		event := ctx.Payload.(*github.PullRequestEvent)
		pr := event.GetPullRequest()
		repo := event.GetRepo()
		log.Printf("Got the PR title: %s", *pr.Title)

		newReview := &github.PullRequestReviewRequest{CommitID: pr.Head.SHA}
		if strings.Contains(*pr.Title, "🤖") {
			newReview.Event = github.String("APPROVE")
			newReview.Body = github.String("LGTM 🚀")
		} else {
			newReview.Event = github.String("REQUEST_CHANGES")
			newReview.Body = github.String("Needs more 🤖")
		}

		review, _, err := ctx.GitHub.PullRequests.CreateReview(context.Background(), *repo.Owner.Login, *repo.Name, pr.GetNumber(), newReview)

		if err != nil {
			return nil
		}

		log.Printf("New review created: %+v", review)

		return nil
	})

	probot.Start()
}
