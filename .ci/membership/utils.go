package membership

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
) 

func onList(user string, userList []string) bool {
	for _, v := range userList {
		if (v == user){
			return true
		}
	}
	return false
}

func RequestCall(url, method, GITHUB_TOKEN string, result interface{}, body interface{}) (int, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return 1, fmt.Errorf("Error marshaling JSON: %s", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 1, fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("token %s", GITHUB_TOKEN))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 1, err
	}
	defer resp.Body.Close()

	if result != nil {
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 1, err
		}
	}

	return resp.StatusCode, nil
}

func GetPullRequestAuthor(prNumber, GITHUB_TOKEN string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/shuyama1/mm-ci-test/issues/%s", prNumber)

	var pullRequest struct {
		User struct {
			Login string `json:"login"`
		} `json:"user"`
	}

	_, err := RequestCall(url, "GET", GITHUB_TOKEN, &pullRequest, nil)
	if err != nil{
		return "", err
	}

	return pullRequest.User.Login, nil
}

