package main

import "testing"

func TestRun(t *testing.T) {
	inTest = true
	_, err := run()
	if err != nil {
		t.Error("Failed to run")
	}
}
