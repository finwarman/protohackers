package main

import "testing"

func TestRewriting(t *testing.T) {

	message := "7W01NPV8BW4xZnyLOBLlN9eQsNwAekbkI Please send the payment 7bWsVRNvWd7PaTlpo6gWX5A7kMHH 7mM5vVlq4TpN5SBnt3KTWLKT5Ywr7OjF of 750 Boguscoins to 7W01NPV8BW4xZnyLOBLlN9eQsNwAekbkI"
	result := RewriteCoins(message)

	expected := "7YWHMfk9JZe0LM0g1ZauHuiSxhI Please send the payment 7YWHMfk9JZe0LM0g1ZauHuiSxhI 7YWHMfk9JZe0LM0g1ZauHuiSxhI of 750 Boguscoins to 7YWHMfk9JZe0LM0g1ZauHuiSxhI"
	if result != expected {
		t.Fatalf("Fail:\nExpected:\n%s\n\nGot:\n%s\n", expected, result)
	}
}
