package test

import (
	"hajime/golangp/libs"
	"testing"
)

func TestHelloFromLibsInDemo(t *testing.T) {
	expected := "Hello from libs!"
	if result := libs.Hello(); result != expected {
		t.Errorf("Expected %s, but got %s", expected, result)
	}
}
