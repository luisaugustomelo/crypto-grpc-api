package upvotetesting

import (
	"reflect"
	"testing"
)

func TestUpVote(t *testing.T) {
	type Test struct {
		cryptocurrency string
		vote           int
	}

	s := []*Test{
		{
			cryptocurrency: "Bitcoin",
			vote:           1,
		},
		{
			cryptocurrency: "Bitcoin",
			vote:           1,
		},
		{
			cryptocurrency: "Bitcoin",
			vote:           -1,
		},
		{
			cryptocurrency: "Bitcoin",
			vote:           1,
		},
		{
			cryptocurrency: "Iota",
			vote:           1,
		},
		{
			cryptocurrency: "Iota",
			vote:           1,
		},
		{
			cryptocurrency: "Ethereum",
			vote:           1,
		},
		{
			cryptocurrency: "Ethereum",
			vote:           -1,
		},
		{
			cryptocurrency: "Ethereum",
			vote:           -1,
		},
	}

	for _, x := range s {
		if !checkTypes(x.cryptocurrency, x.vote) {
			t.Error("Wrong type to validate vote")
		}
	}
}

func checkTypes(a string, b int) bool {
	if reflect.TypeOf(a).String() == "string" && reflect.TypeOf(b).String() == "int" {
		return true
	} else {
		return false
	}
}
