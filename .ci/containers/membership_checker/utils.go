package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"bytes"
	"embed"
) 

var (
    //go:embed *
    embededFiles embed.FS
)

func onList(user string, userList []string) bool {
	for _, v := range userList {
		if (v == user){
			return true
		}
	}
	return false
}

func requestCall(url, method, credentials string, result interface{}, body interface{}) (int, error) {
	client := &http.Client{}
	jsonBody, err := json.Marshal(body)

	if err != nil {
		return 1, fmt.Errorf("Error marshaling JSON: %s", err)
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 2, fmt.Errorf("Error creating request: %s", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", credentials))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return 3, err
	}
	defer resp.Body.Close()

	if result != nil {
		if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return 4, err
		}
	}

	return resp.StatusCode, nil
}

func readFile(filename string) (string, error) {
	contents, err := embededFiles.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(contents), nil
}

