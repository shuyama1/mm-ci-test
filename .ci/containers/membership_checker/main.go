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
	if len(os.Args) <= 3 {
		fmt.Println("Not enough arguments")
		os.Exit(1)
	}

	target := os.Args[1]
	fmt.Println("PR Number: ", target)

	prNumber := os.Args[2]
	fmt.Println("PR Number: ", prNumber)

	commitSha := os.Args[3]
	fmt.Println("Commit SHA: ", commitSha)

	author, err := getPullRequestAuthor(prNumber, GITHUB_TOKEN)
	if err != nil {
		fmt.Println(err)
		return
	}

	if target == "auto_run_tests" {
		err = reviewerAssignment(author, prNumber, GITHUB_TOKEN)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	trusted := isTrustedUser(author, GITHUB_TOKEN)

	if (target == "auto_run_tests" && trusted) || (target == "check_community_contributor" && !trusted) {
		err = triggerMMPresubmitRuns("shuya-terraform-test", "mm-ci-test", commitSha)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
