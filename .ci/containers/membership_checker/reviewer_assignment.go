package main

import (
	"fmt"
	"net/http"
	"strings"
	_ "embed"
)

var (
    //go:embed REVIEWER_ASSIGNMENT_COMMENT.md
    reviewerAssignmentComment string
)

func reviewerAssignment(author, prNumber, GITHUB_TOKEN string) error {
	if isNoAssigneeUser(author) {
		fmt.Println("User is on the list, not assigning")
		return nil
	}

	requestedReviewer, err := getPullRequestRequestedReviewer(prNumber, GITHUB_TOKEN)
	if err != nil {
		return err
	}

	previousAssignedReviewers, err := getPullRequestPreviousAssignedReviewers(prNumber, GITHUB_TOKEN)
	if err != nil {
		return err
	}

	if requestedReviewer != "" {
		fmt.Println("Issue is assigned")
		if previousAssignedReviewers != nil {
			fmt.Println("Retrieving previous reviewers to re-request reviews")
			for _, reviewer := range previousAssignedReviewers {
				if isTeamReviewer(reviewer) {
					err = assignPullRequestReviewer(prNumber, reviewer, GITHUB_TOKEN)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}

	if previousAssignedReviewers == nil {
		assignRandomReviewer(prNumber, GITHUB_TOKEN)
	} else {
		foundTeamReviewer := false
		for _, reviewer := range previousAssignedReviewers {
			if isTeamReviewer(reviewer) {
				foundTeamReviewer = true
				err = assignPullRequestReviewer(prNumber, reviewer, GITHUB_TOKEN)
				if err != nil {
					return err
				}
			}
		}

		if !foundTeamReviewer {
			err = assignRandomReviewer(prNumber, GITHUB_TOKEN)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func getPullRequestAuthor(prNumber, GITHUB_TOKEN string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/shuyama1/mm-ci-test/issues/%s", prNumber)

	var pullRequest struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}

	_, err := requestCall(url, "GET", GITHUB_TOKEN, &pullRequest, nil)
	if err != nil {
		return "", err
	}

	return pullRequest.User.Login, nil
}

func getPullRequestRequestedReviewer(prNumber, GITHUB_TOKEN string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/shuyama1/mm-ci-test/pulls/%s/requested_reviewers", prNumber)

	var requestedReviewers struct {
		Users []struct {
			Login string `json:"login"`
		} `json:"users"`
	}

	_, err := requestCall(url, "GET", GITHUB_TOKEN, &requestedReviewers, nil)
	if err != nil {
		return "", err
	}

	if requestedReviewers.Users == nil || len(requestedReviewers.Users) == 0 {
		return "", nil
	}

	return requestedReviewers.Users[0].Login, nil
}

func getPullRequestPreviousAssignedReviewers(prNumber, GITHUB_TOKEN string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/shuyama1/mm-ci-test/pulls/%s/reviews", prNumber)

	var reviews []struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}

	_, err := requestCall(url, "GET", GITHUB_TOKEN, &reviews, nil)
	if err != nil {
		return nil, err
	}

	previousAssignedReviewers := map[string]struct{}{}
	for _, review := range reviews {
		previousAssignedReviewers[review.User.Login] = struct{}{}
	}

	result := []string{}
	for key, _ := range previousAssignedReviewers {
		result = append(result, key)
	}

	return result, nil
}

func assignPullRequestReviewer(prNumber, assignee, GITHUB_TOKEN string) error {
	url := fmt.Sprintf("https://api.github.com/repos/shuyama1/mm-ci-test/pulls/%s/requested_reviewers", prNumber)

	body := map[string][]string{
		"reviewers":      []string{assignee},
		"team_reviewers": []string{},
	}

	reqStatusCode, err := requestCall(url, "POST", GITHUB_TOKEN, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("Error adding reviewer for PR %s", prNumber)
	}

	fmt.Printf("Successfully added reviewer %s to pull request %s", assignee, prNumber)

	return nil
}

func assignRandomReviewer(prNumber, GITHUB_TOKEN string) error {
	assignee := getRamdomReviewer()
	err := assignPullRequestReviewer(prNumber, assignee, GITHUB_TOKEN)
	if err != nil {
		return err
	}
	err = postComment(prNumber, assignee, GITHUB_TOKEN)
	if err != nil {
		return err
	}
	return nil

}

func postComment(prNumber, reviewer, GITHUB_TOKEN string) error {
	url := fmt.Sprintf("https://api.github.com/repos/shuyama1/mm-ci-test/issues/%s/comments", prNumber)

	comment := strings.Replace(reviewerAssignmentComment, "{{reviewer}}", reviewer, 1)

	body := map[string]string{
		"body": comment,
	}

	reqStatusCode, err := requestCall(url, "POST", GITHUB_TOKEN, nil, body)
	if err != nil {
		return err
	}

	if reqStatusCode != http.StatusCreated {
		return fmt.Errorf("Error posting reviewer assignment comment for PR %s", prNumber)
	}

	fmt.Printf("Successfully posted reviewer assignment comment to pull request %s", prNumber)

	return nil
}
