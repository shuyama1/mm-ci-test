package main

import (
	"fmt"
	"os"
)

func main() {
	GITHUB_TOKEN, ok := os.LookupEnv("GITHUB_TOKEN")
	if !ok {
		fmt.Println("Did not provide GITHUB_TOKEN environment variable")
		os.Exit(1)
	}
	if len(os.Args) <= 4 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	target := os.Args[1]
	fmt.Println("PR Number: ", target)

	prNumber := os.Args[2]
	fmt.Println("PR Number: ", prNumber)

	commitSha := os.Args[3]
	fmt.Println("Commit SHA: ", commitSha)

	branchName := os.Args[4]
	fmt.Println("Branch Name: ", branchName)

	projectId := "shuya-terraform-test"
	repoName := "mm-ci-test"

	author, err := getPullRequestAuthor(prNumber, GITHUB_TOKEN)
	if err != nil {
		fmt.Println(err)
		return
	}

	if target == "check_auto_run_contributor" {
		err = reviewerAssignment(author, prNumber, GITHUB_TOKEN)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	trusted := isTrustedUser(author, GITHUB_TOKEN)

	if (target == "check_auto_run_contributor" && trusted) || (target == "check_community_contributor" && !trusted) {
		err = triggerMMPresubmitRuns(projectId, repoName, commitSha, branchName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// in contributor-membership-checker job:
	// auto approve community-checker run for trusted users
	// add Awaiting Approval label to external contributor PRs
	if target == "check_auto_run_contributor" {
		if trusted {
			approveCommunityChecker(prNumber, projectId, commitSha)
		} else {
			addAwaitingApprovalLabel(prNumber, GITHUB_TOKEN)
		}
	}

	// in community-checker job:
	// remove Awaiting Approval label from external contributor PRs
	if target == "check_community_contributor" && !trusted {
		err = removeAwaitingApprovalLabel(prNumber, GITHUB_TOKEN)
	}
}


