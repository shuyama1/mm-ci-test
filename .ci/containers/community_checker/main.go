package main

import (
	"fmt"
	"os"
	"github.com/shuyama1/mm-ci-test/.ci/membership"
	// "example.com/membership"
)


func main() {
	fmt.Println("check if PR author is trusted user")
    GITHUB_TOKEN, ok := os.LookupEnv("GITHUB_TOKEN")
    if !ok {
        fmt.Println("Did not provide GITHUB_TOKEN environment variable")
        os.Exit(1)
    }
    if (len(os.Args) <= 2) {
        fmt.Println("Not enough arguments")
        os.Exit(1)
    }
    prNumber := os.Args[1]
    fmt.Println("PR Number: ", prNumber)

    commitSha := os.Args[2]
    fmt.Println("Commit SHA: ", commitSha)

	author, err := membership.GetPullRequestAuthor(prNumber, GITHUB_TOKEN)
	if err != nil{
	    fmt.Println(err)
	    return
	}

	if membership.IsTrustedUser(author, GITHUB_TOKEN){
	    fmt.Println("User is a trusted user. Presumit run should be automatically triggered")
	    return
	}

	err = membership.TriggerMMPresubmitRun("shuya-terraform-test", "mm-ci-test", commitSha)
	if err != nil{
	    fmt.Println(err)
	    return
	}
}