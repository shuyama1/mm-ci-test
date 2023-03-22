package main

import (
    "testing"
    "golang.org/x/exp/slices"
)

func TestListCorrect(t *testing.T) {
    for _, member := range reviewerRotationList {
        if !slices.Contains(noAssigneeList, member) {
            t.Fatalf(`%v is not on noAssigneeList`, member)
        }
    }
}
