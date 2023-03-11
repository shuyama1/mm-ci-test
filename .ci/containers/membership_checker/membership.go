package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	// This is where you add users who do not need to have an assignee chosen for them
	noAssigneeList = []string{"megan07", "slevenick", "c2thorn", "rileykarson", "melinath", "ScottSuarez", "shuyama1", "SarahFrench", "roaks3", "zli82016", "trodge", "hao-nan-li"}

	// This is where you add people to the random-assignee rotation.
	// reviewerRotationList = []string{"megan07", "slevenick", "c2thorn", "rileykarson", "melinath", "ScottSuarez", "shuyama1", "SarahFrench", "roaks3", "zli82016", "trodge", "hao-nan-li"}
	reviewerRotationList = []string{"shuyama1"}

	// This is where you add trusted users who do not need to '/gcbrun' comment to run tests
	trustedMemberList = []string{}
)

func isNoAssigneeUser(author string) bool {
	return onList(author, noAssigneeList)
}

func isTeamReviewer(reviewer string) bool {
	return onList(reviewer, reviewerRotationList)
}

// Check if a user is safe to run tests automatically
func isTrustedUser(author, GITHUB_TOKEN string) bool {
	if isTrustedMember(author) {
		fmt.Println("User is on the list")
		return true
	}

	if isOrgMember(author, "GoogleCloudPlatform", GITHUB_TOKEN){
		fmt.Println("User is a GCP org member")
		return true
	}

	if isOrgMember(author, "googlers", GITHUB_TOKEN){
		fmt.Println("User is a googlers org member")
		return true
	}

	return false
}

func isTrustedMember(author string) bool {
	return onList(author, trustedMemberList)
}


func isOrgMember(author, org, GITHUB_TOKEN string) bool {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/members/%s", org, author)
	res, _ := requestCall(url, "GET", GITHUB_TOKEN, nil, nil)

	return res!=404
}

func getRamdomReviewer() string{
	assignee := reviewerRotationList[rand.Intn(len(reviewerRotationList))]
    rand.Seed(time.Now().Unix())
    return assignee
}

