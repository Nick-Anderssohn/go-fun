package fun_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Nick-Anderssohn/go-fun/fun"
)

func removeEmptyKOrVAndCapitalizeMap(testMap map[string]string) (map[string]string, error) {
	return fun.M(testMap).
		Filter(
			func(k, v string) (bool, error) {
				return k != "" && v != "", nil
			},
		).
		Map(
			func(k, v string) (string, string, error) {
				return strings.ToUpper(k), strings.ToUpper(v), nil
			},
		).
		Collect()
}

func TestEmptyMap(t *testing.T) {
	testMap := map[string]string{}

	result, err := removeEmptyKOrVAndCapitalizeMap(testMap)

	if err != nil {
		t.Errorf("empty map slice resulted in error: %v", err)
	}

	if len(result) != 0 {
		t.Errorf("empty map is no longer empty somehow")
	}
}

func TestRemoveEmptyKOrVAndCapitalizeMap(t *testing.T) {
	testMap := map[string]string{
		"":       "Lol",
		"Lol":    "",
		"cheese": "Burger",
		"CHEEse": "burger",
		"YO":     "hey",
	}

	expectedResult := map[string]string{
		"CHEESE": "BURGER",
		"YO":     "HEY",
	}

	result, err := removeEmptyKOrVAndCapitalizeMap(testMap)

	if err != nil {
		t.Errorf("err: %v", err)
	}

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("expected: %v\nactual: %v", expectedResult, testMap)
	}
}
