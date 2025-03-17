package fun_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Nick-Anderssohn/go-fun/fun"
)

func removeEmptyAndCapitalize(testSlice []string) ([]string, error) {
	return fun.NewSliceStream(testSlice).
		Filter(
			func(v string) (bool, error) {
				return v != "", nil
			},
		).
		Map(
			func(v string) (string, error) {
				return strings.ToUpper(v), nil
			},
		).
		Collect()
}

func TestEmptySlice(t *testing.T) {
	testSlice := []string{}

	result, err := removeEmptyAndCapitalize(testSlice)

	if err != nil {
		t.Errorf("empty string slice resulted in error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("empty slice is no longer empty somehow")
	}
}

func TestRemoveEmptyAndCapitalize(t *testing.T) {
	testSlice := []string{
		"",
		"Foo",
		"",
		"bar",
		"",
	}

	expectedResult := []string{
		"FOO",
		"BAR",
	}

	result, err := removeEmptyAndCapitalize(testSlice)

	if err != nil {
		t.Errorf("err: %v", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("expected: %v\nactual: %v", expectedResult, testSlice)
	}
}
