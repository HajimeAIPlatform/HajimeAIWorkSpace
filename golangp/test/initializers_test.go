package test

import (
	"hajime/golangp/common/initializers"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	initializers.LoadEnv(".")
}
